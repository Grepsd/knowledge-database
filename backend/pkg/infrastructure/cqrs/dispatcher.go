package cqrs
//
//import "errors"
//
//var ErrNoCommandHandlerFound = errors.New("no command handler found supporting this command")
//var ErrCommandHandlerAlreadyRegistered = errors.New("command handler already registered")
//
//type QueryDispatcher interface {
//	Handle(query Queryer)
//	Register(handler QueryHandler) error
//}
//
//type CommandDispatcher interface {
//	Handle(command Commander) error
//	Register(handler CommandHandler) error
//}
//
//type SimpleCommandDispatcher struct {
//	handlers []CommandHandler
//}
//
//func (s SimpleCommandDispatcher) Handle(command Commander) error {
//	for _, handler := range s.handlers {
//		if !handler.Supports(command) {
//			continue
//		}
//		return handler.Handle(command)
//	}
//	return ErrNoCommandHandlerFound
//}
//
//func (s SimpleCommandDispatcher) Register(handler CommandHandler) error {
//	for _, h := range s.handlers {
//		if h == handler {
//			return ErrCommandHandlerAlreadyRegistered
//		}
//	}
//	s.handlers = append(s.handlers, handler)
//	return nil
//}