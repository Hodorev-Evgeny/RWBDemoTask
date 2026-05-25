package core_server

import (
	"fmt"
	"net/http"
)

type ApiVersion string

var (
	ApiVersion1 = ApiVersion("v1")
)

type APIVersionRouter struct {
	*http.ServeMux
	version ApiVersion
}

func NewAPIVersionRouter(
	version ApiVersion,
) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux: http.NewServeMux(),
		version:  version,
	}
}

func (r *APIVersionRouter) RegisterAPIRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.Handler)
	}
}
