package binding

import (
	"net/http"

	"github.com/maksgudimov/catalog-service/internal/app/entity"
	"github.com/maksgudimov/catalog-service/internal/pkg/http/httph"
)

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "JSON"
}

func (jsonBinding) Bind(req *http.Request, obj any) error {
	if req == nil || req.Body == nil {
		return entity.ErrIncorrectParameters
	}
	if err := httph.DecodeJSON(req, obj); err != nil {
		return entity.ErrIncorrectParameters
	}
	return validate(obj)
}
