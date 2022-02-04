package response

import (
	"fmt"
	"time"
)

type UserResponse struct {
	Id       int32    `json:"id"`
	NikeName string   `json:"name"`
	Birthday JsonTime `json:"birthday"`
	Gender   string   `json:"gender"`
	Mobile   string   `json:"mobile"`
}

type JsonTime time.Time

func (j JsonTime) MarshalJson() ([]byte, error) {
	var stmp = fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02"))
	return []byte(stmp), nil
}
