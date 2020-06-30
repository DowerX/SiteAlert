package config

import (
	"io/ioutil"

	"github.com/DowerX/SiteAlert/errorcheck"
	"gopkg.in/yaml.v2"
)

// Config _
type Config struct {
	Url      string
	Email    string
	Sender   string
	Password string
	Server   string
	Port     string
	Msg      string
	Time     int
	Log      bool
	Logfile  string
}

func GetConfig(path string) Config {
	var data, err = ioutil.ReadFile(path)
	errorcheck.Check(err)
	c := Config{}
	err = yaml.Unmarshal(data, &c)
	errorcheck.Check(err)
	return c
}
