package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/davlgd/goaccess-logs-converter/internal/config"
	"github.com/davlgd/goaccess-logs-converter/internal/tools"
)

var (
	name        = "goaccess-logs-converter"
	version     = "dev"
	buildTime   = "unknown"
	description = "Convert Clever Cloud access logs to use with GoAccess (COMMON format)"
)

func main() {
	config.InitFlags()
	flag.Usage = config.Usage(name, version, description)
	flag.Parse()

	if config.ShouldShowVersion() {
		fmt.Printf("%s version %s (built at %s)\n", name, version, buildTime)
		os.Exit(0)
	}

	if err := tools.CheckCleverAuth(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg, err := config.New()
	if err != nil {
		fmt.Printf("Configuration error: %v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}

	if err := tools.Process(cfg); err != nil {
		fmt.Printf("Process error: %v\n", err)
		os.Exit(1)
	}
}
