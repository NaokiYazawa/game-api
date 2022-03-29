package ranking

import (
	"log"
	"net/http"
	"strconv"

	rs "game-api/pkg/domain/service/ranking"
	"game-api/pkg/interfaces/api/myerror"

	"github.com/labstack/echo"
)

// 本アプリは「get_ranking_list」の1種類
type Handler interface {
	HandleGet(c echo.Context) error
}

type handler struct {
	useCase rs.Service
}

func NewHandler(rankingUseCase rs.Service) Handler {
	return &handler{
		useCase: rankingUseCase,
	}
}

func (rh handler) HandleGet(c echo.Context) error {
	type ranking struct {
		UserID   string `json:"userId"`
		UserName string `json:"userName"`
		Rank     int32  `json:"rank"`
		Score    int32  `json:"score"`
	}

	type response struct {
		Ranks []*ranking `json:"ranks"`
	}

	start, err := strconv.Atoi(c.FormValue("start"))
	if err != nil {
		log.Printf("strconv.Atoi is failed : %v", err)
		return err
	}

	if start <= 0 {
		return &myerror.InternalServerError{Err: err}
	}

	rankingList, err := rh.useCase.SelectRankingList(int64(start))
	if err != nil {
		log.Println("start must be positive")
		return err
	}

	res := make([]*ranking, len(rankingList))
	for i, rank := range rankingList {
		res[i] = &ranking{
			UserID:   rank.UserID,
			UserName: rank.UserName,
			Rank:     rank.Rank,
			Score:    rank.Score,
		}
	}

	return c.JSON(http.StatusOK, &response{
		res,
	})
}
