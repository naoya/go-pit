package pit

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestGet(t *testing.T) {
	d, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error(err)
	}

	defer func() {
		if err := os.RemoveAll(d); err != nil {
			t.Error(err)
		}
	}()

	if err := os.Setenv("HOME", d); err != nil {
		t.Error(err)
	}

	conf := Get("twitter.com")
	if conf["username"] != "" {
		t.Errorf("unexpected result: %s", conf["username"])
	}

	data := `--- 
twitter.com:
  username: melody
  password: nelson
`

	if err := ioutil.WriteFile(path.Join(d, ".pit", "default.yaml"), []byte(data), 0644); err != nil {
		t.Error(err)
	}

	conf = Get("twitter.com")
	if conf["username"] != "melody" {
		t.Errorf("unexpected result: %s", conf["username"])
	}
}
