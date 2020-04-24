package crawler

import (
	"bytes"
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"
)

type ReadFile func(filename string) ([]byte, error)

type FakeReadFiler struct {
	Str string
}

func (f FakeReadFiler) ReadFile(filename string) ([]byte, error) {
	buf := bytes.NewBufferString(f.Str)
	return ioutil.ReadAll(buf)
}

func TestReadWhitelists(t *testing.T) {
	payload := `- name: pcm
  codice-iPA: pcm
  orgs:
    - https://github.com/italia
  repos: []`

	path := "/dev/null"

	// reassign myReadFile
	fake := FakeReadFiler{Str: payload}
	fileReaderInject = fake.ReadFile
	result2, err := ReadAndParseWhitelist(path)

	if err != nil {
		t.Logf("err: %+v", err)
		t.Fail()
	}

	log.Printf("ReadAndParseWhitelist == %#v, %#v\n", result2, err)
}
