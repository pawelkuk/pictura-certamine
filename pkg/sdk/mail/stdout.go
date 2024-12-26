package mail

import (
	"context"
	"os"

	"github.com/pawelkuk/pictura-certamine/pkg/sdk/logger"
)

type StdoutSender struct {
	log *logger.Logger
}

func (s *StdoutSender) Send(ctx context.Context, email Email) error {
	s.log.Info(ctx, "sending email", "from", email.From, "to", email.To, "subject", email.Subject, "content", email.Content, "html_content", email.HTMLContent)
	return nil
}

func NewStdoutSender() *StdoutSender {
	return &StdoutSender{
		log: logger.New(os.Stdout, logger.LevelInfo, "MAIL"),
	}
}
