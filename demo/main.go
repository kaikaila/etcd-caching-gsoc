// demo/main.go
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	target := flag.String("demo", "memory", "Which demo to run: memory | log | watch")
	flag.Parse()

	switch *target {
	case "memory":
		testWatcherWithMemoryCache()
	case "log":
		testWatcherWithLogOutput()
	case "watch":
		testWatcherWithWatchCache()
	default:
		fmt.Println("Unknown demo target. Use -demo=memory | log | watch")
		os.Exit(1)
	}
}