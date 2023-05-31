package user

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/noahhai/vigil/app/domains/token"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

type Service interface {
	Get(userID uint) (users []User, err error)
	GetByUsernameEmail(username, email string) (users []User, err error)
	Create(u *User) error
	Update(u *User) error
	Authenticate(u *User) (err error)
}

type userService struct {
	db *gorm.DB
	ts token.Service
}

func NewService(db *gorm.DB, ts token.Service) Service {
	return &userService{
		db: db,
		ts: ts,
	}
}

func (s *userService) All() []User {
	var users []User
	s.db.Find(&users)
	return users
}

func (s *userService) Get(userID uint) (users []User, err error) {
	s.db.First(&users, userID)
	if s.db.Error != nil {
		err = s.db.Error
		return
	}
	if len(users) < 1 {
		err = errors.New("user not found")
	}
	return
}

func (s *userService) GetByUsernameEmail(username, email string) (users []User, err error) {
	s.db.Where("email = ? OR username = ?", email, username).First(&users)
	if s.db.Error != nil {
		err = s.db.Error
		return
	}
	if len(users) < 1 {
		err = errors.New("user not found")
	}
	return
}

func (s *userService) Authenticate(u *User) (err error) {
	var users []User
	users, err = s.GetByUsernameEmail(u.Username, u.Username)
	if err != nil {
		log.Printf("login failed: %v\n", err)
		return errors.New("login failed")
	}
	userFound := users[0]
	if err := bcrypt.CompareHashAndPassword([]byte(userFound.PasswordHash), []byte(u.Password)); err != nil {
		log.Printf("login failed: %v\n", err)
		return errors.New("login failed")
	}
	u.Username = userFound.Username
	u.Email = userFound.Email
	return
}

func (s *userService) Create(u *User) (err error) {
	if u.ID != 0 {
		u.ID = 0
	}

	var byEmail []User
	if u.Email != "" {
		s.db.Where("email = ? OR username = ?", u.Email, u.Email).First(&byEmail)
		if len(byEmail) > 0 {
			return errors.New("Account already exists for email")
		}
	}
	var byUsername []User
	s.db.Where("email = ? OR username = ?", u.Username, u.Username).First(&byUsername)
	if len(byUsername) > 0 {
		return errors.New("Account already exists for username")
	}

	if u.Password == "" {
		return errors.New("must specify password")
	}
	if u.PhoneNumber != nil && *u.PhoneNumber != "" && !strings.HasPrefix(*u.PhoneNumber, "+") && !strings.HasPrefix(*u.PhoneNumber, "00") {
		pn := *u.PhoneNumber
		pn = "+1" + pn
		u.PhoneNumber = &pn
	}
	u.PasswordHash = s.hashAndSalt([]byte(u.Password))

	u.Token = s.ts.GenerateAgentToken()

	s.db.Create(u)
	if s.db.Error != nil {
		err = s.db.Error
	}
	return
}

func (s *userService) Update(u *User) (err error) {
	var currUser *User

	var byUsername []User
	byUsername, err = s.GetByUsernameEmail(u.Username, u.Email)
	if err != nil {
		return
	}
	if len(byUsername) > 0 {
		currUser = &byUsername[0]
	} else {
		return errors.New("user not found")
	}

	if u.Password != "" {
		u.PasswordHash = s.hashAndSalt([]byte(u.Password))
	}

	if u.Email != "" {
		if u.Email != currUser.Email && currUser.Username != u.Email {
			var byEmail []User
			s.db.Where("email = ? OR email = ?", u.Email, u.Username).First(&byEmail)
			if len(byEmail) > 0 {
				return errors.New("Account already exists for email")
			}
		}
		currUser.Email = u.Email
	}
	if u.PasswordHash != "" && u.PasswordHash != currUser.PasswordHash {
		currUser.PasswordHash = u.PasswordHash
	}

	if u.PhoneNumber != nil && *u.PhoneNumber != "" && !strings.HasPrefix(*u.PhoneNumber, "+") && !strings.HasPrefix(*u.PhoneNumber, "00") {
		pn := *u.PhoneNumber
		pn = "+1" + pn
		u.PhoneNumber = &pn
	}
	if u.PhoneNumber != nil {
		currUser.PhoneNumber = u.PhoneNumber
	}
	currUser.ResetCode = u.ResetCode
	if u.Token != "" {
		currUser.Token = u.Token
	}
	currUser.NotificationEmail = u.NotificationEmail
	currUser.NotificationMobile = u.NotificationMobile
	currUser.NotificationPhone = u.NotificationPhone
	s.db.Save(currUser)

	if s.db.Error != nil {
		err = s.db.Error
	}
	return
}

func (s *userService) hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}
