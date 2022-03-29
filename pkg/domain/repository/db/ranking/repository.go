package ranking

import (
	rm "game-api/pkg/domain/model/ranking"
	um "game-api/pkg/domain/model/user"
)

type Repository interface {
	SelectRankingList(start int64) ([]*rm.Ranking, error)
	AddRankingList(user *um.User) error
}
