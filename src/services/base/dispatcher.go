package base

const (
	ServiceDispatcher = "dispatcher"
	DefaultBufSize    = 65535
)

type DispatcherHander interface{}

// async producer/consumer
type IDispatcherService interface {
	CreateDispatcher(dispatcherName string, params ...uint) error
	RemoveDispatcher(dispatcherName string) error
	Dispatch(dispatcherName string, data interface{}) error
	AddDispatcherHandler(dispatcherName string, handler DispatcherHander, handlerName ...string) error
	RemoveDispatcherHandler(dispatcherName string, handlerName string) error
}
