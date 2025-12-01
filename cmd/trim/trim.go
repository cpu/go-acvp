package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dsnet/compress/bzip2"
)

func main() {
	algos := flag.String("algorithms", "all", "Comma-separated list of algorithms")

	flag.Parse()

	var algorithms []string

	if *algos != "all" {
		algorithms = strings.Split(*algos, ",")
	}

	if len(algorithms) == 0 {
		log.Fatal("one or more -algorithms must be specified")
	}

	for _, alg := range algorithms {
		alg = strings.ReplaceAll(strings.TrimSpace(alg), "/", "-")
		fmt.Printf("trimming algorithm: %q\n", alg)

		vectorFile := filepath.Join("vectors", alg)
		vectorData, err := os.ReadFile(vectorFile)
		if err != nil {
			log.Fatalf("error: reading vectors for %q: %s", alg, err)
		}

		trimmed, err := trim(vectorData)
		if err != nil {
			log.Fatalf("error: trimming vectors for %q: %s", alg, err)
		}

		if err = os.WriteFile(vectorFile, trimmed, 0644); err != nil {
			log.Fatalf("error: writing trimmed vectors for %q: %s", alg, err)
		}

		if err = compress(trimmed, vectorFile); err != nil {
			log.Fatalf("error: compressing vectors for %q: %s", alg, err)
		}
	}
}

func trim(vectors []byte) ([]byte, error) {
	var vectorSets []any
	if err := json.Unmarshal(vectors, &vectorSets); err != nil {
		return nil, fmt.Errorf("unmarshaling vectors: %w", err)
	}

	// The first element is the metadata which is left unmodified.
	for i := 1; i < len(vectorSets); i++ {
		vectorSet := vectorSets[i].(map[string]any)
		testGroups := vectorSet["testGroups"].([]any)
		for _, testGroupInterface := range testGroups {
			testGroup := testGroupInterface.(map[string]any)
			tests := testGroup["tests"].([]any)

			keepIndex := 10
			if keepIndex >= len(tests) {
				keepIndex = len(tests) - 1
			}

			testGroup["tests"] = []any{tests[keepIndex]}
		}
	}

	trimmed, err := json.MarshalIndent(vectorSets, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshalling trimmed vectors: %w", err)
	}

	return trimmed, nil
}

func compress(data []byte, path string) error {
	if !strings.HasSuffix(path, ".bz2") {
		path = path + ".bz2"
	}

	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating %q: %w", path, err)
	}
	defer func() { _ = outFile.Close() }()

	bw, err := bzip2.NewWriter(outFile, nil)
	if err != nil {
		return fmt.Errorf("constructing bzip2 writer: %w", err)
	}
	defer func() { _ = bw.Close() }()

	if _, err := bw.Write(data); err != nil {
		return fmt.Errorf("compressing data: %w", err)
	}

	return nil
}
