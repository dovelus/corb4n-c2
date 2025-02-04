package comunication

import (
	"time"

	"github.com/rs/xid"
)

// GenerateID generates a random ID using only numbers and letters
func GenerateID() string {
	return xid.New().String()
}

func CurrentUnixTimestamp() int64 {
	return time.Now().Unix()
}
