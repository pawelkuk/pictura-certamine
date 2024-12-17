package mail

import "net/mail"

type Email struct {
	Subject     string
	From        mail.Address
	To          mail.Address
	Content     string
	HTMLContent string
}
