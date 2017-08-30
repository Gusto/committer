package main

import (
	"flag"
	"fmt"
	"github.com/nikhilmat/committer/core"
	// "log"
	"os"
	// "sync"
)

func main() {
	help := flag.Bool("help", false, "Display usage")
	fix := flag.Bool("fix", false, "Run autocorrect for commands that support it")
	config := flag.String("config", "committer.yml", "Location of your config file")

	flag.Parse()

	if *help {
		fmt.Fprintf(os.Stdout, "Usage of committer:\n")
		flag.PrintDefaults()
		return
	}

	parsedConfig, err := core.NewConfig(*config)
	if err != nil {
		return
	}

	core.Run(parsedConfig, *fix)
	// var wg sync.WaitGroup
	// logger := log.New(os.Stdout, "", 0)
	//
	// cmd := core.NewCmd("bundle exec rubocop app/models/company.rb")
	//
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	logger.Print("Running: ", cmd.Command)
	// 	output, _ := cmd.Execute().Output()
	// 	fmt.Println(string(output))
	// }()
	//
	// wg.Wait()
}
