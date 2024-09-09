package db

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DB represents a DB connection that can be used to run SQL queries.
type DB struct {
	db *gorm.DB
}

// TransactionFunc represents a function that will start a transaction and run the given function.
type TransactionFunc func(ctx context.Context, f func(ctx context.Context) error) error

type contextKey int

const (
	TxKey contextKey = iota
)

// New returns a new DB connection that wraps the given dbx.DB instance.
func New(db *gorm.DB) *DB {
	return &DB{db}
}

// DB returns the dbx.DB wrapped by this object.
func (db *DB) DB() *gorm.DB {
	return db.db
}

// With returns a Builder that can be used to build and execute SQL queries.
// With will return the transaction if it is found in the given context.
// Otherwise it will return a DB connection associated with the context.
func (db *DB) With(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(TxKey).(*gorm.DB); ok {
		return tx
	}
	return db.db.WithContext(ctx)
}

// TransactionHandler returns a middleware that starts a transaction.
// The transaction started is kept in the context and can be accessed via With().
func (db *DB) TransactionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db.db.Transaction(func(tx *gorm.DB) error {
			timeoutContext, cancel := context.WithTimeout(c.Request.Context(), time.Second)
			defer cancel()
			ctx := context.WithValue(timeoutContext, TxKey, tx)
			c.Request = c.Request.WithContext(ctx)
			c.Next()

			// Check if the request has been aborted
			if c.IsAborted() {
				// Return an error to abort the transaction
				return errors.New("aborted")
			}

			return nil
		})
	}
}
