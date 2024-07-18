package paho_mqtt

import (
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/mgw-device-manager/handler"
	"github.com/SENERGY-Platform/mgw-device-manager/util"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type Wrapper struct {
	client  mqtt.Client
	timeout time.Duration
}

func NewWrapper(client mqtt.Client, timeout time.Duration) *Wrapper {
	return &Wrapper{
		client:  client,
		timeout: timeout,
	}
}

func (w *Wrapper) Subscribe(topic string, qos byte, msgHandler func(m handler.Message)) error {
	if !w.client.IsConnectionOpen() {
		return util.NotConnectedErr
	}
	t := w.client.Subscribe(topic, qos, func(_ mqtt.Client, message mqtt.Message) {
		msgHandler(&msgWrapper{
			Message:   message,
			timestamp: time.Now(),
		})
	})
	if !t.WaitTimeout(w.timeout) {
		return util.OperationTimeoutErr
	}
	res := t.(*mqtt.SubscribeToken).Result()
	c, ok := res[topic]
	if !ok {
		return errors.New("no result")
	}
	if c < 0 || c > 2 {
		return fmt.Errorf("code=%d", c)
	}
	return t.Error()
}

func (w *Wrapper) Publish(topic string, qos byte, retained bool, payload any) error {
	if !w.client.IsConnectionOpen() {
		return util.NotConnectedErr
	}
	t := w.client.Publish(topic, qos, retained, payload)
	if !t.WaitTimeout(w.timeout) {
		return util.OperationTimeoutErr
	}
	return t.Error()
}

func (w *Wrapper) Connect() mqtt.Token {
	return w.client.Connect()
}

func (w *Wrapper) Disconnect(quiesce uint) {
	w.client.Disconnect(quiesce)
}

type msgWrapper struct {
	mqtt.Message
	timestamp time.Time
}
