package gachaprobability

import (
	gpm "game-api/pkg/domain/model/gacha_probability"
)

type Repository interface {
	SelectAll() ([]*gpm.GachaProbability, error)
}
