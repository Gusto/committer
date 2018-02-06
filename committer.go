package main

import (
	"flag"
	"fmt"
	"github.com/gusto/committer/core"
	"os"
)

const VERSION = "0.1.4"

func main() {
	version := flag.Bool("version", false, "Display version")
	help := flag.Bool("help", false, "Display usage")
	fix := flag.Bool("fix", false, "Run autocorrect for commands that support it")
	configPath := flag.String("config", "committer.yml", "Location of your config file")

	flag.Parse()

	if *help {
		fmt.Fprintf(os.Stdout, "Usage of committer:\n")
		flag.PrintDefaults()
		return
	}

	if *version {
		fmt.Printf(VERSION + "\n")
		return
	}

	parsedConfig, err := core.NewConfigFromFile(*configPath)
	if err != nil {
		panic(err)
	}

	success := core.NewRunner(*parsedConfig, *fix).Run()

	if success {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
