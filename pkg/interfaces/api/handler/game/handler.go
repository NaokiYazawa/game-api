package game

import (
	"net/http"

	gs "game-api/pkg/domain/service/game"
	"game-api/pkg/interfaces/api/dcontext"
	"game-api/pkg/interfaces/api/myerror"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// Handler UserにおけるHandlerのインターフェース
// 本アプリは「get_collection_list」の1種類
type Handler interface {
	HandleFinish(c echo.Context) error
}

type handler struct {
	service gs.Service
}

func NewHandler(collectionItemService gs.Service) Handler {
	return &handler{
		service: collectionItemService,
	}
}

func (gh handler) HandleFinish(c echo.Context) error {
	type (
		// Request body
		request struct {
			Score int32 `json:"score"`
		}
		// Responses
		response struct {
			Coin int32 `json:"coin"`
		}
	)

	requestBody := new(request)
	if err := c.Bind(requestBody); err != nil {
		return &myerror.InternalServerError{Err: err}
	}

	score := requestBody.Score

	if score < 0 {
		return &myerror.BadRequestError{Err: errors.New("score is not positive")}
	}

	user := dcontext.GetUserFromContext(c)
	if user == nil {
		return &myerror.UnauthorizedError{Err: errors.New("user not found")}
	}

	coin, err := gh.service.GameFinish(user.AuthToken, score)
	if err != nil {
		return &myerror.InternalServerError{Err: err}
	}

	return c.JSON(http.StatusOK, &response{
		Coin: coin,
	})
}
