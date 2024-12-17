package mail

import "context"

type FakeSender struct {
	err  error
	Sent []Email
}

func (s *FakeSender) Send(_ context.Context, email Email) error {
	if s.err != nil {
		return s.err
	}
	s.Sent = append(s.Sent, email)
	return nil
}
