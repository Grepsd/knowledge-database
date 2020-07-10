package cqrs

import (
	"errors"
	"time"
)

var ErrCannotStartCommandHandleWhenAlreadyStarted = errors.New("cannot start command handle if already started")
var ErrCannotEndCommandHandleWhenAlreadyEnded = errors.New("cannot end command handle if already ended")
var ErrNoCommandHandlerFound = errors.New("no command handler found supporting this command")
var ErrCommandHandlerAlreadyRegistered = errors.New("command handler already registered")

type Commander interface {
	CreatedOn() time.Time
	HandleStartedOn() time.Time
	HandleEndedOn() time.Time
}

type CommandHandler interface {
	Handle(command Commander) error
	Supports(command Commander) bool
}

type SimpleCommand struct {
	createdOn       time.Time
	handleStartedOn time.Time
	handleEndedOn   time.Time
}

type CommandDispatcher interface {
	Handle(command Commander) error
	Register(handler CommandHandler) error
}

type SimpleCommandDispatcher struct {
	handlers []CommandHandler
}

func NewSimpleCommandDispatcher() SimpleCommandDispatcher {
	return SimpleCommandDispatcher{handlers: []CommandHandler{}}
}

func NewSimpleCommand() SimpleCommand {
	return SimpleCommand{createdOn: time.Time{}, handleStartedOn: time.Time{}, handleEndedOn: time.Time{}}
}

func (s SimpleCommand) CreatedOn() time.Time {
	return s.createdOn
}

func (s SimpleCommand) HandleStartedOn() time.Time {
	return s.handleStartedOn
}

func (s SimpleCommand) HandleEndedOn() time.Time {
	return s.handleEndedOn
}

func (s SimpleCommand) StartHandle() error {
	if !s.handleStartedOn.IsZero() {
		return ErrCannotStartCommandHandleWhenAlreadyStarted
	}
	s.handleStartedOn = time.Now()
	return nil
}

func (s SimpleCommand) EndHandle() error {
	if !s.handleEndedOn.IsZero() {
		return ErrCannotEndCommandHandleWhenAlreadyEnded
	}
	s.handleEndedOn = time.Now()
	return nil
}

func (s SimpleCommandDispatcher) Handle(command Commander) error {
	for _, handler := range s.handlers {
		if !handler.Supports(command) {
			continue
		}
		return handler.Handle(command)
	}
	return ErrNoCommandHandlerFound
}

func (s SimpleCommandDispatcher) Register(handler CommandHandler) error {
	for _, h := range s.handlers {
		if h == handler {
			return ErrCommandHandlerAlreadyRegistered
		}
	}
	s.handlers = append(s.handlers, handler)
	return nil
}
