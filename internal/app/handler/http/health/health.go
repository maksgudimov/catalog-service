package rhealth

import (
	"net/http"

	rhandler "github.com/maksgudimov/catalog-service/internal/app/handler/http"
)

type handler struct{}

func NewHandler() rhandler.Health {
	return &handler{}
}

func (h *handler) LastCheck(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("ok")); err != nil {
		return
	}
}
