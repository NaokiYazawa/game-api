package collectionitem

import (
	"log"
	"strings"

	cim "game-api/pkg/domain/model/collection_item"
	cir "game-api/pkg/domain/repository/db/collection_item"
	"game-api/pkg/infrastructure/mysql"
)

type repositoryImpl struct {
	mysql.SQLHandler
}

func NewRepositoryImpl(sqlHandler mysql.SQLHandler) cir.Repository {
	return &repositoryImpl{
		sqlHandler,
	}
}

func (ciri repositoryImpl) SelectAll() ([]*cim.CollectionItem, error) {
	var collectionItems []*cim.CollectionItem
	rows, err := ciri.SQLHandler.Conn.Query("SELECT id, name, rarity FROM collection_item")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		c := &cim.CollectionItem{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Rarity); err != nil {
			log.Println(err)
			return nil, err
		}
		collectionItems = append(collectionItems, c)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return collectionItems, nil
}

func (ciri repositoryImpl) SelectByCollectionIDs(collectionItemIDs []string) ([]*cim.CollectionItem, error) {
	if len(collectionItemIDs) == 0 {
		// 空のスライスを返すと親切
		return []*cim.CollectionItem{}, nil
	}
	queryString := "SELECT * FROM collection_item WHERE id IN ("
	queryArgs := make([]interface{}, 0, len(collectionItemIDs)*2)
	queryParams := make([]string, 0, len(collectionItemIDs))
	// collectionItemIDs の数だけプレースホルダーを用意
	for _, collectionItemID := range collectionItemIDs {
		queryParams = append(queryParams, "?")
		queryArgs = append(queryArgs, collectionItemID)
	}
	// カンマ区切りで結合させる
	queryString += strings.Join(queryParams, ",")
	// 閉じカッコを追加
	queryString += ")"
	rows, err := ciri.SQLHandler.Conn.Query(queryString, queryArgs...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var collectionItems []*cim.CollectionItem
	for rows.Next() {
		collectionItem := cim.CollectionItem{}
		if err := rows.Scan(&collectionItem.ID, &collectionItem.Name, &collectionItem.Rarity); err != nil {
			return nil, err
		}
		collectionItems = append(collectionItems, &collectionItem)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return collectionItems, nil
}
