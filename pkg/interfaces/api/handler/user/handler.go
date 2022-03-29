package user

import (
	"net/http"

	us "game-api/pkg/domain/service/user"
	"game-api/pkg/interfaces/api/dcontext"
	"game-api/pkg/interfaces/api/myerror"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// Handler UserにおけるHandlerのインターフェース
// 本アプリは「post_user_create」・「get_user_get」・「post_user_update」の3種類
type Handler interface {
	HandleCreate(c echo.Context) error
	HandleGet(c echo.Context) error
	HandleUpdate(c echo.Context) error
}

type handler struct {
	service us.Service
}

// NewHandler Userデータに関するHandlerを生成
func NewHandler(userService us.Service) Handler {
	return &handler{
		service: userService,
	}
}

// HandleCreate ユーザを作成するHandler
func (uh handler) HandleCreate(c echo.Context) error {
	type (
		// Request body
		request struct {
			Name string `json:"name"`
		}
		// Responses
		response struct {
			Token string `json:"token"`
		}
	)

	requestBody := new(request)
	if err := c.Bind(requestBody); err != nil {
		return &myerror.InternalServerError{Err: err}
	}

	authToken, err := uh.service.Create(requestBody.Name)
	if err != nil {
		return &myerror.InternalServerError{Err: err}
	}

	return c.JSON(http.StatusOK, &response{Token: authToken})
}

// HandleGet ユーザー取得処理
func (uh handler) HandleGet(c echo.Context) error {
	type response struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		HighScore int32  `json:"highScore"`
		Coin      int32  `json:"coin"`
	}

	user := dcontext.GetUserFromContext(c)
	if user == nil {
		return &myerror.UnauthorizedError{Err: errors.New("user not found")}
	}

	return c.JSON(http.StatusOK, &response{
		ID:        user.ID,
		Name:      user.Name,
		HighScore: user.HighScore,
		Coin:      user.Coin,
	})
}

// HandleUpdate ユーザー更新処理
func (uh handler) HandleUpdate(c echo.Context) error {
	type (
		request struct {
			Name string `json:"name"`
		}
		response struct{}
	)

	requestBody := new(request)
	if err := c.Bind(requestBody); err != nil {
		return &myerror.InternalServerError{Err: err}
	}

	user := dcontext.GetUserFromContext(c)
	if user == nil {
		return &myerror.UnauthorizedError{Err: errors.New("user not found")}
	}

	if err := uh.service.UpdateName(user, requestBody.Name); err != nil {
		return &myerror.InternalServerError{Err: err}
	}

	return c.JSON(http.StatusOK, &response{})
}
