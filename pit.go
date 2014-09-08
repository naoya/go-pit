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
	self := GetInstance()
	m := self.Config()
	profile = m["profile"].(string)
	return
}

func (pit pit) Load() (profile map[interface{}]interface{}) {
	pit.SetProfilePath(pit.CurrentProfile())

	// TODO: ファイル無いとき -> {} な pit.profile を作る
	b := pit.profilePath.Read()
	err := yaml.Unmarshal(b, &profile)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (pit pit) Config() (profile map[interface{}]interface{}) {
	// TODO: ファイル無いとき -> {"profile" => "default"} な pit.yaml を作る
	b := pit.configPath.Read()
	err := yaml.Unmarshal(b, &profile)
	if err != nil {
		log.Fatal(err)
	}
	return
}

type config struct {
	Profile string `yaml:profile`
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
