package json

import "encoding/json" //nolint:depguard // we need this for aliasing

type (
	RawMessage = json.RawMessage
	Number     = json.Number
)
