package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

const configFile = "./config.yml"

func main() {
	listFlag := flag.Bool("l", false, "list the url mapping")

	flag.Parse()

	if *listFlag {
		listURLMapping()
		return
	}

}

func listURLMapping() {
	bytes, err := ioutil.ReadFile(configFile)

	if err != nil {
		exitWithError(err)
	}

	urlMap := make(map[string]string)

	if err = yaml.Unmarshal(bytes, &urlMap); err != nil {
		exitWithError(err)
	}

	for k, v := range urlMap {
		fmt.Printf("%v - %v\n", k, v)
	}
}

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
