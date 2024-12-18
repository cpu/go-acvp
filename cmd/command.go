package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dsnet/compress/bzip2"
)

var DefaultAlgs = []string{
	"SHA2-224",
	"SHA2-256",
	"SHA2-384",
	"SHA2-512",
	"SHA2-512/224",
	"SHA2-512/256",

	"SHA3-224",
	"SHA3-256",
	"SHA3-384",
	"SHA3-512",

	"SHAKE-128",
	"SHAKE-256",

	"HMAC-SHA2-224",
	"HMAC-SHA2-256",
	"HMAC-SHA2-384",
	"HMAC-SHA2-512",
	"HMAC-SHA2-512/224",
	"HMAC-SHA2-512/256",

	"HMAC-SHA3-224",
	"HMAC-SHA3-256",
	"HMAC-SHA3-384",
	"HMAC-SHA3-512",

	"PBKDF",

	"ECDSA",
	"DetECDSA",
	"EDDSA",

	"CMAC-AES",

	"KDA",

	"TLS-v1.2",
	"TLS-v1.3",

	"kdf-components",

	"ML-KEM",
}

func RunAcvpTool(pathToTool, pathToWrapper, action, input, outputFile string, compress bool) error {
	// All actions require the wrapper.
	args := []string{fmt.Sprintf("-wrapper=%s", pathToWrapper)}

	switch action {
	case "fetch":
		args = append(args, "-fetch", input)
	case "process":
		args = append(args, "-json", input)
	case "upload":
		args = append(args, "-upload", input)
	default:
		return fmt.Errorf("unknown action: %q", action)
	}

	cmd := exec.Command(pathToTool, args...)

	// Signal to the wrapper that it's being used as an acvptool subprocess, not a Go test.
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "ACVP_WRAPPER=1")

	if outputFile != "" {
		outFile, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer outFile.Close()
		cmd.Stdout = outFile
	} else {
		cmd.Stdout = os.Stdout
	}

	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	if compress {
		raw, err := os.ReadFile(outputFile)
		if err != nil {
			return err
		}

		if err = Compress(raw, outputFile); err != nil {
			return err
		}

	}

	return nil
}

func Compress(data []byte, path string) error {
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
