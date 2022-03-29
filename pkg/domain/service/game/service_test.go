package game_test

import (
	"errors"
	"testing"

	gs "game-api/pkg/domain/service/game"
	"game-api/pkg/interfaces/api/myerror"

	rmr "game-api/pkg/domain/repository/db/ranking/mock"
	umr "game-api/pkg/domain/repository/db/user/mock"

	um "game-api/pkg/domain/model/user"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mocks struct {
	userRepository    *umr.MockRepository
	rankingRepository *rmr.MockRepository
}

func newWithMocks(t *testing.T) (gs.Service, *mocks) {
	ctrl := gomock.NewController(t)
	userRepository := umr.NewMockRepository(ctrl)
	rankingRepository := rmr.NewMockRepository(ctrl)
	return gs.NewService(userRepository, rankingRepository), &mocks{
		userRepository:    userRepository,
		rankingRepository: rankingRepository,
	}
}

func TestGameFinish(t *testing.T) {
	type args struct {
		authToken string
		score     int32
	}
	type expected struct {
		coin int32
		err  error
	}
	for name, tt := range map[string]struct {
		args     args
		prepare  func(f *mocks)
		expected expected
	}{
		"【正常系】ゲーム終了（ハイスコア更新）": {
			args: args{authToken: "authToken", score: 1000},
			prepare: func(f *mocks) {
				f.userRepository.EXPECT().SelectByAuthToken("authToken").Return(&um.User{ID: "id", AuthToken: "authToken", Name: "name", HighScore: 0, Coin: 0}, nil).Times(1)
				f.userRepository.EXPECT().Update(&um.User{ID: "id", AuthToken: "authToken", Name: "name", HighScore: 1000, Coin: 1000}).Return(nil).Times(1)
				f.rankingRepository.EXPECT().AddRankingList(&um.User{ID: "id", AuthToken: "authToken", Name: "name", HighScore: 1000, Coin: 1000}).Return(nil).Times(1)
			},
			expected: expected{coin: 1000, err: nil},
		},
		"【正常系】ゲーム終了（ハイスコア未更新）": {
			args: args{authToken: "authToken", score: 100},
			prepare: func(f *mocks) {
				f.userRepository.EXPECT().SelectByAuthToken("authToken").Return(&um.User{ID: "id", AuthToken: "authToken", Name: "name", HighScore: 1000, Coin: 1000}, nil).Times(1)
				f.userRepository.EXPECT().Update(&um.User{ID: "id", AuthToken: "authToken", Name: "name", HighScore: 1000, Coin: 1100}).Return(nil).Times(1)
				f.rankingRepository.EXPECT().AddRankingList(&um.User{ID: "id", AuthToken: "authToken", Name: "name", HighScore: 1000, Coin: 1100}).Return(nil).Times(1)
			},
			expected: expected{coin: 100, err: nil},
		},
		"【異常系】ユーザが見つからなかった場合": {
			args: args{authToken: "invalidToken", score: 100},
			prepare: func(f *mocks) {
				f.userRepository.EXPECT().SelectByAuthToken("invalidToken").Return(nil, nil).Times(1)
			},
			expected: expected{coin: 0, err: &myerror.InternalServerError{Err: errors.New("user not found")}},
		},
	} {
		t.Run(name, func(t *testing.T) {
			u, m := newWithMocks(t)
			tt.prepare(m)
			got, err := u.GameFinish(tt.args.authToken, tt.args.score)
			assert.Equal(t, tt.expected.coin, got)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}
