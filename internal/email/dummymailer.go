package email

import (
	"log"

	"github.com/valli0x/booking-server/pkg/models"
)

type Mailer interface {
	SendConfirmation(email string, o models.Order) error
}

type DummyMailer struct{}

func NewDummyMailer() *DummyMailer {
	return &DummyMailer{}
}

func (m *DummyMailer) SendConfirmation(email string, o models.Order) error {
	log.Printf("Sending confirmation to %s for order %v", email, o)
	return nil
}
