package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func main() {

	genFl := initBoolFlag("g", "genre", "parse genres")
	shoFl := initBoolFlag("s", "show", "parse shows")
	detFl := initBoolFlag("d", "detail", "parse details")
	fedFl := initBoolFlag("f", "feed", "parse feed")
	comFl := initBoolFlag("c", "compact", "generate compact list of shows")

	outFl := flag.String("out", "/tmp", "generated files folder")
	chuFl := flag.Int("chunk", 100, "details parsing chunk")
	delFl := flag.Int("delay", 5, "delay between chunked requests")
	flag.Parse()

	if *genFl == true && *shoFl == true && *detFl == true && *fedFl == true {
		stopOnError(errors.New("Invalid arguments"))
	}

	if *genFl == true {

		actionGenres(*outFl)
	} else if *shoFl == true {

		actionShows(getFilePathFromArg(), *outFl)
	} else if *detFl == true {

		actionDetails(getFilePathFromArg(), *delFl, *chuFl, *outFl)
	} else if *fedFl == true {

		actionFeed(getFilePathFromArg(), *outFl)
	} else if *comFl == true {

		actionCompact(getFilePathFromArg(), *outFl)
	}

	fmt.Println("Done")
	os.Exit(0)
}

func getFilePathFromArg() string {
	if len(flag.Args()) == 0 || flag.Arg(0) == "" {
		stopOnError(errors.New("File path is missing"))
	}
	return flag.Arg(0)
}

func initBoolFlag(short string, full string, desc string) *bool {

	var fl bool
	flag.BoolVar(&fl, full, false, desc)
	flag.BoolVar(&fl, short, false, desc+" (shorthand)")
	return &fl
}

func stopOnErrors(errs []error) {

	if len(errs) == 0 {
		return
	}
	for _, err := range errs {
		fmt.Printf("[ERROR] %s\n", err)
	}
	os.Exit(1)
}

func stopOnError(err error) {

	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
		os.Exit(1)
	}
}
