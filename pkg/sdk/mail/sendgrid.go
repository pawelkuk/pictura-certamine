package mail

import (
	"context"
	"fmt"
	"os"

	"github.com/pawelkuk/pictura-certamine/pkg/sdk/logger"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendgridSender struct {
	client *sendgrid.Client
	log    *logger.Logger
}

func (s *SendgridSender) Send(ctx context.Context, email Email) error {
	from := mail.NewEmail(email.From.Name, email.From.Address)
	to := mail.NewEmail(email.To.Name, email.To.Address)
	message := mail.NewSingleEmail(from, email.Subject, to, email.Content, email.HTMLContent)
	response, err := s.client.Send(message)
	if err != nil {
		return fmt.Errorf("could not send mail: %w", err)
	}
	if response.StatusCode >= 400 {
		s.log.Error(ctx, "could not send mail", "headers", response.Headers, "body", response.Body)
		return fmt.Errorf("returned non 2xx status code %d", response.StatusCode)
	}
	return nil
}

func NewSendgridSender(apiKey string) *SendgridSender {
	return &SendgridSender{
		client: sendgrid.NewSendClient(apiKey),
		log:    logger.New(os.Stdout, logger.LevelInfo, "MAIL"),
	}
}
