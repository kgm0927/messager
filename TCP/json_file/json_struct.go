package jsonfile

import (
	"encoding/json"
)

type Message struct {
	Room_name   string `json:"Room"`
	Client_name string `json:"client"` // 클라이언트 이름
	Args        string `json:"talk"`   //
	Ip          string `json:"ip"`
	Id          int    `json:"id"` // 실행 명령
}

type Making_message interface {
	Initialize(sentences ...string)
	Serialize() (error, []byte)
	UnSerialize(B []byte) (error, *Message)
	Return_self() *Message
}

func (M *Message) Initialize(sentences ...any) {
	M.Room_name = sentences[0].(string)
	M.Client_name = sentences[1].(string)
	M.Args = sentences[2].(string)
	M.Ip = sentences[3].(string)
	M.Id = sentences[4].(int)

}

func (M *Message) Serialize() ([]byte, error) {
	b, err := json.Marshal(M)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (M *Message) UnSerialize(B []byte) (*Message, error) {

	err := json.Unmarshal(B, M)

	if err != nil {
		return nil, err
	}
	return M, nil
}

func (M *Message) Return_self() *Message {
	return M
}
