package token

import "time"

//Maker is a interface manage tokens
type Maker interface {
	//Tạo token dựa trên username, và thời gian tạo
	CreateToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
