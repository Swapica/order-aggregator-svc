package requests

import (
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type GetMatchRequest struct {
	MatchId int64
}

func NewGetMatch(r *http.Request) (GetMatchRequest, error) {
	var request GetMatchRequest

	matchId := chi.URLParam(r, "id")
	matchIdInt, err := strconv.ParseInt(matchId, 10, 64)
	if err != nil {
		return GetMatchRequest{}, errors.Wrap(err, "failed to parse match id to int")
	}
	request.MatchId = matchIdInt

	return request, nil
}
