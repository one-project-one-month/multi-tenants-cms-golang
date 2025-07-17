package nats

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
)

var natsConn *nats.Conn

func InitNATS(url string) error {
	conn, err := nats.Connect(url)
	if err != nil {
		return err
	}
	natsConn = conn
	return nil
}

func Publish(subject string, data map[string]any) error {

	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return natsConn.Publish(subject, marshal)
}

func Subscribe(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	return natsConn.Subscribe(subject, handler)
}

func GetNATS() *nats.Conn {
	return natsConn
}
