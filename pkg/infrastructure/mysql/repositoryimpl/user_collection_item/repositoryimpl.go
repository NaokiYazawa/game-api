package usercollectionitem

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	ucim "game-api/pkg/domain/model/user_collection_item"
	ucir "game-api/pkg/domain/repository/db/user_collection_item"
	"game-api/pkg/infrastructure/mysql"
	"game-api/pkg/infrastructure/mysql/repositoryimpl/transaction"
)

type repositoryImpl struct {
	mysql.SQLHandler
}

// NewRepositoryImpl Userに関するDB更新処理を生成
func NewRepositoryImpl(sqlHandler mysql.SQLHandler) ucir.Repository {
	return &repositoryImpl{
		sqlHandler,
	}
}

// SelectUserCollectionItemsByUserID user_idと一致する行全て取得
func (uciri repositoryImpl) SelectByUserID(userID string) ([]*ucim.UserCollectionItem, error) {
	rows, err := uciri.SQLHandler.Conn.Query("SELECT * FROM user_collection_item WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	return convertToUserCollectionItems(rows)
}

// InsertUserCollectionItemsWithLock
func (uciri repositoryImpl) InsertWithLock(ctx context.Context, userCollectionItems []*ucim.UserCollectionItem) error {
	var dbtx mysql.DBTx
	dbtx, ok := transaction.GetTx(ctx)
	if !ok {
		dbtx = uciri.Conn
		fmt.Println("transaction context is null")
	}
	queryString := "INSERT INTO `user_collection_item`(`user_id`, `collection_item_id`) VALUES"
	queryArgs := make([]interface{}, 0, len(userCollectionItems)*2)
	queryParams := make([]string, 0, len(userCollectionItems))
	for _, userCollectionItem := range userCollectionItems {
		queryParams = append(queryParams, "(?,?)")
		queryArgs = append(queryArgs, userCollectionItem.UserID, userCollectionItem.CollectionItemID)
	}
	queryString += strings.Join(queryParams, ",")
	stmtTx, err := dbtx.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = stmtTx.Exec(queryArgs...)
	return err
}

func (uciri repositoryImpl) SelectByUserIDAndCollectionIDs(userID string, collectionItemIDs []string) ([]*ucim.UserCollectionItem, error) {
	var queryString string
	queryArgs := make([]string, 0, len(collectionItemIDs))
	for _, collectionItemID := range collectionItemIDs {
		queryArgs = append(queryArgs, fmt.Sprintf("('%s', '%s')", userID, collectionItemID))
	}
	queryString += strings.Join(queryArgs, ",")
	rows, err := uciri.SQLHandler.Conn.Query(fmt.Sprintf(`SELECT user_id, collection_item_id FROM user_collection_item WHERE (user_id, collection_item_id) IN (%s)`, queryString))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return convertToUserCollectionItems(rows)
}

// convertToUserCollectionItems rowsデータをUserCollectionItemsデータへ変換する
func convertToUserCollectionItems(rows *sql.Rows) ([]*ucim.UserCollectionItem, error) {
	var userCollectionItems []*ucim.UserCollectionItem
	for rows.Next() {
		userCollectionItem := ucim.UserCollectionItem{}
		if err := rows.Scan(&userCollectionItem.UserID, &userCollectionItem.CollectionItemID); err != nil {
			return nil, err
		}
		userCollectionItems = append(userCollectionItems, &userCollectionItem)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return userCollectionItems, nil
}
