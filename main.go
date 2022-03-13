package main

import (
	"flag"
	"github.com/rtakhautdinov/myhttp/internal/app"
	"log"
)

const (
	// DefaultMaxParallel is max value for parallel runner that tool's user could request
	DefaultMaxParallel = 32
	// DefaultParallelExecution is how many parallel request could be processed(if PC supports)
	DefaultParallelExecution = 10
)

var (
	numOfParallel int
	urls          []string
)

func init() {
	flag.IntVar(&numOfParallel, "parallel", DefaultParallelExecution, "Number of parallel requests to be executed")
}

func processArgs() {
	flag.Parse()
	urls = flag.Args()

	if len(urls) == 0 {
		log.Fatalf("No urls provided")
	}

	if numOfParallel < 0 || numOfParallel > DefaultMaxParallel {
		log.Fatalf("Invalid %d parallel provided. Expected: from 1 to %d", numOfParallel, DefaultMaxParallel)
	}
}

func main() {
	processArgs()
	app.Run(urls, numOfParallel)
}
