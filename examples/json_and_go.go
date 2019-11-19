package examples

import "github.com/philandstuff/dhall-golang"

type Message struct {
	Name string
	Body string
	Time int64
}

var dhallMessage = `{Name = "Alice",Body = "Hello", Time = 1294706395881547000}`

func DecodeMessage() (*Message, error) {
	var m Message
	err := dhall.Unmarshal([]byte(dhallMessage), &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
