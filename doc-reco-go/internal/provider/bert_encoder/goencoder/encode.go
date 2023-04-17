package goencoder

import (
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/bert"
)

var (
	poolingStrategy = bert.ReduceMean
)

func Vectorize(text string) ([]float32, error) {

	if len(text) > 512 {
		text = text[:512]
	}

	mat, err := Model.Vectorize(text, poolingStrategy)
	return mat.Data(), err
}
