package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Phần đại diện cho Payload trong 1 token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  int64     `json:"issued_at"`
	ExpiredAt int64     `json:"expired_at"`
}

// Tạo mới 1 Payload
func NewPayLoad(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now().Unix(),
		ExpiredAt: time.Now().Add(duration).Unix(),
	}
	return payload, nil
}

func (p *Payload) Valid() error {
	now := time.Now().Unix()
	if now > p.ExpiredAt {
		return errors.New("Token is expired")
	}
	return nil
}
