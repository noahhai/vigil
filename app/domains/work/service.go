package work

import (
	"fmt"
	"github.com/noahhai/vigil/app/domains/email"
	"github.com/noahhai/vigil/app/domains/sms"
	"github.com/noahhai/vigil/app/domains/user"
	"github.com/noahhai/vigil/app/types"
)



type Service interface {
	HandleWorkDone(u user.User, work types.Work) (err error)
}

type workService struct {
	email email.Service
	sms sms.Service
}

func NewWorkService(email email.Service, sms sms.Service) Service {
	return &workService{email: email, sms: sms}
}

func (s *workService) HandleWorkDone(u user.User, work types.Work) (errAny error) {
	if u.NotificationEmail {
		if u.Email == "" {
			errAny = fmt.Errorf("email not set for user '%s'", u.Username)
		} else if err := s.email.SendEmail(u.Email, work); err != nil {
			errAny = err
		}
	}
	if u.NotificationPhone {
		if u.PhoneNumber == nil || *u.PhoneNumber == "" {
			errAny = fmt.Errorf("phone number not set for user '%s'", u.Username)
		} else if  err := s.sms.SendSMS(*u.PhoneNumber, work); err != nil {
			errAny = err
		}
	}
	return
}
