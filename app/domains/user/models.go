package user

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email         string
	Username      string
	PasswordHash  string
	Password      string `gorm:"-"`
	PhoneNumber   *string
	Token         string
	LoginFailures string
	ResetCode     string

	NotificationPhone  bool
	NotificationEmail  bool
	NotificationMobile bool
}
