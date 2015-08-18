/*
Package pit is the main interface to the Go version of pit.

pit is simple account management tool by ruby, perl or something.
You can see more detail on ruby version: https://github.com/cho45/pit
*/
package pit

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

type filePath string

func (f filePath) Read() (data []byte, err error) {
	return ioutil.ReadFile(string(f))
}

func (f filePath) Write(data []byte, perm os.FileMode) (err error) {
	return ioutil.WriteFile(string(f), data, perm)
}

func (f filePath) Exists() bool {
	_, err := os.Stat(string(f))
	return !os.IsNotExist(err)
}

func (f filePath) LoadYaml(out interface{}) (err error) {
	b, err := f.Read()
	return yaml.Unmarshal(b, out)
}

type pit struct {
	directory   string
	configPath  filePath
	profilePath filePath
}

var instance *pit

func GetInstance() *pit {
	if instance == nil {
		d := os.Getenv("HOME")
		if d == "" {
			usr, err := user.Current()
			if err != nil {
				panic(err)
			}
			d = usr.HomeDir
		}
		d = filepath.Join(d, ".pit")
		instance = &pit{
			directory: d,
		}
		instance.SetProfilePath("default")
		instance.configPath = filePath(filepath.Join(d, "pit.yaml"))
	}
	return instance
}

func (pit *pit) SetProfilePath(name string) {
	pit.profilePath = filePath(filepath.Join(pit.directory, name+".yaml"))
}

func (pit pit) CurrentProfile() (profile string) {
	profile = GetInstance().Config().Profile
	return
}

func (pit pit) Load() (profile map[interface{}]interface{}, err error) {
	if _, e := os.Stat(pit.directory); os.IsNotExist(e) {
		err = os.MkdirAll(pit.directory, 0700)
		if err != nil {
			return
		}
	}

	if !pit.configPath.Exists() {
		pit.configPath.Write([]byte("---\nprofile: default\n"), 0600)
	}

	pit.SetProfilePath(pit.CurrentProfile())

	if !pit.profilePath.Exists() {
		pit.profilePath.Write([]byte("--- {}\n"), 0600)
	}
	pit.profilePath.LoadYaml(&profile)
	return
}

type config struct {
	Profile string `yaml:profile`
}

func (pit pit) Config() (c config) {
	pit.configPath.LoadYaml(&c)
	return
}

func (pit *pit) UpdateConfig(name string) (err error) {
	c := config{
		Profile: name,
	}
	b, err := yaml.Marshal(&c)
	if err != nil {
		return
	}
	pit.configPath.Write(b, 0600)
	return
}

type Profile map[string]string

// Get retrieves account information saved under ~/.pit directory. Default profile is ~/.pit/default.yaml
func Get(name string) (profile Profile, err error) {
	profile = make(Profile)
	self := GetInstance()
	m, err := self.Load()

	if err != nil {
		return
	}

	// これもちっとマシに型変換できないのかな...
	if m[name] != nil {
		for k, v := range m[name].(map[interface{}]interface{}) {
			profile[k.(string)] = v.(string)
		}
	}
	return
}

// Switch profile to specified name.
func Switch(name string) (prev string, err error) {
	self := GetInstance()
	prev = self.CurrentProfile()
	self.SetProfilePath(name)
	err = self.UpdateConfig(name)
	return
}
