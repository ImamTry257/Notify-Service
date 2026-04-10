package nats

import (
	"log"

	"github.com/nats-io/nats.go"
	"github.com/ImamTry257/Notify-Service/config"
)

func NewNATSConnection(cfg config.NATSConfig) (nats.JetStreamContext, *nats.Conn, error) {
	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, nc, err
	}

	// Ensure stream exists
	_, err = js.StreamInfo(cfg.StreamName)
	if err != nil {
		log.Printf("Stream %s not found, attempting to create...", cfg.StreamName)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     cfg.StreamName,
			Subjects: []string{cfg.Subject},
		})
		if err != nil {
			return nil, nc, err
		}
	}

	return js, nc, nil
}
