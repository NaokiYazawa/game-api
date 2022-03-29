package usercollectionitem

type UserCollectionItem struct {
	UserID           string
	CollectionItemID string
}

type UserCollection struct {
	CollectionID string
	Name         string
	Rarity       int
	HasItem      bool
}
