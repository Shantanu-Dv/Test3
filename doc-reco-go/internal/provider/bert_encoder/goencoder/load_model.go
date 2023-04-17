package goencoder

import (
	"fmt"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/bert"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/huggingface"
	"os"
	"path"
	"path/filepath"
)

const defaultModelFile = "spago_model.bin"

const (
	modelDir = "model"
)

var Model *bert.Model

func LoadModel(modelName string) error {

	if Model != nil {
		return nil
	}

	modelPath := filepath.Join(modelDir, modelName)
	var err error

	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		fmt.Printf("Unable to find `%s` locally.\n", modelPath)
		fmt.Printf("Pulling `%s` from Hugging Face models hub...\n", modelDir)
		// make sure the models path exists
		if _, err := os.Stat(modelDir); os.IsNotExist(err) {
			if err := os.MkdirAll(modelDir, 0755); err != nil {
				return err
			}
		}
		err = huggingface.NewDownloader(modelDir, modelName, false).Download()
		if err != nil {
			return err
		}
		fmt.Printf("Converting modelName...\n")
		err = huggingface.NewConverter(modelDir, modelName).Convert()
		if err != nil {
			return err
		}
	} else if _, err := os.Stat(path.Join(modelPath, defaultModelFile)); os.IsNotExist(err) {
		fmt.Printf("Unable to find `%s` in the modelName directory.\n", defaultModelFile)
		fmt.Printf("Assuming there is a Hugging Face modelName to convert...\n")
		err = huggingface.NewConverter(modelDir, modelName).Convert()
		if err != nil {
			return err
		}
	}

	Model, err = bert.LoadModel(modelPath)
	if err != nil {
		return err
	}

	return nil
}
