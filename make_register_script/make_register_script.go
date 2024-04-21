// Copyright 2024 Ross Light
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//		 https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

func main() {
	vaultExe := flag.String("vault", "vault", "`path` to Vault executable")
	flag.Parse()

	outputPath := os.Getenv("out")
	if outputPath == "" {
		fmt.Fprintln(os.Stderr, "$out not set")
		os.Exit(1)
	}

	pluginsData, err := os.ReadFile(os.Getenv("pluginsPath"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var plugins []struct {
		Type        string `json:"type"`
		ProgramName string `json:"pname"`
		Version     string `json:"version"`
	}
	if err := json.Unmarshal(pluginsData, &plugins); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	signal.Ignore(syscall.SIGPIPE)

	buf := new(bytes.Buffer)
	for _, p := range plugins {
		command := p.ProgramName
		if p.Version != "" {
			command += "-" + p.Version
		}
		sum, err := sha256sum(filepath.Join(outputPath, "bin", command))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		buf.Reset()
		buf.WriteString(escapeShellArg(*vaultExe))
		buf.WriteString(" plugin register -sha256=")
		buf.WriteString(sum)
		buf.WriteString(" -command=")
		buf.WriteString(escapeShellArg(command))
		if p.Version != "" {
			buf.WriteString(" -version=")
			buf.WriteString(escapeShellArg(p.Version))
		}
		buf.WriteString(" ")
		buf.WriteString(escapeShellArg(p.Type))
		buf.WriteString(" ")
		buf.WriteString(escapeShellArg(p.ProgramName))
		buf.WriteString("\n")
		if _, err := os.Stdout.Write(buf.Bytes()); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func sha256sum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("sha256sum: %w", err)
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("sha256sum %s: %v", path, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func escapeShellArg(arg string) string {
	const singleQuoteReplacement = `'\''`
	sb := new(strings.Builder)
	sb.Grow(len("''") + len(arg) + len(singleQuoteReplacement)*strings.Count(arg, "'"))
	sb.WriteString("'")
	for _, b := range []byte(arg) {
		if b == '\'' {
			sb.WriteString(singleQuoteReplacement)
		} else {
			sb.WriteByte(b)
		}
	}
	sb.WriteString("'")
	return sb.String()
}
