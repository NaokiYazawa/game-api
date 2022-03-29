package collectionitem

import (
	"net/http"

	cis "game-api/pkg/domain/service/collection_item"
	"game-api/pkg/interfaces/api/dcontext"
	"game-api/pkg/interfaces/api/myerror"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// Handler UserにおけるHandlerのインターフェース
// 本アプリは「get_collection_list」の1種類
type Handler interface {
	HandleGet(c echo.Context) error
}

type handler struct {
	service cis.Service
}

func NewHandler(collectionItemService cis.Service) Handler {
	return &handler{
		service: collectionItemService,
	}
}

func (cih handler) HandleGet(c echo.Context) error {
	type collectionItem struct {
		CollectionID string `json:"collectionID"`
		Name         string `json:"name"`
		Rarity       int32  `json:"rarity"`
		HasItem      bool   `json:"hasItem"`
	}

	type response struct {
		Collections []*collectionItem `json:"collections"`
	}

	user := dcontext.GetUserFromContext(c)
	if user == nil {
		return &myerror.UnauthorizedError{Err: errors.New("user not found")}
	}

	collections, err := cih.service.SelectAllUserCollectionItems(user.ID)

	if err != nil {
		return &myerror.InternalServerError{Err: err}
	}

	res := make([]*collectionItem, len(collections))
	for i, c := range collections {
		res[i] = &collectionItem{
			CollectionID: c.CollectionID,
			Name:         c.Name,
			Rarity:       int32(c.Rarity),
			HasItem:      c.HasItem,
		}
	}

	return c.JSON(http.StatusOK, &response{
		res,
	})
}
