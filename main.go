package main

import (
	"fmt"

	gotmux "github.com/jubnzv/go-tmux"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
)

type sessionConfig struct {
    Name string `yaml:"name"`
	Cwd string `yaml:"cwd"`
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
		window := gotmux.Window {
			Id: idx,
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
	loadCmd := &cobra.Command{
		Use:   "load",
		Short: "short",
		Long:  "long",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				panic("need a config file name")
			}
			fd, err := os.Open(args[0])
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
			fmt.Println(sess)
			// err = sess.AttachSession()
			// if err != nil {
			// 	panic(err)
			// }

		},
	}
	rootCmd := &cobra.Command{
		Use:   "tmuxify",
		Short: "short",
		Long:  "long",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	rootCmd.AddCommand(loadCmd)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
