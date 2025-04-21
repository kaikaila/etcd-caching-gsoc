package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	fmt.Println(">>> main started")
	target := flag.String("demo", "memory", "Which demo to run: memory | ctlLog | watch | withLog")
	flag.Parse()

	switch *target {
	case "memory":
		testWatcherWithMemoryCache()
	case "log":
		testWatcherWithCtlOutput()
	case "watch":
		testWatcherWithWatchCache()
	case "withLog":
		testWatcherWithLogOutput()
	default:
		fmt.Println("Unknown demo target. Use -demo=memory | ctlLog | watch | withLog")
		os.Exit(1)
	}
}