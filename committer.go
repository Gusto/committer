package main

import (
	"flag"
	"fmt"
	"github.com/gusto/committer/core"
	"io/ioutil"
	"os"
)

func main() {
	content, _ := ioutil.ReadFile("VERSION")
	VERSION := string(content)

	version := flag.Bool("version", false, "Display version")
	help := flag.Bool("help", false, "Display usage")
	fix := flag.Bool("fix", false, "Run autocorrect for commands that support it")
	changed := flag.Bool("changed", false, "Run autocorrect for commands that support it")
	configPath := flag.String("config", "committer.yml", "Location of your config file")

	flag.Parse()

	if *help {
		fmt.Fprintf(os.Stdout, "Usage of committer:\n")
		flag.PrintDefaults()
		return
	}

	if *version {
		fmt.Printf(VERSION)
		return
	}

	parsedConfig, err := core.NewConfigFromFile(*configPath)
	if err != nil {
		return
	}

	success := core.NewRunner(*parsedConfig, *fix, *changed).Run()

	if success {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
