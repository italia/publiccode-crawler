package crawler

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
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
	payload := `---
- name: pcm
  codice-iPA: pcm
  orgs:
    - https://github.com/italia
  repos: []`

	path := "/dev/null"

	// reassign myReadFile
	fake := FakeReadFiler{Str: payload}
	fileReaderInject = fake.ReadFile
	result, err := ReadAndParseWhitelist(path)

	if err != nil {
		t.Logf("err: %+v", err)
		t.Fail()
	}

	assert.Len(t, result, 1)
	assert.Nil(t, err)
}

func TestReadBlacklists(t *testing.T) {
	payload := `---
repos:
  - url: https://github.com/italia/repo1
    reason: 1
    description: GitHub takedown
  - url: https://github.com/italia/repo2
    reason: 2
    description: GitHub takedown`

	path := "/dev/null"

	// reassign myReadFile
	fake := FakeReadFiler{Str: payload}
	fileReaderInject = fake.ReadFile
	result, err := ReadAndParseBlacklist(path)

	if err != nil {
		t.Logf("err: %+v", err)
		t.Fail()
	}

	assert.Len(t, result, 2)
	assert.Nil(t, err)

	// disabled at the moment, need to figure out how to mock filepath.*
	// assert.True(t, IsRepoInBlackList("https://github.com/italia/repo1"))
	// assert.True(t, IsRepoInBlackList("https://github.com/italia/repo2"))
	// assert.False(t, IsRepoInBlackList("https://github.com/italia/repo3"))
}
