package resources

import (
	"encoding/json"
	"time"
)

type RequestResource struct {
	Request json.RawMessage `json:"request"`
	Headers json.RawMessage `json:"headers"`
	Created time.Time       `json:"created"`
}
