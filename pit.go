package pit

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type pit struct {
	directory   string
	configPath  string
	profilePath string
}

var instance *pit

func GetInstance() *pit {
	if instance == nil {
		d := path.Join(os.Getenv("HOME"), ".pit")
		instance = &pit{
			directory: d,
		}
		instance.SetProfilePath("default")
		instance.configPath = path.Join(d, "pit.yaml")
	}
	return instance
}

func (pit *pit) SetProfilePath(name string) {
	pit.profilePath = path.Join(pit.directory, name+".yaml")
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
	b, err := ioutil.ReadFile(pit.profilePath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(b, &profile)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (pit pit) Config() (profile map[interface{}]interface{}) {
	b, err := ioutil.ReadFile(pit.configPath)

	// TODO: ファイル無いとき -> {"profile" => "default"} な pit.yaml を作る
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(b, &profile)
	if err != nil {
		log.Fatal(err)
	}
	return
}

type Profile map[string]string
type Config struct {
	Profile string `yaml:profile`
}

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

	c := Config{
		Profile: name,
	}

	// FIXME: エラー無視してる
	b, _ := yaml.Marshal(&c)
	err := ioutil.WriteFile(self.configPath, b, 0600)
	if err != nil {
		log.Fatal(err)
	}
	return
}
