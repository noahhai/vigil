package user

import (
	"encoding/json"
	"github.com/noahhai/vigil/app/domains/token"
	"github.com/noahhai/vigil/app/types/web"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/noahhai/vigil/app/utils"
)

type handler struct {
	us Service
	ts token.Service
}

type Handler interface {
	HandleUpdate(w http.ResponseWriter, r *http.Request)
	HandleCreate(w http.ResponseWriter, r *http.Request)
	HandleLogin(w http.ResponseWriter, r *http.Request)
	HandleNotification(w http.ResponseWriter, r *http.Request)
	HandleGetByUsernameEmail(w http.ResponseWriter, r *http.Request)
	HandleGetByID(w http.ResponseWriter, r *http.Request)
}

func NewHandler(s Service, t token.Service) Handler {
	return &handler{
		us: s,
		ts: t,
	}
}

func (h *handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var webUser web.User
	err := decoder.Decode(&webUser)
	if err != nil {
		log.Printf("Error handlng upsert: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u := User{
		Email: webUser.Email,
		Username: webUser.Username,
		Password: webUser.Password,
		PhoneNumber: webUser.PhoneNumber,
		Token: webUser.Token,
		NotificationPhone: webUser.NotificationPhone,
		NotificationEmail: webUser.NotificationEmail,
		NotificationMobile: webUser.NotificationMobile,
	}
	err = h.us.Update(&u)
	webUser = web.User{
			Email: u.Email,
			Username: u.Username,
			PhoneNumber: u.PhoneNumber,
			Token: u.Token,
			NotificationPhone: u.NotificationPhone,
			NotificationEmail: u.NotificationEmail,
			NotificationMobile: u.NotificationMobile,
		}
	utils.WriteDataResponse(w, webUser, err)
}

func (h *handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var webUser web.User
	err := decoder.Decode(&webUser)
	if err != nil {
		log.Printf("Error decoding create request: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u := User{
		Email: webUser.Email,
		Username: webUser.Username,
		Password: webUser.Password,
		PhoneNumber: webUser.PhoneNumber,
		Token: webUser.Token,
		NotificationPhone: webUser.NotificationPhone,
		NotificationEmail: webUser.NotificationEmail,
		NotificationMobile: webUser.NotificationMobile,
	}
	err = h.us.Create(&u)
	webUser = web.User{
			Email: u.Email,
			Username: u.Username,
			PhoneNumber: u.PhoneNumber,
			Token: u.Token,
			NotificationPhone: u.NotificationPhone,
			NotificationEmail: u.NotificationEmail,
			NotificationMobile: u.NotificationMobile,
		}
	if err != nil {
		log.Printf("Error handlng create: %v\n", err)
		utils.WriteDataResponse(w, webUser, err)
		return
	}

	tk, err := h.ts.Generate(u.Username, u.Email)
	webUser.Token = tk

	utils.WriteDataResponse(w, webUser, err)
}

func (h *handler) HandleGetByUsernameEmail(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	email := mux.Vars(r)["email"]
	users, err := h.us.GetByUsernameEmail(username, email)
	var userModels []web.User
	for _, u := range users {
		userModels = append(userModels, web.User{
			Email: u.Email,
			Username: u.Username,
			PhoneNumber: u.PhoneNumber,
			Token: u.Token,
			NotificationPhone: u.NotificationPhone,
			NotificationEmail: u.NotificationEmail,
			NotificationMobile: u.NotificationMobile,
		})
	}
	utils.WriteDataResponse(w, userModels, err)
}

func (h *handler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	userIDstr := mux.Vars(r)["id"]
	var userID uint64
	var err error
	userID, err = strconv.ParseUint(userIDstr, 10, 32)
	if err != nil {
		log.Printf("invalid user id: '%d'\n", userID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	users, err := h.us.Get(uint(userID))
	var userModels []web.User
	for _, u := range users {
		userModels = append(userModels, web.User{
			Email: u.Email,
			PhoneNumber: u.PhoneNumber,
			Token: u.Token,
			NotificationPhone: u.NotificationPhone,
			NotificationEmail: u.NotificationEmail,
			NotificationMobile: u.NotificationMobile,
		})
	}
	utils.WriteDataResponse(w, userModels, err)
}

func (h *handler) HandleNotification(w http.ResponseWriter, r *http.Request) {

}
func (h *handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var wu web.UserAuth
	err := decoder.Decode(&wu)
	if err != nil {
		log.Printf("Error decoding request: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u := User{
		Username:wu.Username,
		Password:wu.Password,
	}

	if err := h.us.Authenticate(&u); err != nil {
		utils.WriteDataResponseWithStatus(w, nil, err, http.StatusForbidden)
		return
	}


	tk, err := h.ts.Generate(u.Username, u.Email)
	token := web.Token{
		Token: tk,
		Username: u.Username,
		Email: u.Email,
	}
	utils.WriteDataResponse(w, token, err)
}
