package transaction

import (
	"context"
	"database/sql"
	"log"

	"game-api/pkg/domain/repository"
	"game-api/pkg/infrastructure/mysql"
)

var txKey = struct{}{}

type repositoryImpl struct {
	mysql.Tx
}

// NewRepositoryImpl Userに関するDB更新処理を生成
func NewRepositoryImpl(tx mysql.Tx) repository.Repository {
	return &repositoryImpl{
		tx,
	}
}

func (t *repositoryImpl) DoInTx(ctx context.Context, f func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	conn := t.GetDBConn()
	tx, err := conn.Begin() // {<nil>}
	// tx, err := t.Tx.Conn.Begin() // BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	// txを入れる
	ctx = context.WithValue(ctx, &txKey, tx)
	v, err := f(ctx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Printf("rollback is failed: %v", err)
		}
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Printf("rollback is failed: %v", err)
		}
		return nil, err
	}
	return v, nil
}

func GetTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(&txKey).(*sql.Tx)
	return tx, ok
}
