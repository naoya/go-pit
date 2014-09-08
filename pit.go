package pit

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type filePath string

func (f filePath) Read() (data []byte) {
	data, err := ioutil.ReadFile(string(f))
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (f filePath) Write(data []byte, perm os.FileMode) {
	err := ioutil.WriteFile(string(f), data, perm)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (f filePath) Exists() bool {
	_, err := os.Stat(string(f))
	return !os.IsNotExist(err)
}

func (f filePath) LoadYaml(out interface{}) {
	err := yaml.Unmarshal(f.Read(), out)
	if err != nil {
		log.Fatal(err)
	}
}

type pit struct {
	directory   string
	configPath  filePath
	profilePath filePath
}

var instance *pit

func GetInstance() *pit {
	if instance == nil {
		d := path.Join(os.Getenv("HOME"), ".pit")
		instance = &pit{
			directory: d,
		}
		instance.SetProfilePath("default")
		instance.configPath = filePath(path.Join(d, "pit.yaml"))
	}
	return instance
}

func (pit *pit) SetProfilePath(name string) {
	pit.profilePath = filePath(path.Join(pit.directory, name+".yaml"))
}

func (pit pit) CurrentProfile() (profile string) {
	profile = GetInstance().Config().Profile
	return
}

func (pit pit) Load() (profile map[interface{}]interface{}) {
	if _, err := os.Stat(pit.directory); os.IsNotExist(err) {
		err = os.MkdirAll(pit.directory, 0700)
		if err != nil {
			log.Fatal(err)
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
	m := make(map[interface{}]interface{})
	pit.configPath.LoadYaml(&m)
	c = config{
		Profile: m["profile"].(string),
	}
	return
}

func (pit *pit) UpdateConfig(name string) {
	c := config{
		Profile: name,
	}
	b, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatal(err)
	}
	pit.configPath.Write(b, 0600)
}

type Profile map[string]string

func Get(name string) (profile Profile) {
	self := GetInstance()
	m := self.Load()

	// これもちっとマシに型変換できないのかな...
	profile = make(Profile)
	for k, v := range m[name].(map[interface{}]interface{}) {
		profile[k.(string)] = v.(string)
	}
	return
}

func Switch(name string) (prev string) {
	self := GetInstance()
	prev = self.CurrentProfile()
	self.SetProfilePath(name)
	self.UpdateConfig(name)
	return
}
