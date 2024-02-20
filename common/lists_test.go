package common

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ReadFile func(filename string) ([]byte, error)

type FakeReadFiler struct {
	Str string
}

func (f FakeReadFiler) ReadFile(filename string) ([]byte, error) {
	buf := bytes.NewBufferString(f.Str)
	return io.ReadAll(buf)
}

func TestReadPublishers(t *testing.T) {
	payload := `---
- name: pcm
  id: pcm
  orgs:
    - https://github.com/italia
  repos: []`

	path := "/dev/null"

	// reassign myReadFile
	fake := FakeReadFiler{Str: payload}
	fileReaderInject = fake.ReadFile
	result, err := LoadPublishers(path)

	if err != nil {
		t.Logf("err: %+v", err)
		t.Fail()
	}

	assert.Len(t, result, 1)
	assert.Nil(t, err)
}
