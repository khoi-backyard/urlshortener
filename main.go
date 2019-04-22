package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	yaml "gopkg.in/yaml.v2"
)

const configFile = "./config.yml"

func main() {
	if len(os.Args) == 1 {
		printHelp()
		return
	}

	bytes, err := ioutil.ReadFile(configFile)

	if err != nil {
		exitWithError(err)
	}

	urlMap := make(map[string]string)

	if err = yaml.Unmarshal(bytes, &urlMap); err != nil {
		exitWithError(err)
	}

	runFlagSet := flag.NewFlagSet("run", flag.ExitOnError)
	runPort := runFlagSet.Int("p", 8080, "port of the server")

	configureFlagSet := flag.NewFlagSet("configure", flag.ExitOnError)
	addingKey := configureFlagSet.String("a", "", "the key")
	addingURL := configureFlagSet.String("u", "", "the url")

	listFlag := flag.Bool("l", false, "list the url mapping")
	helpFlag := flag.Bool("h", false, "print the help message")
	deleteFlag := flag.String("d", "", "delete an entry using key")

	switch os.Args[1] {
	case "run":
		runFlagSet.Parse(os.Args[2:])
	case "configure":
		configureFlagSet.Parse(os.Args[2:])
	default:
		flag.Parse()
	}

	// Run subcommand
	if runFlagSet.Parsed() {
		http.HandleFunc("/", redirectionHandler(urlMap))
		port := fmt.Sprintf(":%d", *runPort)
		fmt.Printf("Starting server on port %v\n", port)
		log.Fatal(http.ListenAndServe(port, nil))
		return
	}

	// Configure subcommand
	if configureFlagSet.Parsed() {
		if *addingKey == "" || *addingURL == "" {
			printHelp()
			os.Exit(1)
		}
		urlMap[*addingKey] = *addingURL

		if err := writeURLMappingToFile(urlMap, configFile); err != nil {
			exitWithError(err)
		}

		return
	}

	if flag.Parsed() {
		if *listFlag {
			printURLMapping(urlMap)
			return
		}
		if *helpFlag {
			printHelp()
			return
		}
		if *deleteFlag != "" {
			delete(urlMap, *deleteFlag)

			if err := writeURLMappingToFile(urlMap, configFile); err != nil {
				exitWithError(err)
			}

			return
		}
	}

}

func redirectionHandler(urlMap map[string]string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]
		if url, ok := urlMap[path]; ok {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			return
		}
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func writeURLMappingToFile(urlMap map[string]string, filePath string) error {
	b, err := yaml.Marshal(urlMap)

	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filePath, b, 0644); err != nil {
		return err
	}

	return nil
}

func printURLMapping(urlMap map[string]string) {
	for k, v := range urlMap {
		fmt.Printf("%v - %v\n", k, v)
	}
}

func printHelp() {
	fmt.Printf("urlshortener\n\n")

	fmt.Println("To add new entry:")
	fmt.Println("\t urlshortener configure -a dogs -u www.dogs.com")

	fmt.Println("To delete an entry:")
	fmt.Println("\t urlshortener -d dogs")

	fmt.Println("To list all entries:")
	fmt.Println("\t urlshortener -l")

	fmt.Println("To start the server:")
	fmt.Println("\t urlshortener run -p 8080")
}

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
