package nats

import (
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

func Publish(subject string, data []byte) error {
	return natsConn.Publish(subject, data)
}

func Subscribe(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	return natsConn.Subscribe(subject, handler)
}

func GetNATS() *nats.Conn {
	return natsConn
}
