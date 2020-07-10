package cqrs

import (
	"errors"
	"time"
)

var ErrCannotStartQueryHandleWhenAlreadyStarted = errors.New("cannot start Query handle if already started")
var ErrCannotEndQueryHandleWhenAlreadyEnded = errors.New("cannot end Query handle if already ended")
var ErrNoQueryHandlerFound = errors.New("no query handler found supporting this query")
var ErrQueryHandlerAlreadyRegistered = errors.New("query handler already registered")

type Queryer interface {
	CreatedOn() time.Time
	HandleStartedOn() time.Time
	HandleEndedOn() time.Time
}

type QueryHandler interface {
	Handle(query Queryer) error
	Supports(query Queryer) bool
}

type SimpleQuery struct {
	createdOn       time.Time
	handleStartedOn time.Time
	handleEndedOn   time.Time
}

type QueryDispatcher interface {
	Handle(query Queryer)
	Register(handler QueryHandler) error
}

type SimpleQueryDispatcher struct {
	handlers []QueryHandler
}

func NewSimpleQueryDispatcher() SimpleQueryDispatcher {
	return SimpleQueryDispatcher{handlers: []QueryHandler{}}
}

func NewSimpleQuery() SimpleQuery {
	return SimpleQuery{createdOn: time.Time{}, handleStartedOn: time.Time{}, handleEndedOn: time.Time{}}
}

func (s SimpleQuery) CreatedOn() time.Time {
	return s.createdOn
}

func (s SimpleQuery) HandleStartedOn() time.Time {
	return s.handleStartedOn
}

func (s SimpleQuery) HandleEndedOn() time.Time {
	return s.handleEndedOn
}

func (s SimpleQuery) StartHandle() error {
	if !s.handleStartedOn.IsZero() {
		return ErrCannotStartQueryHandleWhenAlreadyStarted
	}
	s.handleStartedOn = time.Now()
	return nil
}

func (s SimpleQuery) EndHandle() error {
	if !s.handleEndedOn.IsZero() {
		return ErrCannotEndQueryHandleWhenAlreadyEnded
	}
	s.handleEndedOn = time.Now()
	return nil
}

func (s SimpleQueryDispatcher) Handle(query Queryer) error {
	for _, handler := range s.handlers {
		if !handler.Supports(query) {
			continue
		}
		return handler.Handle(query)
	}
	return ErrNoQueryHandlerFound
}

func (s SimpleQueryDispatcher) Register(handler QueryHandler) error {
	for _, h := range s.handlers {
		if h == handler {
			return ErrQueryHandlerAlreadyRegistered
		}
	}
	s.handlers = append(s.handlers, handler)
	return nil
}
