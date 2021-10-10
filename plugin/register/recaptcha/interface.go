package recaptcha

import (
	"github.com/line-api/model/go/model"
	"net/http"
)

type Solver interface {
	Solve(details *model.WebAuthDetails, client *http.Client) (string, error)
}
