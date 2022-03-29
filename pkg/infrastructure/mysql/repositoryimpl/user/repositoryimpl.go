package user

import (
	"context"
	"database/sql"
	"fmt"

	um "game-api/pkg/domain/model/user"
	ur "game-api/pkg/domain/repository/db/user"
	"game-api/pkg/infrastructure/mysql"
	"game-api/pkg/infrastructure/mysql/repositoryimpl/transaction"
	"game-api/pkg/interfaces/api/myerror"
)

type repositoryImpl struct {
	mysql.SQLHandler
}

// NewRepositoryImpl Userに関するDB更新処理を生成
func NewRepositoryImpl(sqlHandler mysql.SQLHandler) ur.Repository {
	return &repositoryImpl{
		sqlHandler,
	}
}

// Create ユーザ登録処理
func (uri *repositoryImpl) Create(id, authToken, name string) error {
	prep, err := uri.SQLHandler.Conn.Prepare("INSERT INTO `user` (`id`, `auth_token`, `name`, `high_score`, `coin`) VALUES (?, ?, ?, 0, 0)")
	if err != nil {
		return err
	}

	_, err = prep.Exec(id, authToken, name)
	return err
}

// SelectByAuthToken auth_tokenを条件にUserを取得する
func (uri *repositoryImpl) SelectByAuthToken(authToken string) (*um.User, error) {
	row := uri.SQLHandler.Conn.QueryRow("SELECT * FROM `user` WHERE `auth_token` = ?", authToken)
	return convertToUser(row)
}

// SelectByPrimaryKey user_IDを条件にUserを取得する
func (uri *repositoryImpl) SelectByPrimaryKey(userID string) (*um.User, error) {
	row := uri.SQLHandler.Conn.QueryRow("SELECT * FROM `user` WHERE `id` = ?", userID)
	return convertToUser(row)
}

// Update ユーザーを更新する
func (uri *repositoryImpl) Update(record *um.User) error {
	prep, err := uri.SQLHandler.Conn.Prepare("UPDATE user SET name=?, high_score=?, coin=? WHERE id=?")
	if err != nil {
		return err
	}
	_, err = prep.Exec(record.Name, record.HighScore, record.Coin, record.ID)
	return err
}

// UpdateWithLock ユーザを排他制御で更新する
func (uri *repositoryImpl) UpdateWithLock(ctx context.Context, record *um.User) error {
	var dbtx mysql.DBTx
	dbtx, ok := transaction.GetTx(ctx)
	if !ok {
		dbtx = uri.Conn
		fmt.Println("transaction context is null")
	}
	stmt, err := dbtx.Prepare("UPDATE `user` SET `auth_token`=?, `name`=?, `high_score`=?, `coin`=? WHERE `id`=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(record.AuthToken, record.Name, record.HighScore, record.Coin, record.ID)
	return err
}

func convertToUser(row *sql.Row) (*um.User, error) {
	user := um.User{}
	if err := row.Scan(&user.ID, &user.AuthToken, &user.Name, &user.HighScore, &user.Coin); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, &myerror.InternalServerError{Err: err}
	}
	return &user, nil
}
