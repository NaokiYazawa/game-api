package gacha

import (
	"context"
	"errors"
	"fmt"
	"log"

	"game-api/pkg/constant"
	cim "game-api/pkg/domain/model/collection_item"
	gpm "game-api/pkg/domain/model/gacha_probability"
	um "game-api/pkg/domain/model/user"
	ucim "game-api/pkg/domain/model/user_collection_item"
	"game-api/pkg/interfaces/api/myerror"

	txr "game-api/pkg/domain/repository"
	cir "game-api/pkg/domain/repository/db/collection_item"
	gpr "game-api/pkg/domain/repository/db/gacha_probability"
	ur "game-api/pkg/domain/repository/db/user"
	ucir "game-api/pkg/domain/repository/db/user_collection_item"

	rewarddto "game-api/pkg/domain/dto/reward"
	"game-api/pkg/domain/enum"
	fs "game-api/pkg/domain/service/component"
)

type Service interface {
	GachaDraw(ctx context.Context, authToken string, times int32) ([]*gpm.GachaDrawResult, error)
}

type service struct {
	userRepository               ur.Repository
	collectionItemRepository     cir.Repository
	userCollectionItemRepository ucir.Repository
	gachaProbabilityRepository   gpr.Repository
	transactionRepository        txr.Repository
	facade                       fs.Facade
}

func NewService(userRepo ur.Repository, collectionItemRepo cir.Repository, userCollectionItemRepo ucir.Repository, gachaProbabilityRepo gpr.Repository, txRepo txr.Repository, facade fs.Facade) Service {
	return &service{
		userRepository:               userRepo,
		collectionItemRepository:     collectionItemRepo,
		userCollectionItemRepository: userCollectionItemRepo,
		gachaProbabilityRepository:   gachaProbabilityRepo,
		transactionRepository:        txRepo,
		facade:                       facade,
	}
}

// updateCoinAndItem は，ユーザのコイン消費とアイテム更新を行います．
func (gs *service) updateCoinAndItem(user *um.User, userCollectionItems []*ucim.UserCollectionItem) func(ctx context.Context) (interface{}, error) {
	return func(ctx context.Context) (interface{}, error) {
		fmt.Println(ctx)
		userRepository := gs.userRepository
		userCollectionItemRepository := gs.userCollectionItemRepository
		// ユーザの更新処理を排他制御で実装する
		if err := userRepository.UpdateWithLock(ctx, user); err != nil {
			return nil, err
		}
		// アイテム更新
		// itemRepository.InsertUserItemTx では，contextからtxを取得して実行する
		if len(userCollectionItems) > 0 {
			if err := userCollectionItemRepository.InsertWithLock(ctx, userCollectionItems); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}
}

func (gs *service) GachaDraw(ctx context.Context, authToken string, times int32) ([]*gpm.GachaDrawResult, error) {
	user, err := gs.userRepository.SelectByAuthToken(authToken)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, &myerror.InternalServerError{Err: errors.New("user not found")}
	}
	if user.Coin < (constant.GachaCoinConsumption * times) {
		return nil, errors.New("you don't have enough coins")
	}
	// 全ての GachaProbability を取得する
	gachaProbabilities, err := gs.gachaProbabilityRepository.SelectAll()
	if err != nil {
		return nil, err
	}

	// Rewardへ詰め替える処理
	rewards := make(rewarddto.Rewards, 0, len(gachaProbabilities))
	for _, v := range gachaProbabilities {
		rewards = append(rewards, &rewarddto.Reward{ResourceID: v.CollectionItemID, ResourceType: enum.ResourceTypeCollectionItem, Ratio: v.Ratio})
	}

	fixedRewards, err := gs.facade.ChooseByRatio(rewards, times) // 注入したfacedeの抽選関数を呼ぶ
	if err != nil {
		return nil, err
	}

	collectionItemIDs := make([]string, 0, len(fixedRewards))
	for _, v := range fixedRewards {
		collectionItemIDs = append(collectionItemIDs, v.ResourceID)
	}
	// collectionItemIDs := lotteryCollectionItems(times, gachaProbabilities)

	//　ガチャを引いて取得した collectionItemIDs から、CollectionItem を取得する
	collectionItems, err := gs.collectionItemRepository.SelectByCollectionIDs(collectionItemIDs)
	if err != nil {
		return nil, err
	}

	// user_collection_item テーブルから userID をもとにユーザのアイテム所持情報を取得
	// 工夫: collectionItemIDs を元に抽選で当たった必要最低限のデータを取得する
	// userCollectionItems には、UserID と CollectionID をフィールドに持つ UserCollectionItem の配列が返ってくる
	userCollectionItems, err := gs.userCollectionItemRepository.SelectByUserIDAndCollectionIDs(user.ID, collectionItemIDs)
	if err != nil {
		return nil, err
	}

	// 新たに獲得したアイテムを格納する
	var newUserCollectionItems []*ucim.UserCollectionItem

	var gachaResults []*gpm.GachaDrawResult

	// - ガチャの結果の作成
	// - newUserCollectionItems に今回獲得したアイテムを格納する
	gachaResults, newUserCollectionItems = createGachaResultAndNewUserCollectionItems(user, collectionItems, userCollectionItems)

	// コインを消費する
	user.Coin -= constant.GachaCoinConsumption * times

	// トランザクション処理
	txRepository := gs.transactionRepository
	if _, err := txRepository.DoInTx(ctx, gs.updateCoinAndItem(user, newUserCollectionItems)); err != nil {
		log.Println("updateCoinAndItem is failed")
		return nil, err
	}

	return gachaResults, nil
}

func createGachaResultAndNewUserCollectionItems(
	user *um.User,
	collectionItems []*cim.CollectionItem,
	userCollectionItems []*ucim.UserCollectionItem,
) ([]*gpm.GachaDrawResult, []*ucim.UserCollectionItem) {
	// ユーザにとって新しいアイテムを格納する
	newUserCollectionItems := make([]*ucim.UserCollectionItem, 0, len(collectionItems))
	userCollectionItemMap := make(map[string]struct{}, len(userCollectionItems))
	for _, userCollectionItem := range userCollectionItems {
		userCollectionItemMap[userCollectionItem.CollectionItemID] = struct{}{}
	}
	// gachaResults の初期化
	gachaResults := make([]*gpm.GachaDrawResult, 0, len(collectionItems))
	for _, collectionItem := range collectionItems {
		_, hasItem := userCollectionItemMap[collectionItem.ID]
		gachaResult := &gpm.GachaDrawResult{
			CollectionID: collectionItem.ID,
			Name:         collectionItem.Name,
			Rarity:       collectionItem.Rarity,
			IsNew:        !hasItem,
		}
		// newUserCollectionItems への追加
		// 今回獲得したアイテムが新しいかつ、重複がない場合、newUserCollectionItems に追加する
		if gachaResult.IsNew {
			newUserCollectionItems = append(newUserCollectionItems, &ucim.UserCollectionItem{
				UserID:           user.ID,
				CollectionItemID: collectionItem.ID,
			})
			// 今回獲得したアイテムを userCollectionItemMap に追加
			userCollectionItemMap[collectionItem.ID] = struct{}{}
		}
		// gachaResults への追加
		gachaResults = append(gachaResults, gachaResult)
	}
	return gachaResults, newUserCollectionItems
}
