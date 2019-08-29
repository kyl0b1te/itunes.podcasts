package main

import (
	"flag"

	"github.com/pkg/errors"
)

type Options struct {
	Type *string
	Chunk *int
	Out *string
}

func newOptions() (*Options, error) {

	opt := &Options{
		Type: flag.String("parse", "", "content type to parse (genre|shows|details|feed)"),
		Chunk: flag.Int("chunk", 0, "data parsing chunk"),
		Out: flag.String("out", "/tmp/data.json", "path to the file with results"),
	}
	flag.Parse()

	if !opt.isTypeValid() {
		return &Options{}, errors.New("Invalid data type")
	}

	return opt, nil
}

func (o *Options) isTypeValid() bool {

	types := []string{"genres", "shows", "details", "feed"}
	for _, t := range types {
		if *o.Type == t {
			return true
		}
	}
	return false
}
