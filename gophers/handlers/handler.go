package handlers

import (
	"gcalsync/gophers/core"
	"net/http"
)

// TODO: Refactor signatures as Writes happen only in response middleware

//go:generate mockery --name=Handler --dir=./ --output=mocks --outpkg=mocks
type Handler interface {
	CallbackHandler(http.ResponseWriter, *http.Request) (interface{}, error)
	ListEventsHandler(http.ResponseWriter, *http.Request) (interface{}, error)
	ConnectHandler(http.ResponseWriter, *http.Request) (interface{}, error)
	NotifyHandler(http.ResponseWriter, *http.Request) (interface{}, error)
}

type handler struct {
	core core.Core
}

func New() (Handler, error) {
	service, err := core.New()
	if err != nil {
		return nil, err
	}
	return &handler{
		core: service,
	}, nil
}
