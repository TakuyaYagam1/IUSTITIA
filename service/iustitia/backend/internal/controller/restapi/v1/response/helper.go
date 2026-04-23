package response

type Message struct {
	Message string `json:"message"`
}

func Msg(text string) Message { return Message{Message: text} }
