package sms

import (
	"fmt"
	"github.com/noahhai/vigil/app/consts"
	"github.com/noahhai/vigil/app/types"
)

type Service interface {
	SendSMS(recipient string, w types.Work) error
}

func NewService(cloudType consts.CloudType) (svc Service, err error) {
	if cloudType == consts.CloudTypeAWS {
		svc, err = newSnsService()
	} else {
		err = fmt.Errorf("Service not found for cloud type %s", cloudType)
	}
	return
}