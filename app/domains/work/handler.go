package work

import (
	"encoding/json"
	"errors"
	"app/domains/user"
	"github.com/noahhai/vigil/app/types"
	"github.com/noahhai/vigil/app/types/web"
	"github.com/noahhai/vigil/app/utils"
	"log"
	"net/http"
)

type handler struct {
	ws Service
	us user.Service
}

type Handler interface {
	HandlePost(w http.ResponseWriter, r *http.Request)
}

func NewHandler(ws Service, us user.Service) Handler {
	return &handler{
		ws: ws,
		us: us,
	}
}

func (h *handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var am web.AgentMessage
	err := decoder.Decode(&am)
	if err != nil {
		log.Printf("Error handlng post: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username, token := am.Username, am.Token
	if username == "" {
		err = errors.New("username cannot be empty")
	}
	if token == "" {
		err = errors.New("token cannot be empty")
	}
	users, e := h.us.GetByUsernameEmail(username, username)
	if e != nil {
		err = e
	} else if len(users) < 1 {
		err = errors.New("user not found")
	} else if users[0].Token != token {
		err = errors.New("invalid auth")
	}
	if err != nil {
		utils.WriteDataResponseWithStatus(w, nil, err, http.StatusUnauthorized)
		return
	}

	work := types.Work{
		Name:     am.Name,
		Duration: am.Duration,
		Status:   am.Status,
	}
	err = h.ws.HandleWorkDone(users[0], work)
	utils.WriteDataResponse(w, nil, err)
}