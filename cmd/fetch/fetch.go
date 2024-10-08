package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cpu/go-acvp/cmd"
)

func main() {
	toolPath := flag.String("tool", "", "Path to the boringssl acvptool binary")
	wrapperPath := flag.String("wrapper", "", "Path to the Go acvp_wrapper test binary")
	algos := flag.String("algorithms", "all", "Comma-separated list of algorithms")

	flag.Parse()

	if *toolPath == "" || *wrapperPath == "" {
		fmt.Println("error: both -tool and -wrapper are required")
		os.Exit(1)
	}

	var algorithms []string

	if *algos != "all" {
		algorithms = strings.Split(*algos, ",")
	} else {
		algorithms = cmd.DefaultAlgs
	}

	err := os.MkdirAll("vectors", os.ModePerm)
	if err != nil {
		log.Fatalf("error: creating vectors dir: %s", err)
	}

	err = os.MkdirAll("expected", os.ModePerm)
	if err != nil {
		log.Fatalf("error: creating expected dir: %s", err)
	}

	for _, alg := range algorithms {
		alg = strings.TrimSpace(alg)
		algFile := strings.ReplaceAll(alg, "/", "-")
		fmt.Printf("processing algorithm: %q\n", alg)

		vectorFile := filepath.Join("vectors", algFile)
		err := cmd.RunAcvpTool(*toolPath, *wrapperPath, "fetch", alg, vectorFile, true)
		if err != nil {
			log.Fatalf("error: fetching vectors for %q: %s\n", alg, err)
		}
	}
}
