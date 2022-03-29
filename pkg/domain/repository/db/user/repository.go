package user

import (
	"context"

	"game-api/pkg/domain/model/user"
)

type Repository interface {
	Create(ID, authToken, name string) error
	SelectByAuthToken(authToken string) (*user.User, error)
	SelectByPrimaryKey(userID string) (*user.User, error)
	Update(user *user.User) error
	UpdateWithLock(ctx context.Context, user *user.User) error
}
