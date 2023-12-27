package email

import (
	"log"
	"net/smtp"

	"github.com/valli0x/booking-server/pkg/models"
)

type MailerSMTP struct {
	smtpServer, smtpPort string
	from, password       string
}

func NewMailerSMPT(smtpServer, smtpPort, from, password string) *MailerSMTP {
	return &MailerSMTP{
		smtpServer: smtpServer,
		smtpPort:   smtpPort,
		from:       from,
		password:   password,
	}
}

func (m *MailerSMTP) SendConfirmation(email string, o models.Order) error {
	to := email
	msg := []byte("The order has been accepted")

	auth := smtp.PlainAuth("", m.from, m.password, m.smtpServer)

	err := smtp.SendMail(m.smtpServer+":"+m.smtpPort, auth, m.from, []string{to}, msg)
	if err != nil {
		return err
	}

	log.Printf("Sending confirmation to %s for order %v", email, o)
	return nil
}
