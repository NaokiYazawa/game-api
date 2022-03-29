package usercollectionitem

import (
	"context"

	ucim "game-api/pkg/domain/model/user_collection_item"
)

type Repository interface {
	SelectByUserID(userID string) ([]*ucim.UserCollectionItem, error)
	InsertWithLock(ctx context.Context, userCollectionItems []*ucim.UserCollectionItem) error
	SelectByUserIDAndCollectionIDs(userID string, collectionItemIDs []string) ([]*ucim.UserCollectionItem, error)
}
