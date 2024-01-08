/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/enp0s3/wasp-operator/internal/defaults"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"os"
	"os/exec"
	"strings"
)

var (
	startOpts struct {
		fsRoot       string
		strategy     string
		nodeName     string
		swappiness   string
		swapSizeMb   string
		swapFilePath string
		swapFileName string
		verbose      bool
		dryRun       bool
	}
)

func main() {
	val, ok := os.LookupEnv("FSROOT")
	if !ok || val == "" {
		startOpts.fsRoot = defaults.CmdFsRoot
	} else {
		startOpts.fsRoot = val
	}

	val, ok = os.LookupEnv("STRATEGY")
	if !ok || val == "" {
		startOpts.strategy = defaults.CmdStrategy
	} else {
		startOpts.strategy = val
	}

	val, ok = os.LookupEnv("SWAPINESS")
	if !ok || val == "" {
		startOpts.swappiness = defaults.CmdSwappiness
	} else {
		startOpts.swappiness = val
	}

	val, ok = os.LookupEnv("SWAP_SIZE_MB")
	if !ok || val == "" {
		startOpts.swapSizeMb = defaults.CmdSwapSizeMb
	} else {
		startOpts.swapSizeMb = val
	}

	val, ok = os.LookupEnv("VERBOSE")
	if !ok || val == "" {
		startOpts.verbose = false
	} else {
		startOpts.verbose = true
	}

	val, ok = os.LookupEnv("DRY_RUN")
	if !ok || val == "" {
		startOpts.dryRun = false
	} else {
		startOpts.dryRun = true
	}

	startOpts.swapFilePath = fmt.Sprintf("%s/var/tmp", startOpts.fsRoot)
	startOpts.swapFileName = "wasp.file"

	absPath := fmt.Sprintf("%s/%s", startOpts.swapFilePath, startOpts.swapFileName)
	if isSwapOn(startOpts.swapFileName) {
		out, err := RunCommand("swapoff", []string{
			"-v",
			absPath,
		})
		if err != nil {
			klog.Fatal(err)
		}
		klog.Info(out)
	}
	createSwap(absPath, startOpts.swapSizeMb)
}

func isSwapOn(swapfile string) bool {
	f, err := os.Open("/proc/swaps")
	if err != nil {
		klog.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), swapfile) {
			return true
		}
	}

	return false
}

func createSwap(swapFile, swapSize string) {
	out, err := RunCommand("dd", []string{
		"if=/dev/zero",
		fmt.Sprintf("of=%s", swapFile),
		"bs=1M",
		fmt.Sprintf("count=%s", startOpts.swapSizeMb),
	})
	if err != nil {
		klog.Fatal(err)
	}
	klog.Info(out)

	out, err = RunCommand("chmod", []string{"0600", swapFile})
	if err != nil {
		klog.Fatal(err)
	}
	klog.Info(out)

	out, err = RunCommand("mkswap", []string{swapFile})
	if err != nil {
		klog.Fatal(err)
	}
	klog.Info(out)

	out, err = RunCommand("swapon", []string{swapFile})
	if err != nil {
		klog.Fatal(err)
	}
	klog.Info(out)
}

func RunCommand(command string, args []string) (string, error) {
	if startOpts.verbose {
		klog.Infof("[*] %s %s", command, strings.Join(args, " "))
	}

	if startOpts.dryRun {
		return "", nil
	}

	return runCommand(command, args)
}

func runCommand(command string, args []string) (string, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(command, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", errors.Wrapf(err, "Failed to run command: %s %s\nStdout:\n%s\nStderr:\n%s",
			command, strings.Join(args, " "), cmd.Stdout, cmd.Stderr)
	}

	return stdout.String(), nil
}
