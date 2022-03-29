package ranking

import (
	rm "game-api/pkg/domain/model/ranking"
	um "game-api/pkg/domain/model/user"
	rr "game-api/pkg/domain/repository/db/ranking"
)

type Service interface {
	SelectRankingList(start int64) ([]*rm.Ranking, error)
	AddRankingList(user *um.User) error
}

type service struct {
	rankingRepository rr.Repository
}

func NewService(rankingRepo rr.Repository) Service {
	return &service{
		rankingRepository: rankingRepo,
	}
}

func (rs *service) SelectRankingList(start int64) ([]*rm.Ranking, error) {
	return rs.rankingRepository.SelectRankingList(start)
}

func (rs *service) AddRankingList(user *um.User) error {
	return rs.rankingRepository.AddRankingList(user)
}
