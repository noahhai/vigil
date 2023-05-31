package sms

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/noahhai/vigil/app/types"
)

type snsService struct {
	sns *sns.SNS
}

func newSnsService() (svc Service, err error) {
	if sess, e := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	},
	); e != nil {
		err = e
		return
	} else {
		sns := sns.New(sess)
		s := snsService{sns: sns}
		svc = &s
	}
	return
}

func (s *snsService) SendSMS(recipient string, w types.Work) error {
	textBody := fmt.Sprintf("Task `%s` completed in %s.", w.Name, w.Duration)
	if w.Status != "" {
		textBody += fmt.Sprintf(" Status: %s.", w.Status)
	}

	params := &sns.PublishInput{
		Message:     aws.String(textBody),
		PhoneNumber: aws.String(recipient),
	}
	result, err := s.sns.Publish(params)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
			return aerr
		} else {
			fmt.Println(err.Error())
			return err
		}
	}
	fmt.Println("SMS Sent to recipient: " + recipient + ". Result: " + result.String())

	return nil
}
