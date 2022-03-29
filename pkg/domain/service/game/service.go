package game

import (
	"errors"

	rr "game-api/pkg/domain/repository/db/ranking"
	ur "game-api/pkg/domain/repository/db/user"
	"game-api/pkg/interfaces/api/myerror"
)

type Service interface {
	GameFinish(authToken string, score int32) (coin int32, err error)
}

type service struct {
	userRepository    ur.Repository
	rankingRepository rr.Repository
}

func NewService(userRepo ur.Repository, rankingRepo rr.Repository) Service {
	return &service{
		userRepository:    userRepo,
		rankingRepository: rankingRepo,
	}
}

func (gs *service) GameFinish(authToken string, score int32) (int32, error) {
	user, err := gs.userRepository.SelectByAuthToken(authToken)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, &myerror.InternalServerError{Err: errors.New("user not found")}
	}
	// 今回獲得したスコアが最高得点かどうかを判定
	if user.HighScore < score {
		user.HighScore = score
	}
	// 今回獲得したコインを計算（今回はscoreの1倍とする）
	coin := score
	// 実行userにコインを付与
	user.Coin += coin
	if err := gs.userRepository.Update(user); err != nil {
		return 0, err
	}
	// ユーザをランキングに追加
	if err := gs.rankingRepository.AddRankingList(user); err != nil {
		return 0, err
	}
	// 今回獲得したコインを返す
	return coin, nil
}
