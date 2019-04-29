package config

import (
	yaml "gopkg.in/yaml.v2"
	"log"
	"os"
)

var (
	maxConfigBytes int    = 1024 * 8 // max config file size
	readBytesCount int               // number of bytes read from open config file
	readBytes      []byte            // the byte array of size maxConfigBytes
	Config         ConfigFile
)

func openConfig(config string) []byte {
	file, err := os.Open(config)
	if err != nil {
		log.Fatal(err)
	}

	readBytes := make([]byte, maxConfigBytes) // set aside 8kb of memory
	readBytesCount, err := file.Read(readBytes)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("read %d bytes: %s\n", readBytesCount, config)
	return readBytes[:readBytesCount]
}

type Metric struct {
	Description string `yaml:"description`
	Query       string `yaml:"query"`
}

type ConfigFile map[string]Metric // set up a new type

func ProcessConfig(config string) interface{} {
	f := openConfig(config)           // set f as our byte array
	err := yaml.Unmarshal(f, &Config) // read yaml into the map
	if err != nil {
		log.Fatal(err)
	}
	return Config
}
