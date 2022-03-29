package middleware

import (
	us "game-api/pkg/domain/service/user"
	"game-api/pkg/interfaces/api/dcontext"
	"game-api/pkg/interfaces/api/myerror"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// Middleware middlewareのインターフェース
type Middleware interface {
	Authenticate(echo.HandlerFunc) echo.HandlerFunc
}

type middleware struct {
	service us.Service
}

// NewMiddleware userUseCaseと疎通
func NewMiddleware(service us.Service) Middleware {
	return &middleware{
		service: service,
	}
}

// Authenticate ユーザ認証を行ってContextへユーザID情報を保存する
func (m middleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// リクエストヘッダからx-token(認証トークン)を取得
		token := c.Request().Header.Get("x-token")
		if token == "" {
			return &myerror.UnauthorizedError{Err: errors.New("x-token is empty")}
		}
		// データベースから認証トークンに紐づくユーザの情報を取得
		user, err := m.service.SelectByAuthToken(token)
		if err != nil {
			return &myerror.InternalServerError{Err: err}
		}
		if user == nil {
			return &myerror.UnauthorizedError{Err: errors.Errorf(`user is not found: token="%s"`, token)}
		}
		dcontext.SetUser(c, *user)
		return next(c)
	}
}
