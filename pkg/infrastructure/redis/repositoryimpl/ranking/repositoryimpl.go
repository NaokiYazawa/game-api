package ranking

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"game-api/pkg/constant"
	rm "game-api/pkg/domain/model/ranking"
	um "game-api/pkg/domain/model/user"
	ri "game-api/pkg/domain/repository/db/ranking"
	redisInfra "game-api/pkg/infrastructure/redis"

	"github.com/go-redis/redis"
)

type repositoryImpl struct {
	redisInfra.CacheHandler
}

func NewRepositoryImpl(cacheHandler redisInfra.CacheHandler) ri.Repository {
	return &repositoryImpl{
		cacheHandler,
	}
}

// GetRankingList はソート済みマップから指定範囲のランキングを取得する
func (rri repositoryImpl) SelectRankingList(start int64) ([]*rm.Ranking, error) {
	redisZList, err := rri.Client.ZRevRangeWithScores(constant.RankingKey, start-1, start+constant.RankingUserLimit-2).Result()
	if err != nil {
		return nil, err
	}
	return convertToRanking(redisZList, start)
}

func (rri repositoryImpl) AddRankingList(user *um.User) error {
	userJSON, _ := json.Marshal(user)
	if err := rri.Client.ZAdd(constant.RankingKey, redis.Z{
		Score:  float64(user.HighScore),
		Member: userJSON,
	}).Err(); err != nil {
		log.Printf("ZAdd userdata is failed: %v", err)
		return err
	}
	return nil
}

func convertToRanking(zlist []redis.Z, start int64) ([]*rm.Ranking, error) {
	rankingList := make([]*rm.Ranking, 0, constant.RankingUserLimit)
	for i, Z := range zlist {
		user := um.User{}
		userJSON, ok := Z.Member.(string)
		if !ok {
			log.Println("casting is failed ... not string")
			return nil, errors.New("casting is failed")
		}

		dec := json.NewDecoder(strings.NewReader(userJSON))
		if err := dec.Decode(&user); err != nil {
			log.Printf("Decode error :%v", err)
			return nil, err
		}

		rankingList = append(rankingList, &rm.Ranking{
			UserID:   user.ID,
			UserName: user.Name,
			Rank:     int32(i + int(start)),
			Score:    user.HighScore,
		})
	}
	return rankingList, nil
}
