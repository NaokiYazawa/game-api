package facade

import (
	"math/rand"

	rewarddto "game-api/pkg/domain/dto/reward"
)

type Facade interface {
	ChooseByRatio(rewards rewarddto.Rewards, times int32) (rewarddto.Rewards, error)
}

type facade struct{}

func NewFacade() Facade {
	return &facade{}
}

func (f *facade) ChooseByRatio(rewards rewarddto.Rewards, times int32) (rewarddto.Rewards, error) {
	// ..ratioによって抽選する処理...
	var sumRatio int32
	// sumRatio：Ratio の合計
	for _, reward := range rewards {
		sumRatio += reward.Ratio
	}
	// rewardsSlice に獲得した Item の CollectionItemID を格納していく
	rewardsSlice := make(rewarddto.Rewards, 0, times)

	for i := 0; i < int(times); i++ {
		// 負でない疑似乱数を [0,n) で int32 で返す。 n <= 0 の場合はパニックする。
		randomRatio := rand.Int31n(sumRatio)
		// ratio の初期化
		var ratio int32

		for _, reward := range rewards {
			// 仕様書：あるコレクションアイテムの排出確率=あるコレクションアイテムのratio/全体のratio合計
			// したがって、gachaProbability.Ratio をプラスしていく必要がある
			ratio += reward.Ratio
			if randomRatio <= ratio {
				rewardsSlice = append(rewardsSlice, reward)
				break
			}
		}
	}
	return rewardsSlice, nil
}
