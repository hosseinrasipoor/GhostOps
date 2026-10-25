package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	outFile := flag.String("o", "", "export json file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: ghostops [-o out.json] <logfile>\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("error: enter log file path")
		flag.Usage()
		os.Exit(1)
	}

	path := args[0]
	fmt.Printf("reading file: %s (out: %s)\n", path, *outFile)
}
