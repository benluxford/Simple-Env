package senv_test

import (
	"log"

	"github.com/benluxford/Senv"
)

func ExampleLoad() {
	if err := senv.Load("absolute/path/to/env.file"); err != nil {
		// handle error...
	}
}

func ExampleGetVar() {
	if variable, err := senv.GetVar("some key"); err == nil {
		// do something with variable
		log.Println(variable)
	} else {
		// handle error...
	}
}

func ExampleSetVar() {
	if err := senv.SetVar("key string", "value string"); err != nil {
		// handle error...
	}
}
