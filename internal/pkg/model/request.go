package model

// Request ...
type Request struct {
	ID        int32  `json:"id,omitempty"`
	GroupUUID string `json:"group_uuid,omitempty"`
	UUID      string `json:"uuid,omitempty"`
	Body      []byte `json:"body,omitempty"`
}
