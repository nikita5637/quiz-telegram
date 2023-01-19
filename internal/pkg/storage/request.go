//go:generate mockery --case underscore --name RequestStorage --with-expecter

package storage

import (
	"context"
	"database/sql"

	"github.com/nikita5637/quiz-telegram/internal/config"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"github.com/nikita5637/quiz-telegram/internal/pkg/storage/mysql"

	"github.com/go-xorm/builder"
)

// RequestStorage ...
type RequestStorage interface {
	Delete(ctx context.Context, id int32) error
	Find(ctx context.Context, q builder.Cond) ([]model.Request, error)
	GetRequestByUUID(ctx context.Context, uuid string) (model.Request, error)
	Insert(ctx context.Context, request model.Request) (uint64, string, error)
}

// NewRequestStorage ...
func NewRequestStorage(db *sql.DB) RequestStorage {
	switch config.GetValue("Driver").String() {
	case config.DriverMySQL:
		return mysql.NewRequestsStorage(db)
	}

	return nil
}
