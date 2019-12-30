package main

import (
	"bytes"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"strings"

	"github.com/gessnerfl/instana-terraform-dump/modules"
	"github.com/gessnerfl/instana-terraform-dump/rest"
)

func main() {
	apiKeyFlag := flag.String("key", "", "The api key of the Instana REST interface")
	hostFlag := flag.String("host", "", "The hostname of the Instana REST interface")
	outFileFlag := flag.String("out", "", "The target file path")
	flag.Parse()

	outFile := strings.TrimSpace(*outFileFlag)
	if len(outFile) == 0 {
		panic(errors.New("Output file not provided"))
	}

	restClient, err := createClient(apiKeyFlag, hostFlag)
	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer
	moduleFactory := modules.NewFactory(restClient)

	log.Println("Start processing of custom events")
	err = moduleFactory.CustomEvents().AppendInstanaResourcesTo(&buffer)
	log.Println("Processing of custom events completed")
	if err != nil {
		panic(err)
	}

	log.Printf("Write data to output file %s", outFile)
	ioutil.WriteFile(outFile, buffer.Bytes(), 0644)
}

func createClient(apiKeyFlag *string, hostFlag *string) (rest.Client, error) {
	apiKey := strings.TrimSpace(*apiKeyFlag)
	host := strings.TrimSpace(*hostFlag)

	if len(apiKey) == 0 {
		return nil, errors.New("API Key not provided")
	}

	if len(host) == 0 {
		return nil, errors.New("Instana hostname not provided")
	}

	return rest.NewClient(apiKey, host), nil
}
