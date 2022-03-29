package user

import (
	um "game-api/pkg/domain/model/user"
	ur "game-api/pkg/domain/repository/db/user"

	"github.com/google/uuid"
)

type Service interface {
	Create(name string) (authToken string, err error)
	SelectByAuthToken(authToken string) (user *um.User, err error)
	UpdateName(user *um.User, name string) error
}

type service struct {
	repository ur.Repository
}

func NewService(userRepo ur.Repository) Service {
	return &service{
		repository: userRepo,
	}
}

func (us *service) Create(name string) (string, error) {
	userID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	authToken, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	if err := us.repository.Create(userID.String(), authToken.String(), name); err != nil {
		return "", err
	}

	return authToken.String(), nil
}

func (us *service) SelectByAuthToken(authToken string) (*um.User, error) {
	return us.repository.SelectByAuthToken(authToken)
}

func (us *service) UpdateName(user *um.User, name string) error {
	user.Name = name
	return us.repository.Update(user)
}
