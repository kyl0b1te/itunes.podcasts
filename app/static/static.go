package static

import (
	"errors"
	"io/ioutil"
	"os"
)

type DataEncoder func() ([]byte, error)
type DataDecoder func(body []byte) error

func Save(path string, encoder DataEncoder) error {

	data, err := encoder()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}

func Load(path string, decoder DataDecoder) error {

	if exists, err := isExists(path); err != nil || !exists {
		if err == nil {
			err = errors.New("File is not exist")
		}
		return err
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
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
