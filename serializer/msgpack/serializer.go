package msgpack

import (
	"errors"
	"github.com/DiasOrazbaev/url-shortner/shortner"
	"github.com/vmihailenco/msgpack/v5"
)

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*shortner.Redirect, error) {
	redirect := &shortner.Redirect{}
	if err := msgpack.Unmarshal(input, redirect); err != nil {
		return nil, errors.New(err.Error() + " serializer.json.Redirect.Decode")
	}
	return redirect, nil
}

func (r *Redirect) Encode(input *shortner.Redirect) ([]byte, error) {
	rawMsg, err := msgpack.Marshal(input)
	if err != nil {
		return nil, errors.New(err.Error() + " serializer.json.Redirect.Encode")
	}
	return rawMsg, nil
}
