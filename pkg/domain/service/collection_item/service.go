package collectionitem

import (
	ucim "game-api/pkg/domain/model/user_collection_item"
	cir "game-api/pkg/domain/repository/db/collection_item"
	ucir "game-api/pkg/domain/repository/db/user_collection_item"
)

type Service interface {
	SelectAllUserCollectionItems(userID string) ([]*ucim.UserCollection, error)
}

type service struct {
	collectionRepository         cir.Repository
	userCollectionItemRepository ucir.Repository
}

// NewUseCase Userデータに関するユースケースを生成
func NewService(collectionItemRepo cir.Repository, userCollectionItemRepo ucir.Repository) Service {
	return &service{
		collectionRepository:         collectionItemRepo,
		userCollectionItemRepository: userCollectionItemRepo,
	}
}

func (cis *service) SelectAllUserCollectionItems(userID string) ([]*ucim.UserCollection, error) {
	collectionItems, err := cis.collectionRepository.SelectAll()
	if err != nil {
		return nil, err
	}

	// userIDに紐づくUserCollectionItemを全件取得
	userCollectionItems, err := cis.userCollectionItemRepository.SelectByUserID(userID)
	if err != nil {
		return nil, err
	}

	// CollectionItem_IDでkey判別するためのMAPを作成
	userCollectionItemMap := make(map[string]struct{}, len(userCollectionItems))
	for _, userCollectionItem := range userCollectionItems {
		// struct{}{} は空の構造体
		// userCollectionItemMap の key に CollectionID を設定して、value に空の構造体を設定
		userCollectionItemMap[userCollectionItem.CollectionItemID] = struct{}{}
	}

	// response作成
	collections := make([]*ucim.UserCollection, 0, len(collectionItems))
	for _, collectionItem := range collectionItems {
		_, hasItem := userCollectionItemMap[collectionItem.ID]
		collections = append(collections, &ucim.UserCollection{
			CollectionID: collectionItem.ID,
			Name:         collectionItem.Name,
			Rarity:       int(collectionItem.Rarity),
			HasItem:      hasItem,
		})
	}

	return collections, nil
}
