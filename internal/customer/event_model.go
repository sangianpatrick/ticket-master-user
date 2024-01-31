package customer

import "time"

type RegisterEvent struct {
	Customer     Customer `json:"customer"`
	Verification struct {
		URL       string    `json:"url"`
		ExpiredAt time.Time `json:"expiredAt"`
	}
}
