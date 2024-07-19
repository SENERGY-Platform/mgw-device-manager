package msg_relay_hdl

import (
	"github.com/SENERGY-Platform/mgw-device-manager/handler"
	"reflect"
	"testing"
	"time"
)

type mockMessage struct {
	topic     string
	payload   []byte
	timestamp time.Time
}

func (m *mockMessage) Topic() string {
	return m.topic
}

func (m *mockMessage) Payload() []byte {
	return m.payload
}

func TestHandler(t *testing.T) {
	msg := &mockMessage{
		topic:     "test",
		payload:   []byte("test"),
		timestamp: time.Now(),
	}
	testMsgHdl := func(m handler.Message) {
		if !reflect.DeepEqual(m, msg) {
			t.Error("expected", msg, "got", m)
		}
	}
	h := New(1, testMsgHdl)
	err := h.Put(msg)
	if err != nil {
		t.Error(err)
	}
	if len(h.messages) != 1 {
		t.Error("message not in channel")
	}
	h.Start()
	time.Sleep(1 * time.Second)
	if len(h.messages) > 0 {
		t.Error("message not consumed")
	}
	h.Stop()
	t.Run("buffer full", func(t *testing.T) {
		testMsgHdl = func(m handler.Message) {}
		h = New(1, testMsgHdl)
		err = h.Put(msg)
		if err != nil {
			t.Error(err)
		}
		err = h.Put(msg)
		if err == nil {
			t.Error(err)
		}
		if len(h.messages) != 1 {
			t.Error("message not in channel")
		}
		h.Start()
		time.Sleep(1 * time.Second)
		if len(h.messages) > 0 {
			t.Error("message not consumed")
		}
		h.Stop()
	})
}
