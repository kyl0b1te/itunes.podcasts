package static

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type DataEncoder func() ([]byte, error)
type DataDecoder func(body []byte) error

func Save(path string, encoder DataEncoder) error {

	data, err := encoder()
	if err != nil {
		return errors.Wrap(err, "Cannot save data")
	}

	return ioutil.WriteFile(path, data, 0644)
}

func Load(path string, decoder DataDecoder) error {

	if exists, err := isExists(path); err != nil || !exists {
		return errors.Wrap(err, "File is not exists")
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "Cannot read file")
	}

	return decoder(file)
}

func isExists(path string) (bool, error) {

	_, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}
