package email

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/noahhai/vigil/app/types"
)


const (
	Sender = "notification@getvigil.io"
	Subject = "Task Finished"
	CharSet = "UTF-8"
)

type sesService struct {
	ses *ses.SES
}

func newSesService() (svc Service, err error) {
	if sess, e := session.NewSession(&aws.Config{
		Region:aws.String("us-east-1"),
	},
	); e != nil {
		err = e
		return
	} else {
		ses := ses.New(sess)
		s := sesService{ses:ses}
		svc = &s
	}
	return
}

func (s *sesService) SendEmail(recipient string, w types.Work) error {
	textBody := fmt.Sprintf("Task `%s` completed in %s.", w.Name, w.Duration)
	if w.Status != "" {
		textBody += fmt.Sprintf(" Status: %s.", w.Status)
	}

	// TODO : change htmlBody to better format
	htmlBody := textBody

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{
			},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(textBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(Sender),
	}

	// Attempt to send the email.
	result, err := s.ses.SendEmail(input)

	// Display error messages if they occur.
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
	fmt.Println("Email Sent to address: " + recipient + ". Result: " + result.String())

	return nil
}