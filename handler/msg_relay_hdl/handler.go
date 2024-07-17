package msg_relay_hdl

import (
	"errors"
	"github.com/SENERGY-Platform/mgw-device-manager/handler"
	"github.com/SENERGY-Platform/mgw-device-manager/util"
)

const logPrefix = "[relay-hdl]"

type Handler struct {
	messages   chan handler.Message
	handleFunc handler.MessageHandler
	dChan      chan struct{}
}

func New(buffer int, handleFunc handler.MessageHandler) *Handler {
	return &Handler{
		messages:   make(chan handler.Message, buffer),
		handleFunc: handleFunc,
		dChan:      make(chan struct{}),
	}
}

func (h *Handler) Put(m handler.Message) error {
	select {
	case h.messages <- m:
	default:
		return errors.New("buffer full")
	}
	return nil
}

func (h *Handler) Start() {
	go h.run()
}

func (h *Handler) Stop() {
	close(h.messages)
	<-h.dChan
}

func (h *Handler) run() {
	for message := range h.messages {
		err := h.handleFunc(message)
		if err != nil {
			util.Logger.Errorf("%s handle message: %s", logPrefix, err)
		}
	}
	h.dChan <- struct{}{}
}
