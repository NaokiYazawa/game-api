package reward

import "game-api/pkg/domain/enum"

// Reward 抽選対象
type Reward struct {
	ResourceID   string
	ResourceType enum.ResourceType
	Ratio        int32
}

type Rewards []*Reward
