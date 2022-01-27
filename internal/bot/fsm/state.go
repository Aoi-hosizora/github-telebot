package fsm

import (
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
	"sync"
)

const (
	None xtelebot.ChatState = iota

	// BindingUsername (controller.Subscribe)
	BindingUsername

	// BindingToken (controller.Subscribe)
	BindingToken
)

func StateString(state xtelebot.ChatState) string {
	switch state {
	case BindingUsername, BindingToken:
		return "subscribe"
	default:
		return "?"
	}
}

var (
	_handlers  = make(map[xtelebot.ChatState]xtelebot.MessageHandler)
	_handlerMu = sync.RWMutex{}
)

func GetStateHandler(state xtelebot.ChatState) xtelebot.MessageHandler {
	_handlerMu.RLock()
	handler, ok := _handlers[state]
	_handlerMu.Unlock()
	if !ok {
		return nil
	}
	return handler
}

func IsHandlerRegistered(state xtelebot.ChatState) bool {
	_handlerMu.RLock()
	_, ok := _handlers[state]
	_handlerMu.RUnlock()
	return ok
}

func RegisterHandler(state xtelebot.ChatState, handler xtelebot.MessageHandler) {
	_handlerMu.Lock()
	_handlers[state] = handler
	_handlerMu.Unlock()
}
