package server

import (
	"fmt"
	"log"
	"net/http"

	cis "game-api/pkg/domain/service/collection_item"
	cs "game-api/pkg/domain/service/component"
	gcs "game-api/pkg/domain/service/gacha"
	gs "game-api/pkg/domain/service/game"
	rs "game-api/pkg/domain/service/ranking"
	us "game-api/pkg/domain/service/user"
	"game-api/pkg/infrastructure/mysql"

	cir "game-api/pkg/infrastructure/mysql/repositoryimpl/collection_item"
	gpr "game-api/pkg/infrastructure/mysql/repositoryimpl/gacha_probability"
	txr "game-api/pkg/infrastructure/mysql/repositoryimpl/transaction"
	ur "game-api/pkg/infrastructure/mysql/repositoryimpl/user"
	ucir "game-api/pkg/infrastructure/mysql/repositoryimpl/user_collection_item"
	"game-api/pkg/infrastructure/redis"
	rr "game-api/pkg/infrastructure/redis/repositoryimpl/ranking"
	"game-api/pkg/interfaces/api/dcontext"
	cih "game-api/pkg/interfaces/api/handler/collection_item"
	gch "game-api/pkg/interfaces/api/handler/gacha"
	gh "game-api/pkg/interfaces/api/handler/game"
	rh "game-api/pkg/interfaces/api/handler/ranking"
	uh "game-api/pkg/interfaces/api/handler/user"
	authMiddleware "game-api/pkg/interfaces/api/middleware"
	"game-api/pkg/interfaces/api/myerror"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Serve サーバー起動処理
func Serve(addr string) {
	// 依存性の注入（tx）
	transactionRepoImpl := txr.NewRepositoryImpl(mysql.Tx(mysql.NewSQLHandler()))

	// 依存性の注入（user）
	userRepoImpl := ur.NewRepositoryImpl(mysql.NewSQLHandler())
	userService := us.NewService(userRepoImpl)
	userHandler := uh.NewHandler(userService)
	auth := authMiddleware.NewMiddleware(userService)

	// 依存性の注入（ranking）
	rankingRepoImpl := rr.NewRepositoryImpl(redis.NewCacheHandler())
	rankingService := rs.NewService(rankingRepoImpl)
	rankingHandler := rh.NewHandler(rankingService)

	// 依存性の注入（collection_item）
	collectionItemRepoImpl := cir.NewRepositoryImpl(mysql.NewSQLHandler())
	userCollectionItemRepoImpl := ucir.NewRepositoryImpl(mysql.NewSQLHandler())
	collectionItemService := cis.NewService(collectionItemRepoImpl, userCollectionItemRepoImpl)
	collectionItemHandler := cih.NewHandler(collectionItemService)

	// 依存性の注入(game)
	gameService := gs.NewService(userRepoImpl, rankingRepoImpl)
	gameHandler := gh.NewHandler(gameService)

	// 依存性の注入（gacha）
	facade := cs.NewFacade()
	gachaRepoImpl := gpr.NewRepositoryImpl(mysql.NewSQLHandler())
	gachaService := gcs.NewService(userRepoImpl, collectionItemRepoImpl, userCollectionItemRepoImpl, gachaRepoImpl, transactionRepoImpl, facade)
	gachaHandler := gch.NewHandler(gachaService)

	echo.NotFoundHandler = func(c echo.Context) error {
		return &myerror.NotFoundError{Err: fmt.Errorf(`URL is invalid: (url="%s")`, c.Request().URL)}
	}
	e := echo.New()
	e.Use(
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{"Content-Type", "Accept", "Origin", "x-token"},
		}),
		middleware.Logger(),
		middleware.Recover(),
	)
	e.HTTPErrorHandler = errorHandler

	e.POST("/user/create", userHandler.HandleCreate)
	e.GET("/user/get", auth.Authenticate(userHandler.HandleGet))
	e.PATCH("/user/update", auth.Authenticate(userHandler.HandleUpdate))
	e.GET("collection/list", auth.Authenticate(collectionItemHandler.HandleGet))
	e.GET("ranking/list", auth.Authenticate(rankingHandler.HandleGet))
	e.POST("/game/finish", auth.Authenticate(gameHandler.HandleFinish))
	e.POST("/gacha/draw", auth.Authenticate(gachaHandler.HandleGacha))

	log.Println("Server running...")
	if err := e.Start(addr); err != nil {
		log.Fatalf("Listen and serve failed. %+v", err)
	}
}

func errorHandler(err error, c echo.Context) {
	type response struct {
		Message string `json:"message"`
	}
	var (
		userID  string
		code    int
		msg     string
		errInfo error
	)

	if user := dcontext.GetUserFromContext(c); user != nil {
		userID = user.ID
	}

	switch e := err.(type) {
	case *myerror.BadRequestError:
		code = http.StatusBadRequest
		msg = e.Message
		errInfo = e.Err
	case *myerror.UnauthorizedError:
		code = http.StatusUnauthorized
		msg = "401 Authentication error"
		errInfo = e.Err
	case *myerror.NotFoundError:
		code = http.StatusNotFound
		msg = "404 not found"
		errInfo = e.Err
	case *myerror.InternalServerError:
		code = http.StatusInternalServerError
		msg = "InternalServerError"
		errInfo = e.Err
	default:
		code = http.StatusInternalServerError
		msg = "ServerError"
		errInfo = err
	}

	log.Printf(`access:"%s", userID:"%s", errorCode:%d, errorMessage:"%s", error="%+v"`, c.Request().URL, userID, code, msg, errInfo)

	if !c.Response().Committed {
		if err := c.JSON(code, &response{
			Message: msg,
		}); err != nil {
			log.Print("Json conversion error in errorResponse")
		}
	}
}
