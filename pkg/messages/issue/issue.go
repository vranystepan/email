package issue

import (
	"encoding/json"
	"time"

	"github.com/vranystepan/email/pkg/messages/message"
)

// New creates a new Payload struct
func New(email string) Payload {
	return Payload{
		Email:       email,
		TimeCreated: time.Now(),
	}
}

// Payload holds registration data
type Payload struct {
	message.Payload
	Email       string    `json:"email"`
	TimeCreated time.Time `json:"timeCreated"`
}

// Parse converts JSON paylod to Payload struct
func Parse(message string) (Payload, error) {
	var m Payload
	err := json.Unmarshal([]byte(message), &m)
	if err != nil {
		return m, err
	}
	return m, nil
}
