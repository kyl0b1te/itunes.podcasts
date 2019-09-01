package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zhikiri/itunes.podcasts/app/genre"
)

func main() {

	opt, err := newOptions()
	stopOnError(err)

	log.Println("--- Start ---")

	switch *opt.Type {
	case "genres":
		log.Println("Parsing genres...")
		genres, errs := genre.GetGenres(genre.GetRequestOptions())
		stopOnErrors(errs)

		log.Println(len(genres), "was successfully parsed")
		genre.Save(*opt.Out, genres)
	default:
		//
	}

	log.Println("--- END ---")
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
