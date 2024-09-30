package jsonfile

import (
	"encoding/json"
)

type Message struct {
	Room_name   string `json:"Room"`
	Client_name string `json:"client"`
	Talk        string `json:"talk"`
}

type Making_message interface {
	Initialize(sentences ...string)
	Serialize() (error, []byte)
	UnSerialize(B []byte) (error, *Message)
	Return_self() *Message
}

func (M *Message) Initialize(sentences ...string) {
	M.Room_name = sentences[0]
	M.Client_name = sentences[1]
	M.Talk = sentences[2]
}

func (M *Message) Serialize() (error, []byte) {
	b, err := json.Marshal(M)
	if err != nil {
		return err, nil
	}
	return nil, b
}

func (M *Message) UnSerialize(B []byte) (error, *Message) {

	err := json.Unmarshal(B, M)

	if err != nil {
		return err, nil
	}
	return nil, M
}

func (M *Message) Return_self() *Message {
	return M
}
