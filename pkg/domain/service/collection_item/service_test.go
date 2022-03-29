package collectionitem_test

import (
	"testing"

	cim "game-api/pkg/domain/model/collection_item"
	ucim "game-api/pkg/domain/model/user_collection_item"
	cimr "game-api/pkg/domain/repository/db/collection_item/mock"
	ucimr "game-api/pkg/domain/repository/db/user_collection_item/mock"
	cis "game-api/pkg/domain/service/collection_item"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mocks struct {
	collectionItemRepository     *cimr.MockRepository
	userCollectionItemRepository *ucimr.MockRepository
}

func newWithMocks(t *testing.T) (cis.Service, *mocks) {
	ctrl := gomock.NewController(t)
	collectionItemRepository := cimr.NewMockRepository(ctrl)
	userCollectionItemRepository := ucimr.NewMockRepository(ctrl)
	return cis.NewService(collectionItemRepository, userCollectionItemRepository), &mocks{
		collectionItemRepository:     collectionItemRepository,
		userCollectionItemRepository: userCollectionItemRepository,
	}
}

func TestSelectAllUserCollectionItems(t *testing.T) {
	type args struct {
		userID string
	}
	type expected struct {
		userCollection []*ucim.UserCollection
	}
	for name, tt := range map[string]struct {
		args     args
		prepare  func(f *mocks)
		expected expected
	}{
		"【正常系】ユーザコレクションの取得": {
			args: args{userID: "1234"},
			prepare: func(f *mocks) {
				f.collectionItemRepository.EXPECT().SelectAll().Return([]*cim.CollectionItem{
					{ID: "1001", Name: "スゴリラ01", Rarity: 1},
					{ID: "2001", Name: "レアスゴリラ01", Rarity: 2},
					{ID: "3001", Name: "超スゴリラ01", Rarity: 3},
				}, nil).Times(1)
				f.userCollectionItemRepository.EXPECT().SelectByUserID("1234").Return([]*ucim.UserCollectionItem{
					{UserID: "1234", CollectionItemID: "1001"},
					{UserID: "1234", CollectionItemID: "3001"},
				}, nil)
			},
			expected: expected{userCollection: []*ucim.UserCollection{
				{CollectionID: "1001", Name: "スゴリラ01", Rarity: 1, HasItem: true},
				{CollectionID: "2001", Name: "レアスゴリラ01", Rarity: 2, HasItem: false},
				{CollectionID: "3001", Name: "超スゴリラ01", Rarity: 3, HasItem: true},
			}},
		},
	} {
		t.Run(name, func(t *testing.T) {
			u, m := newWithMocks(t)
			tt.prepare(m)
			got, err := u.SelectAllUserCollectionItems(tt.args.userID)
			assert.Equal(t, tt.expected.userCollection, got)
			assert.NoError(t, err)
		})
	}
}
