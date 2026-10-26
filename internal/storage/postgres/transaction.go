package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transactor struct {
	db *pgxpool.Pool
}

func NewTransactor(db *pgxpool.Pool) *Transactor {
	return &Transactor{db: db}
}

type executor interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

type txKey struct{}

// Executor возвращает транзакцию из ctx, а если её нет, использует пул.
// Передавай контекст из InTx, чтобы запрос участвовал в его транзакции.
func Executor(ctx context.Context, db *pgxpool.Pool) executor {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}

	return db
}

// InTx выполняет fn в транзакции, передавая её через контекст.
// Вложенный вызов использует ту же транзакцию без отдельной точки сохранения.
// Для отката ошибка вложенного вызова должна дойти до внешнего InTx.
func (t *Transactor) InTx(ctx context.Context, fn func(context.Context) error) error {
	if _, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return fn(ctx)
	}

	err := pgx.BeginFunc(ctx, t.db, func(tx pgx.Tx) error {
		return fn(context.WithValue(ctx, txKey{}, tx))
	})
	if err != nil {
		return MapError(ctx, fmt.Errorf("transaction: %w", err))
	}

	return nil
}
