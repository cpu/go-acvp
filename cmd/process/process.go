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
	toolPath := flag.String("tool", "", "Path to the boringssl trim_vectors tool")
	wrapperPath := flag.String("wrapper", "", "Path to the Go acvp_wrapper test binary")
	algos := flag.String("algorithms", "all", "Comma-separated list of algorithms")
	upload := flag.Bool("upload", false, "Whether to upload answers to NIST ACVTS demo server")

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

	for _, alg := range algorithms {
		alg = strings.ReplaceAll(strings.TrimSpace(alg), "/", "-")
		fmt.Printf("processing algorithm: %q\n", alg)

		vectorFile := filepath.Join("vectors", alg)
		answersFile := filepath.Join("expected", alg)
		err := cmd.RunAcvpTool(*toolPath, *wrapperPath, "process", vectorFile, answersFile, true)
		if err != nil {
			log.Fatalf("error: processing vectors for %q: %s\n", alg, err)
		}

		if *upload {
			err = cmd.RunAcvpTool(*toolPath, *wrapperPath, "upload", answersFile, "", false)
			if err != nil {
				log.Fatalf("error: uploading vector answers for %q: %s\n", alg, err)
			}
		}
	}
}
