package register

import (
	"encoding/json"
	"time"

	"github.com/vranystepan/email/pkg/messages/message"
)

// Payload holds registration data
type Payload struct {
	message.Payload
	TimeCreated time.Time `json:"timeCreated"`
	Source      string    `json:"source"`
	Email       string    `json:"email"`
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
