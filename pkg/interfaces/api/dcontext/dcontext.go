package dcontext

import (
	um "game-api/pkg/domain/model/user"

	"github.com/labstack/echo"
)

var userKey = "userKey"

// SetUser Contextへユーザを保存する
// ユーザIDではなく、ユーザを保存することで、Contextからの取得が容易になる
func SetUser(c echo.Context, user um.User) {
	c.Set(userKey, user)
}

// GetUserFromContext Contextからユーザを取得する
func GetUserFromContext(c echo.Context) *um.User {
	var user um.User
	if c.Get(userKey) == nil {
		return nil
	}

	user = c.Get(userKey).(um.User)
	return &user
}
