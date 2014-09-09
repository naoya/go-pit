package pit_test

import (
	"fmt"
	"github.com/naoya/go-pit"
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

	// Test initialization
	conf := pit.Get("twitter.com")
	if conf["username"] != "" {
		t.Errorf("unexpected result: %s", conf["username"])
	}

	// Test default profile
	data := `--- 
twitter.com:
  username: melody
  password: nelson
`

	if err := ioutil.WriteFile(path.Join(d, ".pit", "default.yaml"), []byte(data), 0644); err != nil {
		t.Error(err)
	}

	conf = pit.Get("twitter.com")
	if conf["username"] != "melody" {
		t.Errorf("unexpected result: %s", conf["username"])
	}

	// Test Switch()
	data = `
twitter.com:
  username: development
  password: barbaz
`

	if err := ioutil.WriteFile(path.Join(d, ".pit", "development.yaml"), []byte(data), 0644); err != nil {
		t.Error(err)
	}

	pit.Switch("development")
	conf = pit.Get("twitter.com")

	if conf["username"] != "development" {
		t.Errorf("unexpected result: %s", conf["username"])
	}

	// Test switch back to default
	pit.Switch("default")
	conf = pit.Get("twitter.com")

	if conf["username"] != "melody" {
		t.Errorf("unexpected result: %s", conf["username"])
	}
}

func Example() {
	// Read account information from ~/.pit/default.yaml
	conf := pit.Get("twitter.com")
	fmt.Printf(conf["username"])
	fmt.Printf(conf["password"])

	// Switch profile to development, now using ~/.pit/development.yaml
	pit.Switch("development")
	conf = pit.Get("twitter.com")

	// Switch back to default profile
	pit.Switch("default")
}
