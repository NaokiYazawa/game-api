package ranking_test

import (
	"testing"

	rm "game-api/pkg/domain/model/ranking"
	um "game-api/pkg/domain/model/user"
	rmr "game-api/pkg/domain/repository/db/ranking/mock"

	rs "game-api/pkg/domain/service/ranking"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mocks struct {
	rankingRepository *rmr.MockRepository
}

func newWithMocks(t *testing.T) (rs.Service, *mocks) {
	ctrl := gomock.NewController(t)
	rankingRepository := rmr.NewMockRepository(ctrl)
	return rs.NewService(rankingRepository), &mocks{
		rankingRepository: rankingRepository,
	}
}
func TestSelectRankingList(t *testing.T) {
	type args struct {
		start int64
	}
	type expected struct {
		rankingList []*rm.Ranking
	}
	for name, tt := range map[string]struct {
		args     args
		prepare  func(f *mocks)
		expected expected
	}{
		// 今回は同率順位を考慮していません
		// Score が同じ場合には、redisに登録された順にランキングがつきます
		"【正常系】ランキングリストの取得": {
			args: args{start: int64(2)},
			prepare: func(f *mocks) {
				f.rankingRepository.EXPECT().SelectRankingList(int64(2)).Return([]*rm.Ranking{
					{UserID: "1111", UserName: "jack", Rank: 2, Score: 10000},
					{UserID: "2222", UserName: "john", Rank: 3, Score: 6000},
					{UserID: "3333", UserName: "george", Rank: 4, Score: 2000},
				}, nil).Times(1)
			},
			expected: expected{rankingList: []*rm.Ranking{
				{UserID: "1111", UserName: "jack", Rank: 2, Score: 10000},
				{UserID: "2222", UserName: "john", Rank: 3, Score: 6000},
				{UserID: "3333", UserName: "george", Rank: 4, Score: 2000},
			}},
		},
	} {
		t.Run(name, func(t *testing.T) {
			u, m := newWithMocks(t)
			tt.prepare(m)
			got, err := u.SelectRankingList(tt.args.start)
			assert.Equal(t, tt.expected.rankingList, got)
			assert.NoError(t, err)
		})
	}
}

func TestAddRankingList(t *testing.T) {
	type args struct {
		user *um.User
	}
	for name, tt := range map[string]struct {
		args    args
		prepare func(f *mocks)
	}{
		"【正常系】ランキングリストへの追加": {
			args: args{user: &um.User{ID: "1234", AuthToken: "authToken", Name: "sample", HighScore: 1000, Coin: 1000}},
			prepare: func(f *mocks) {
				f.rankingRepository.EXPECT().AddRankingList(&um.User{ID: "1234", AuthToken: "authToken", Name: "sample", HighScore: 1000, Coin: 1000}).Return(nil).Times(1)
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			u, m := newWithMocks(t)
			tt.prepare(m)
			err := u.AddRankingList(tt.args.user)
			assert.NoError(t, err)
		})
	}
}
