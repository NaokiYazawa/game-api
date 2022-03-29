package gachaprobability

import (
	"log"

	gpm "game-api/pkg/domain/model/gacha_probability"
	gpr "game-api/pkg/domain/repository/db/gacha_probability"
	"game-api/pkg/infrastructure/mysql"
)

type repositoryImpl struct {
	mysql.SQLHandler
}

func NewRepositoryImpl(sqlHandler mysql.SQLHandler) gpr.Repository {
	return &repositoryImpl{
		sqlHandler,
	}
}

func (gpri repositoryImpl) SelectAll() ([]*gpm.GachaProbability, error) {
	rows, err := gpri.SQLHandler.Conn.Query("SELECT collection_item_id, ratio FROM gacha_probability")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	// gachaProbabilitySlice：全ての GachaProbability が格納される
	var gachaProbabilitySlice []*gpm.GachaProbability
	for rows.Next() {
		// gachaProbabilitySlice に 1 つ 1 つの gachaProbability を追加していく
		gachaProbability := &gpm.GachaProbability{}
		if err := rows.Scan(&gachaProbability.CollectionItemID, &gachaProbability.Ratio); err != nil {
			log.Println(err)
			return nil, err
		}
		gachaProbabilitySlice = append(gachaProbabilitySlice, gachaProbability)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return gachaProbabilitySlice, nil
}
