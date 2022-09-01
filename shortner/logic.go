package shortner

import (
	"errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
	"time"
)

var (
	ErrRedirectNotFound = errors.New("redirect not found")
	ErrRedirectInvalid  = errors.New("redirect invalid")
)

type redirectService struct {
	redirectRepository RedirectRepository
}

func NewRedirectService(redirectRepo RedirectRepository) *redirectService {
	return &redirectService{redirectRepo}
}

func (r *redirectService) Find(code string) (*Redirect, error) {
	return r.redirectRepository.Find(code)
}

func (r *redirectService) Save(redirect *Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return errors.New(err.Error() + " " + ErrRedirectNotFound.Error() + " service.Redirect.Store")
	}
	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC().Unix()
	return r.redirectRepository.Save(redirect)
}
