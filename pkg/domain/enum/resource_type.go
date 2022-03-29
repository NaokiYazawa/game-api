package enum

type ResourceType int32

const (
	ResourceTypeCollectionItem ResourceType = iota + 1 // 1と同義、0スタートだと未定義時(ゼロ値)と判断つかないため
	// ResourceType_Card (将来増えたら以下に記述していく)
	// ResourceType_Card
)
