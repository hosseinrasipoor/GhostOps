package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"ghostops/analyzer"
	"ghostops/parser"
	"os"
	"time"
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
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	an := analyzer.NewAnalyzer()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		entry, err := parser.ParseLine(sc.Text())
		if err != nil {
			an.AddBad()
			continue
		}
		an.AddGood(entry)
	}

	if err := sc.Err(); err != nil {
		fmt.Printf("error reading file: %v\n", err)
		os.Exit(1)
	}

	res := an.GetResult()

	fmt.Println("--- ghostops report ---")
	fmt.Printf("total lines: %d\n", res.TotalLines)
	fmt.Printf("valid entries: %d\n", res.ValidEntries)
	fmt.Printf("malformed lines: %d\n\n", res.MalformedLines)

	fmt.Println("severities:")
	for lvl, count := range res.Severities {
		fmt.Printf("%s: %d\n", lvl, count)
	}

	fmt.Println("\nservices:")
	for svc, count := range res.Services {
		fmt.Printf("%s: %d\n", svc, count)
	}

	fmt.Println("\nincidents:")
	if len(res.Incidents) == 0 {
		fmt.Println("no incidents found")
	} else {
		for _, inc := range res.Incidents {
			fmt.Printf("- type: %s\n", inc.Type)
			fmt.Printf("  severity: %s\n", inc.Severity)
			fmt.Printf("  service: %s\n", inc.Service)
			fmt.Printf("  count: %d\n", inc.Count)
			fmt.Printf("  first: %s\n", inc.FirstSeen.Format(time.RFC3339))
			fmt.Printf("  last: %s\n", inc.LastSeen.Format(time.RFC3339))
			fmt.Printf("  desc: %s\n", inc.Desc)
		}
	}

	if *outFile != "" {
		data, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			fmt.Printf("error marshaling json: %v\n", err)
			os.Exit(1)
		}
		err = os.WriteFile(*outFile, data, 0644)
		if err != nil {
			fmt.Printf("error writing file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\njson saved in %s\n", *outFile)
	}
}
