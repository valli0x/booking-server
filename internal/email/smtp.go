package email

import (
	"log"
	"net/smtp"

	"github.com/valli0x/booking-server/pkg/models"
)

type MailerSMTP struct {
	from, password string
}

func NewMailerSMPT(from, password string) *MailerSMTP {
	return &MailerSMTP{
		from:     from,
		password: password,
	}
}

func (m *MailerSMTP) SendConfirmation(email string, o models.Order) error {
	to := email
	msg := []byte("The order has been accepter")

	smtpServer := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", m.from, m.password, smtpServer)

	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, m.from, []string{to}, msg)
	if err != nil {
		return err
	}

	log.Printf("Sending confirmation to %s for order %v", email, o)
	return nil
}
