package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	CNIVersion string `json:"cniVersion"`
	Name       string `json:"name"`
}

func main() {
	command := os.Getenv("CNI_COMMAND")
	if command == "" {
		fmt.Fprintln(os.Stderr, "CNI_COMMAND is required")
		os.Exit(1)
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil || cfg.CNIVersion == "" || cfg.Name == "" {
		fmt.Fprintln(os.Stderr, "invalid CNI configuration")
		os.Exit(1)
	}

	switch command {
	case "ADD":
		fmt.Fprintf(os.Stderr, "XYZ CNI demo: ADD for network %s (simulation only)\n", cfg.Name)
		fmt.Fprintln(os.Stderr, "No interfaces or routes were changed.")
		fmt.Printf("{\"cniVersion\":%q}\n", cfg.CNIVersion)
	case "CHECK", "DEL":
		fmt.Fprintf(os.Stderr, "XYZ CNI demo: %s for network %s (simulation only)\n", command, cfg.Name)
	case "VERSION":
		fmt.Println(`{"cniVersion":"1.0.0","supportedVersions":["1.0.0"]}`)
	default:
		fmt.Fprintln(os.Stderr, "unsupported CNI command")
		os.Exit(1)
	}
}
