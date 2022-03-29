package gacha

import (
	"net/http"

	gs "game-api/pkg/domain/service/gacha"
	"game-api/pkg/interfaces/api/dcontext"
	"game-api/pkg/interfaces/api/myerror"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// Handler UserにおけるHandlerのインターフェース
// 本アプリは「get_collection_list」の1種類
type Handler interface {
	HandleGacha(c echo.Context) error
}

type handler struct {
	service gs.Service
}

func NewHandler(collectionItemService gs.Service) Handler {
	return &handler{
		service: collectionItemService,
	}
}

func (gh handler) HandleGacha(c echo.Context) error {
	type (
		// Request body
		request struct {
			Times int32 `json:"times"`
		}
		// Responses
		gachaDrawResult struct {
			CollectionID string `json:"collectionID"`
			Name         string `json:"name"`
			Rarity       int    `json:"rarity"`
			IsNew        bool   `json:"isNew"`
		}
		response struct {
			Results []*gachaDrawResult `json:"results"`
		}
	)

	requestBody := new(request)
	if err := c.Bind(requestBody); err != nil {
		return &myerror.InternalServerError{Err: err}
	}

	// timesのバリデーション
	times := requestBody.Times
	if !inBetween(int(times), 1, 10) {
		return &myerror.BadRequestError{Err: errors.New("Invalid times. times must be between 1 and 10.")}
	}

	user := dcontext.GetUserFromContext(c)
	if user == nil {
		return &myerror.UnauthorizedError{Err: errors.New("user not found")}
	}

	gachaDrawResults, err := gh.service.GachaDraw(c.Request().Context(), user.AuthToken, times)
	if err != nil {
		return &myerror.InternalServerError{Err: err}
	}

	res := make([]*gachaDrawResult, len(gachaDrawResults))
	for i, g := range gachaDrawResults {
		res[i] = &gachaDrawResult{
			CollectionID: g.CollectionID,
			Name:         g.Name,
			Rarity:       int(g.Rarity),
			IsNew:        g.IsNew,
		}
	}

	return c.JSON(http.StatusOK, &response{
		res,
	})
}

func inBetween(i, min, max int) bool {
	return i >= min && i <= max
}
