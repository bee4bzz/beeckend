package testutils

import (
	"context"
	"net/http"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func (r *Repository) Get(ctx context.Context, cheptel *entity.Cheptel) error {
	args := r.Called(*cheptel)
	return args.Error(0)
}

func (r *Repository) QueryByUser(ctx context.Context, user *entity.User, cheptels *[]entity.Cheptel) error {
	args := r.Called(*user, *cheptels)
	*cheptels = args.Get(0).([]entity.Cheptel)
	return args.Error(1)
}

func (r *Repository) Create(ctx context.Context, cheptel *entity.Cheptel) error {
	args := r.Called(*cheptel)
	return args.Error(0)
}

func (r *Repository) Update(ctx context.Context, cheptel *entity.Cheptel) error {
	args := r.Called(*cheptel)
	return args.Error(0)
}
func (r *Repository) SoftDelete(ctx context.Context, cheptel *entity.Cheptel) error {
	args := r.Called(*cheptel)
	return args.Error(0)
}

func (r *Repository) Subscribe(ctx context.Context, user *entity.User, cheptels chan<- *[]entity.Cheptel) error {
	args := r.Called(*user)
	return args.Error(0)
}

type Service struct {
	mock.Mock
}

func (s *Service) Subscribe(ctx context.Context, user *entity.User, cheptels chan<- *[]entity.Cheptel) error {
	args := s.Called(user)
	if args.Error(1) != nil {
		return args.Error(1)
	}
	values := args.Get(0).([]*[]entity.Cheptel)
	go func() {
		for _, v := range values {
			cheptels <- v
		}
		close(cheptels)
	}()
	return nil
}

type Pool struct {
	mock.Mock
}

func (p *Pool) Acquire(ctx context.Context) (*ConnPool, error) {
	args := p.Called()
	return args.Get(0).(*ConnPool), args.Error(1)
}

type ConnPool struct {
	mock.Mock
}

func (p *ConnPool) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	args := p.Called(sql)
	return args.Get(0).(pgconn.CommandTag), args.Error(1)
}

func (p *ConnPool) Conn() *Conn {
	args := p.Called()
	return args.Get(0).(*Conn)
}

func (p *ConnPool) Release() {
	p.Called()
}

type Conn struct {
	mock.Mock
}

func (p *Conn) WaitForNotification(ctx context.Context) (*pgconn.Notification, error) {
	args := p.Called()
	return args.Get(0).(*pgconn.Notification), args.Error(1)
}

type Upgrader struct {
	mock.Mock
}

func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*WSConn, error) {
	args := u.Called(w, r, responseHeader)
	return args.Get(0).(*WSConn), args.Error(1)
}

type WSConn struct {
	mock.Mock
}

func (c *WSConn) WriteMessage(messageType int, data []byte) error {
	args := c.Called(messageType, data)
	return args.Error(0)
}

func (c *WSConn) Close() error {
	args := c.Called()
	return args.Error(0)
}
