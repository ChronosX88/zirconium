package models

// BaseMessage is a basic message model, basis of the whole protocol. It is used for a very easy protocol extension process.
type BaseMessage struct {
	ID          string                 `json:"id"`
	MessageType string                 `json:"type"`
	From        string                 `json:"from,omitempty"`
	To          string                 `json:"to,omitempty"`
	Ok          bool                   `json:"ok"`
	Payload     map[string]interface{} `json:"payload"`
}

func NewBaseMessage(id, messageType, from string, to string, ok bool, payload map[string]interface{}) BaseMessage {
	return BaseMessage{
		ID:          id,
		MessageType: messageType,
		From:        from,
		To:          to,
		Ok:          ok,
		Payload:     payload,
	}
}

func NewEmptyBaseMessage() BaseMessage {
	return BaseMessage{
		ID:          "",
		MessageType: "",
		From:        "",
		To:          "",
		Ok:          true,
		Payload:     map[string]interface{}{},
	}
}
