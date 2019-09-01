package static

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {

	path := "/tmp/static.save.test.txt"

	err := Save(path, func() ([]byte, error) {

		return []byte{}, errors.New("Encoding error")
	})
	assert.NotNil(t, err)
	assert.Equal(t, "Encoding error", errors.Cause(err).Error())

	err = Save(path, func() ([]byte, error) {

		return []byte("test file"), nil
	})
	assert.Nil(t, err)
	assert.FileExists(t, path)

	os.Remove(path)
}

func TestLoad(t *testing.T) {

	path := "/tmp/static.load.test.txt"
	os.Remove(path)

	err := Load(path, func(body []byte) error { return nil })
	msg := fmt.Sprintf("stat %s: no such file or directory", path)
	assert.Equal(t, msg, errors.Cause(err).Error())

	ioutil.WriteFile(path, []byte("test load"), 0644)
	err = Load(path, func(body []byte) error {
		return errors.New("Decoder error")
	})
	assert.Equal(t, "Decoder error", err.Error())

	Load(path, func(body []byte) error {
		assert.Equal(t, []byte("test load"), body)
		return nil
	})

	defer os.Remove(path)
}
