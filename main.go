package main

import (
	"bufio"
	"csvReaderTemplateApp/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	// evaluate command line flags
	var help bool
	var verbose bool
	var file string
	flags := flag.NewFlagSet("csvReaderTemplateApp", flag.ContinueOnError)
	flags.BoolVar(&help, "help", false, "Show this help message")
	flags.BoolVar(&help, "h", false, "")
	flags.BoolVar(&verbose, "v", false, "Show verbose logging.")
	flags.StringVar(&file, "file", "", "csv filename to parse")
	if len(os.Args) < 2 {
		printHelp(flags)
		os.Exit(1)
	}
	err := flags.Parse(os.Args[1:])
	switch err {
	case flag.ErrHelp:
		help = true
	case nil:
	default:
		log.Fatalf("error parsing flags: %v", err)
	}
	// If the help flag was set, just show the help message and exit.
	if help {
		printHelp(flags)
		os.Exit(0)
	}
	if file == "" {
		file = os.Args[1]
	}

	// check if csv file exists
	if !fileExists(file) {
		log.Println("CSV file", file, "does not exist or is a directory")
		os.Exit(1)
	}

	os.Exit(csvReader(file))
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// printHelp prints command line parameter help
func printHelp(flags *flag.FlagSet) {
	fmt.Fprintf(flags.Output(), "\nUsage of %s:\n", os.Args[0])
	flags.PrintDefaults()
}

func csvReader(file string) int {
	reader, err := getResponseReader(file)
	if err != nil {
		log.Println("Error opening file:", err)
	}

	r := csv.NewFieldReader(reader)
	r.Comma = ','
	r.Comment = '#'

	for r.Scan() {
		// print map of fieldnames: values
		log.Println(r.Fields())

		// print specific fields for every line in file
		log.Println(r.Field("username"), r.Field("firstname"), r.Field("lastname"))
	}

	if err := r.Err(); err != nil {
		log.Println(err)
		return 1
	}

	return 0
}

func getResponseReader(filename string) (io.Reader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return bufio.NewReader(file), nil
}
