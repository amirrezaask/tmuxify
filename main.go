package main

import (
	"fmt"

	gotmux "github.com/jubnzv/go-tmux"
	"gopkg.in/yaml.v3"
	"os"
)

type sessionConfig struct {
	Name    string   `yaml:"name"`
	Cwd     string   `yaml:"cwd"`
	Windows []string `yaml:"windows"`
}

func load(config sessionConfig) (*gotmux.Session, error) {
	server := new(gotmux.Server)
	has, err := server.HasSession(config.Name)
	if err != nil {
		panic(err)
	}

	if has {
		return nil, nil
	}
	sess := gotmux.Session{
		Name: config.Name, StartDirectory: config.Cwd,
	}

	for idx, window := range config.Windows {
		window := gotmux.Window{
			Id:   idx,
			Name: window,
		}
		window.StartDirectory = config.Cwd
		sess.AddWindow(window)
	}

	server.AddSession(sess)

	conf := gotmux.Configuration{
		Server:        server,
		Sessions:      []*gotmux.Session{&sess},
		ActiveSession: nil,
	}

	// Setup this configuration.
	err = conf.Apply()
	if err != nil {
		return nil, fmt.Errorf("Can't apply prepared configuration: %s", err)
	}
	return &sess, nil
}

func main() {
	var filename string
	if len(os.Args) > 1 {
		filename = os.Args[1]	
	} else {
		filename = ".tmuxify.yml"
	}
	
	fd, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	var config sessionConfig
	err = yaml.NewDecoder(fd).Decode(&config)
	if err != nil {
		panic(err)
	}
	sess, err := load(config)

	if err != nil {
		panic(err)
	}
	err = sess.AttachSession()
	if err != nil {
		panic(err)
	}

}
