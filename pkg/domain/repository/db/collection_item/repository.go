package collectionitem

import (
	cim "game-api/pkg/domain/model/collection_item"
)

type Repository interface {
	SelectAll() ([]*cim.CollectionItem, error)
	SelectByCollectionIDs(collectionItemIDs []string) ([]*cim.CollectionItem, error)
}
