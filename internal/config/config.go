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

type conf struct {
	Hits int64 `yaml:"hits"`
	Time int64 `yaml:"time"`
}

func ProcessConfig(config string) *conf {
	f := openConfig(config)      // set f as our byte array
	var c *conf                  // init a struct
	err := yaml.Unmarshal(f, &c) // read yaml into the struct
	if err != nil {
		log.Fatal(err)
	}
	return c
}
