package gacha_test

import (
	"context"
	"errors"
	"testing"

	rewarddto "game-api/pkg/domain/dto/reward"
	cim "game-api/pkg/domain/model/collection_item"
	gpm "game-api/pkg/domain/model/gacha_probability"
	um "game-api/pkg/domain/model/user"
	ucim "game-api/pkg/domain/model/user_collection_item"
	tmr "game-api/pkg/domain/repository"
	cimr "game-api/pkg/domain/repository/db/collection_item/mock"
	gpmr "game-api/pkg/domain/repository/db/gacha_probability/mock"
	umr "game-api/pkg/domain/repository/db/user/mock"
	ucimr "game-api/pkg/domain/repository/db/user_collection_item/mock"
	fms "game-api/pkg/domain/service/component/mock"
	gs "game-api/pkg/domain/service/gacha"
	"game-api/pkg/interfaces/api/myerror"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mocks struct {
	userRepository               *umr.MockRepository
	collectionItemRepository     *cimr.MockRepository
	userCollectionItemRepository *ucimr.MockRepository
	gachaProbabilityRepository   *gpmr.MockRepository
	transactionRepository        *tmr.MockRepository
	facade                       *fms.MockFacade
}

func newWithMocks(t *testing.T) (gs.Service, *mocks) {
	ctrl := gomock.NewController(t)
	userRepository := umr.NewMockRepository(ctrl)
	collectionItemRepository := cimr.NewMockRepository(ctrl)
	userCollectionItemRepository := ucimr.NewMockRepository(ctrl)
	gachaProbabilityRepository := gpmr.NewMockRepository(ctrl)
	transactionRepository := tmr.NewMockRepository(ctrl)
	facade := fms.NewMockFacade(ctrl)
	return gs.NewService(userRepository, collectionItemRepository, userCollectionItemRepository, gachaProbabilityRepository, transactionRepository, facade), &mocks{
		userRepository:               userRepository,
		collectionItemRepository:     collectionItemRepository,
		userCollectionItemRepository: userCollectionItemRepository,
		gachaProbabilityRepository:   gachaProbabilityRepository,
		transactionRepository:        transactionRepository,
		facade:                       facade,
	}
}

func TestGachaDraw(t *testing.T) {
	type args struct {
		ctx       context.Context
		authToken string
		times     int32
	}
	type expected struct {
		gachaDrawResults []*gpm.GachaDrawResult
		err              error
	}
	for name, tt := range map[string]struct {
		args     args
		prepare  func(f *mocks)
		expected expected
	}{
		"【正常系】重複がないアイテムが当選するガチャ（3回）": {
			args: args{ctx: context.Background(), authToken: "authToken", times: 3},
			prepare: func(f *mocks) {
				f.userRepository.EXPECT().SelectByAuthToken("authToken").Return(&um.User{ID: "1234", AuthToken: "authToken", Name: "sample", HighScore: 1000, Coin: 1000}, nil).Times(1)
				f.gachaProbabilityRepository.EXPECT().SelectAll().Return([]*gpm.GachaProbability{
					{CollectionItemID: "1001", Ratio: 6},
					{CollectionItemID: "1002", Ratio: 6},
					{CollectionItemID: "2001", Ratio: 3},
					{CollectionItemID: "2002", Ratio: 3},
					{CollectionItemID: "3001", Ratio: 1},
					{CollectionItemID: "3002", Ratio: 1},
				}, nil).Times(1)
				// ガチャを引く
				f.facade.EXPECT().ChooseByRatio(rewarddto.Rewards{
					{ResourceID: "1001", ResourceType: 1, Ratio: 6},
					{ResourceID: "1002", ResourceType: 1, Ratio: 6},
					{ResourceID: "2001", ResourceType: 1, Ratio: 3},
					{ResourceID: "2002", ResourceType: 1, Ratio: 3},
					{ResourceID: "3001", ResourceType: 1, Ratio: 1},
					{ResourceID: "3002", ResourceType: 1, Ratio: 1},
				}, int32(3)).Return([]*rewarddto.Reward{
					{ResourceID: "1001", ResourceType: 1, Ratio: 6},
					{ResourceID: "2001", ResourceType: 1, Ratio: 3},
					{ResourceID: "3001", ResourceType: 1, Ratio: 1},
				}, nil).Times(1)
				// 当選したアイテム
				f.collectionItemRepository.EXPECT().SelectByCollectionIDs([]string{"1001", "2001", "3001"}).Return([]*cim.CollectionItem{
					{ID: "1001", Name: "スゴリラ01", Rarity: 1},
					{ID: "2001", Name: "レアスゴリラ01", Rarity: 2},
					{ID: "3001", Name: "超スゴリラ01", Rarity: 3},
				}, nil).Times(1)
				// ユーザが所持しているアイテム
				f.userCollectionItemRepository.EXPECT().SelectByUserIDAndCollectionIDs("1234", []string{"1001", "2001", "3001"}).Return([]*ucim.UserCollectionItem{
					{UserID: "1234", CollectionItemID: "1001"},
					{UserID: "1234", CollectionItemID: "2001"},
				}, nil)
				// トランザクションのモック
				f.transactionRepository.EXPECT().DoInTx(context.Background(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) (interface{}, error)) (interface{}, error) {
						return f(ctx)
					}).Times(1)
				// ユーザ更新処理
				f.userRepository.EXPECT().UpdateWithLock(context.Background(), &um.User{ID: "1234", AuthToken: "authToken", Name: "sample", HighScore: 1000, Coin: 700}).Return(nil).Times(1)
				// ユーザアイテム更新処理
				f.userCollectionItemRepository.EXPECT().InsertWithLock(context.Background(), []*ucim.UserCollectionItem{
					{UserID: "1234", CollectionItemID: "3001"},
				}).Return(nil).Times(1)
			},
			expected: expected{gachaDrawResults: []*gpm.GachaDrawResult{
				{CollectionID: "1001", Name: "スゴリラ01", Rarity: 1, IsNew: false},
				{CollectionID: "2001", Name: "レアスゴリラ01", Rarity: 2, IsNew: false},
				{CollectionID: "3001", Name: "超スゴリラ01", Rarity: 3, IsNew: true},
			}, err: nil},
		},
		"【正常系】重複するアイテムが当選するガチャ（3回）": {
			args: args{ctx: context.Background(), authToken: "authToken", times: 3},
			prepare: func(f *mocks) {
				f.userRepository.EXPECT().SelectByAuthToken("authToken").Return(&um.User{ID: "1234", AuthToken: "authToken", Name: "sample", HighScore: 1000, Coin: 1000}, nil).Times(1)
				f.gachaProbabilityRepository.EXPECT().SelectAll().Return([]*gpm.GachaProbability{
					{CollectionItemID: "1001", Ratio: 6},
					{CollectionItemID: "1002", Ratio: 6},
					{CollectionItemID: "2001", Ratio: 3},
					{CollectionItemID: "2002", Ratio: 3},
					{CollectionItemID: "3001", Ratio: 1},
					{CollectionItemID: "3002", Ratio: 1},
				}, nil).Times(1)
				// ガチャを引く
				f.facade.EXPECT().ChooseByRatio(rewarddto.Rewards{
					{ResourceID: "1001", ResourceType: 1, Ratio: 6},
					{ResourceID: "1002", ResourceType: 1, Ratio: 6},
					{ResourceID: "2001", ResourceType: 1, Ratio: 3},
					{ResourceID: "2002", ResourceType: 1, Ratio: 3},
					{ResourceID: "3001", ResourceType: 1, Ratio: 1},
					{ResourceID: "3002", ResourceType: 1, Ratio: 1},
				}, int32(3)).Return([]*rewarddto.Reward{
					{ResourceID: "1001", ResourceType: 1, Ratio: 6},
					{ResourceID: "2001", ResourceType: 1, Ratio: 3},
					{ResourceID: "2001", ResourceType: 1, Ratio: 3},
				}, nil).Times(1)
				// 当選したアイテム
				f.collectionItemRepository.EXPECT().SelectByCollectionIDs([]string{"1001", "2001", "2001"}).Return([]*cim.CollectionItem{
					{ID: "1001", Name: "スゴリラ01", Rarity: 1},
					{ID: "2001", Name: "レアスゴリラ01", Rarity: 2},
				}, nil).Times(1)
				// ユーザが所持しているアイテム
				f.userCollectionItemRepository.EXPECT().SelectByUserIDAndCollectionIDs("1234", []string{"1001", "2001", "2001"}).Return([]*ucim.UserCollectionItem{
					{UserID: "1234", CollectionItemID: "1001"},
				}, nil)
				// トランザクションのモック
				f.transactionRepository.EXPECT().DoInTx(context.Background(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) (interface{}, error)) (interface{}, error) {
						return f(ctx)
					}).Times(1)
				// ユーザ更新処理
				f.userRepository.EXPECT().UpdateWithLock(context.Background(), &um.User{ID: "1234", AuthToken: "authToken", Name: "sample", HighScore: 1000, Coin: 700}).Return(nil).Times(1)
				// ユーザアイテム更新処理
				f.userCollectionItemRepository.EXPECT().InsertWithLock(context.Background(), []*ucim.UserCollectionItem{
					{UserID: "1234", CollectionItemID: "2001"},
				}).Return(nil).Times(1)
			},
			expected: expected{gachaDrawResults: []*gpm.GachaDrawResult{
				{CollectionID: "1001", Name: "スゴリラ01", Rarity: 1, IsNew: false},
				{CollectionID: "2001", Name: "レアスゴリラ01", Rarity: 2, IsNew: true},
			}, err: nil},
		},
		"【正常系】新たに獲得するアイテムが存在しないガチャ（3回）": {
			args: args{ctx: context.Background(), authToken: "authToken", times: 3},
			prepare: func(f *mocks) {
				f.userRepository.EXPECT().SelectByAuthToken("authToken").Return(&um.User{ID: "1234", AuthToken: "authToken", Name: "sample", HighScore: 1000, Coin: 1000}, nil).Times(1)
				f.gachaProbabilityRepository.EXPECT().SelectAll().Return([]*gpm.GachaProbability{
					{CollectionItemID: "1001", Ratio: 6},
					{CollectionItemID: "1002", Ratio: 6},
					{CollectionItemID: "2001", Ratio: 3},
					{CollectionItemID: "2002", Ratio: 3},
					{CollectionItemID: "3001", Ratio: 1},
					{CollectionItemID: "3002", Ratio: 1},
				}, nil).Times(1)
				// ガチャを引く
				f.facade.EXPECT().ChooseByRatio(rewarddto.Rewards{
					{ResourceID: "1001", ResourceType: 1, Ratio: 6},
					{ResourceID: "1002", ResourceType: 1, Ratio: 6},
					{ResourceID: "2001", ResourceType: 1, Ratio: 3},
					{ResourceID: "2002", ResourceType: 1, Ratio: 3},
					{ResourceID: "3001", ResourceType: 1, Ratio: 1},
					{ResourceID: "3002", ResourceType: 1, Ratio: 1},
				}, int32(3)).Return([]*rewarddto.Reward{
					{ResourceID: "1001", ResourceType: 1, Ratio: 6},
					{ResourceID: "2001", ResourceType: 1, Ratio: 3},
					{ResourceID: "3001", ResourceType: 1, Ratio: 1},
				}, nil).Times(1)
				// 当選したアイテム
				f.collectionItemRepository.EXPECT().SelectByCollectionIDs([]string{"1001", "2001", "3001"}).Return([]*cim.CollectionItem{
					{ID: "1001", Name: "スゴリラ01", Rarity: 1},
					{ID: "2001", Name: "レアスゴリラ01", Rarity: 2},
					{ID: "3001", Name: "超スゴリラ01", Rarity: 3},
				}, nil).Times(1)
				// ユーザが所持しているアイテム
				f.userCollectionItemRepository.EXPECT().SelectByUserIDAndCollectionIDs("1234", []string{"1001", "2001", "3001"}).Return([]*ucim.UserCollectionItem{
					{UserID: "1234", CollectionItemID: "1001"},
					{UserID: "1234", CollectionItemID: "2001"},
					{UserID: "1234", CollectionItemID: "3001"},
				}, nil)
				// トランザクションのモック
				f.transactionRepository.EXPECT().DoInTx(context.Background(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, f func(context.Context) (interface{}, error)) (interface{}, error) {
						return f(ctx)
					}).Times(1)
				// ユーザ更新処理
				f.userRepository.EXPECT().UpdateWithLock(context.Background(), &um.User{ID: "1234", AuthToken: "authToken", Name: "sample", HighScore: 1000, Coin: 700}).Return(nil).Times(1)
			},
			expected: expected{gachaDrawResults: []*gpm.GachaDrawResult{
				{CollectionID: "1001", Name: "スゴリラ01", Rarity: 1, IsNew: false},
				{CollectionID: "2001", Name: "レアスゴリラ01", Rarity: 2, IsNew: false},
				{CollectionID: "3001", Name: "超スゴリラ01", Rarity: 3, IsNew: false},
			}, err: nil},
		},
		"【異常系】ユーザのコインが足りない場合のガチャ（3回）": {
			args: args{ctx: context.Background(), authToken: "authToken", times: 3},
			prepare: func(f *mocks) {
				f.userRepository.EXPECT().SelectByAuthToken("authToken").Return(&um.User{ID: "1234", AuthToken: "authToken", Name: "sample", HighScore: 1000, Coin: 200}, nil).Times(1)
			},
			expected: expected{gachaDrawResults: nil, err: errors.New("you don't have enough coins")},
		},
		"【異常系】ユーザが見つからなかった場合（3回）": {
			args: args{ctx: context.Background(), authToken: "invalidToken", times: 3},
			prepare: func(f *mocks) {
				f.userRepository.EXPECT().SelectByAuthToken("invalidToken").Return(nil, nil).Times(1)
			},
			expected: expected{gachaDrawResults: nil, err: &myerror.InternalServerError{Err: errors.New("user not found")}},
		},
	} {
		t.Run(name, func(t *testing.T) {
			u, m := newWithMocks(t)
			tt.prepare(m)
			got, err := u.GachaDraw(tt.args.ctx, tt.args.authToken, tt.args.times)
			assert.Equal(t, tt.expected.gachaDrawResults, got)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}
