package repository

import (
	"context"
	"fmt"
	"strings"

	dbx "github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/log"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm/clause"
)

type Conn interface {
	WaitForNotification(ctx context.Context) (*pgconn.Notification, error)
}

type ConnPool[T Conn] interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Conn() T
	Release()
}

type Pool[T ConnPool[V], V Conn] interface {
	Acquire(ctx context.Context) (T, error)
}

type GormRepository[T ConnPool[V], V Conn] struct {
	*repository.Repository[entity.Cheptel]
	Pool[T, V]
	Logger log.Logger
}

func NewGormRepository[T ConnPool[V], V Conn](db *dbx.DB, pool Pool[T, V], logger log.Logger) *GormRepository[T, V] {
	return &GormRepository[T, V]{
		Repository: repository.NewRepository[entity.Cheptel](db),
		Pool:       pool,
		Logger:     logger,
	}
}

func (r *GormRepository[T, V]) QueryByUser(ctx context.Context, user *entity.User, cheptels *[]entity.Cheptel) error {
	return r.DB().Model(user).
		Preload(clause.Associations).Association(entity.CheptelsKey).Find(cheptels)
}

func (r *GormRepository[T, V]) Subscribe(ctx context.Context, user *entity.User, cheptels chan<- *[]entity.Cheptel) error {
	logger := r.Logger.With(
		ctx, "channel", fmt.Sprintf("%s%v", entity.CheptelsKey, user.ID),
	)

	conn, err := r.Pool.Acquire(ctx)
	if err != nil {
		logger.Error("Error acquiring connection:", err)
		return err
	}

	_, err = conn.Exec(ctx, fmt.Sprintf("LISTEN %s%v", entity.CheptelsKey, user.ID))
	if err != nil {
		logger.Error("cannot listen to channel: ", err)
		return err
	}

	newCheptels := &[]entity.Cheptel{}
	err = r.QueryByUser(ctx, user, newCheptels)
	if err != nil {
		logger.Error("Error querying cheptels:", err)
		return err
	}

	uuids := []string{}
	for _, cheptel := range *newCheptels {
		uuids = append(uuids, fmt.Sprint(cheptel.ID))
	}

	uuidsForSQL := "'" + strings.Join(uuids, "', '") + "'"
	_, err = conn.Exec(ctx, fmt.Sprintf(`
			CREATE OR REPLACE TRIGGER trigger_%s_user_%v
			AFTER UPDATE OR DELETE ON cheptels
			FOR EACH ROW
			WHEN (OLD.id IN (%s))
			EXECUTE FUNCTION notify_changes(%s%v);
			
			CREATE OR REPLACE TRIGGER trigger_%s_user_%v
			AFTER UPDATE OR DELETE ON user_cheptels
			FOR EACH ROW
			WHEN (OLD.cheptel_id IN (%s)
			OR OLD.user_id = %v)
			EXECUTE FUNCTION notify_changes(%s%v);

			CREATE OR REPLACE TRIGGER trigger_%s_user_%v_insert
			AFTER INSERT ON user_cheptels
			FOR EACH ROW
			WHEN (NEW.cheptel_id IN (%s)
			OR NEW.user_id = %v)
			EXECUTE FUNCTION notify_changes(%s%v);
			`,
		entity.CheptelsKey, user.ID, uuidsForSQL, entity.CheptelsKey, user.ID,
		entity.UserCheptelsKey, user.ID, uuidsForSQL, user.ID, entity.CheptelsKey, user.ID,
		entity.UserCheptelsKey, user.ID, uuidsForSQL, user.ID, entity.CheptelsKey, user.ID,
	))
	if err != nil {
		logger.Error("cannot create trigger: ", err)
		return err
	}

	c := conn.Conn()

	newContext := context.WithoutCancel(ctx)

	go func() {
		defer conn.Release()
		defer close(cheptels)
		defer func() {
			_, err = conn.Exec(newContext, fmt.Sprintf("UNLISTEN %s%v", entity.CheptelsKey, user.ID))
			if err != nil {
				logger.Error("cannot unlisten channel: ", err)
			}

			_, err = conn.Exec(newContext, fmt.Sprintf("DROP TRIGGER trigger_%s_user_%v ON cheptels",
				entity.CheptelsKey, user.ID))
			if err != nil {
				logger.Error("cannot drop trigger: ", err)
			}
		}()

		cheptels <- newCheptels

		for {
			_, err := c.WaitForNotification(ctx)
			if err != nil {
				logger.Error("Error waiting for notification:", err)
				return
			}

			newCheptels := &[]entity.Cheptel{}
			err = r.QueryByUser(ctx, user, newCheptels)
			if err != nil {
				logger.Error("Error querying cheptels:", err)
				return
			}
			cheptels <- newCheptels
		}
	}()

	return nil
}
