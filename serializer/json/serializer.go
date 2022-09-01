package json

import (
	"encoding/json"
	"errors"
	"github.com/DiasOrazbaev/url-shortner/shortner"
)

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*shortner.Redirect, error) {
	redirect := &shortner.Redirect{}
	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errors.New(err.Error() + " serializer.json.Redirect.Decode")
	}
	return redirect, nil
}

func (r *Redirect) Encode(input *shortner.Redirect) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, errors.New(err.Error() + " serializer.json.Redirect.Encode")
	}
	return rawMsg, nil
}
