package gachaprobability

type GachaProbability struct {
	CollectionItemID string
	Ratio            int32
}

type GachaDrawResult struct {
	CollectionID string
	Name         string
	Rarity       int32
	IsNew        bool
}
