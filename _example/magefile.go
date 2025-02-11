//go:build mage
// +build mage

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var (
	noRunModules = []string{
		"common",
	}
	noApiModules = []string{
		"common",
	}
	modules = []string{"example-service", "common"}
)

func Run(serviceName string) error {
	if !slices.Contains(noRunModules, serviceName) {
		fmt.Printf("Running %s...\n", serviceName)
		cmd := exec.Command("go", "run", ".")
		err := runCmd(
			cmd,
			getServiceDir(serviceName, true),
			nil,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func RunAll() {
	for _, serviceName := range modules {
		Run(serviceName)
	}
}

// A build step that requires additional params, or platform specific steps for example
func Build(serviceName string) error {
	//mg.Deps(CodeGen)
	fmt.Printf("Building %s...\n", serviceName)
	platforms := map[string]string{
		"windows": "GOOS=windows",
	}
	fileName := fmt.Sprintf("../../../%s/example_%s", getTargetDir(serviceName), serviceName)
	for platform, os := range platforms {
		fileName = fmt.Sprintf("%s_%s", fileName, platform)
		if platform == "windows" {
			fileName = fmt.Sprintf("%s.exe", fileName)
		}
		fmt.Printf("Building %s on %s\n", fileName, platform)
		cmd := exec.Command("go", "build", "-o", fileName, ".")
		envs := []string{
			os,
		}
		err := runCmd(
			cmd,
			getServiceDir(serviceName, true),
			envs,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// A build all step that requires additional params, or platform specific steps for example
func BuildAll() {
	for _, serviceName := range modules {
		Build(serviceName)
	}
}

// A test step to run tests for a specific module
func Test(serviceName string) error {
	fmt.Printf("Testing %s...\n", serviceName)
	cmd := exec.Command("go", "test", "./...", "-v", "--coverprofile", fmt.Sprintf("../../%s/cover.out", getTargetDir(serviceName)))
	err := runCmd(
		cmd,
		getServiceDir(serviceName, false),
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// A test all step to run tests for all modules
func TestAll() {
	for _, serviceName := range modules {
		Test(serviceName)
	}
}

func TidyGoMod(serviceName string) error {
	err := updateGoMod(serviceName)
	if err != nil {
		return err
	}
	fmt.Printf("Tidying Go Mod for %s...\n", serviceName)
	cmd := exec.Command("go", "mod", "tidy")
	err = runCmd(
		cmd,
		fmt.Sprintf("services/%s", serviceName),
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func TidyGoModAll() {
	for _, serviceName := range modules {
		TidyGoMod(serviceName)
	}
}

func GenerateApi(serviceName string) error {
	if !slices.Contains(noApiModules, serviceName) {
		err := cleanUpApi(serviceName)
		if err != nil {
			return err
		}
		fmt.Printf("Generating API for %s...\n", serviceName)
		cmd := exec.Command("swag", "init", "-o", "../docs", "--parseDependency", "--parseInternal")
		err = runCmd(
			cmd,
			getServiceDir(serviceName, true),
			nil,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateApiAll() {
	for _, serviceName := range modules {
		GenerateApi(serviceName)
	}
}

func Clean(serviceName string) error {
	fmt.Printf("Cleaning %s...\n", serviceName)
	err := os.RemoveAll(getTargetDir(serviceName))
	if err != nil {
		return err
	}
	return nil
}

func CleanAll() {
	for _, serviceName := range modules {
		Clean(serviceName)
	}
}

func Mock(serviceName string) error {
	if !slices.Contains(noRunModules, serviceName) {
		err := cleanUpMock(serviceName)
		if err != nil {
			return err
		}
		fmt.Printf("Generating Mock Files for %s...\n", serviceName)
		cmd := exec.Command("go", "generate", "./...")
		err = runCmd(
			cmd,
			getServiceDir(serviceName, true),
			nil,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func MockAll() {
	for _, serviceName := range modules {
		Mock(serviceName)
	}
}

func updateGoMod(serviceName string) error {
	fmt.Printf("Updating Go Mod for %s...\n", serviceName)
	cmd := exec.Command("go", "get", "-u", "./...")
	err := runCmd(
		cmd,
		fmt.Sprintf("services/%s", serviceName),
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func cleanUpApi(serviceName string) error {
	fmt.Printf("Cleaning API Docs for %s...\n", serviceName)
	err := os.RemoveAll(fmt.Sprintf("%s/../docs", getServiceDir(serviceName, true)))
	if err != nil {
		return err
	}
	return nil
}

func cleanUpMock(serviceName string) error {
	fmt.Printf("Cleaning Mock Files for %s...\n", serviceName)
	root := getServiceDir(serviceName, true)
	entries, err := os.ReadDir(
		root,
	)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		subDirectory, err := os.Open(
			fmt.Sprintf("%s/%s", root, entry.Name()),
		)
		if err != nil {
			return err
		}
		defer subDirectory.Close()
		subEntries, err := subDirectory.Readdir(-1)
		if err != nil {
			return err
		}
		for _, subEntry := range subEntries {
			if !subEntry.Mode().IsRegular() {
				continue
			}
			if strings.Contains(subEntry.Name(), "_mock_test.go") {
				err = os.Remove(
					fmt.Sprintf("%s/%s/%s", root, entry.Name(), subEntry.Name()),
				)
				if err != nil {
					return err
				}
			}

		}
	}
	return nil
}

func getServiceDir(serviceName string, withSrc bool) string {
	serviceDir := fmt.Sprintf("./services/%s", serviceName)
	if withSrc {
		serviceDir = fmt.Sprintf("%s/src", serviceDir)
	}
	return serviceDir
}

func getTargetDir(serviceName string) string {
	target := fmt.Sprintf("./target/%s", serviceName)
	_ = os.MkdirAll(target, os.ModePerm)
	return target
}

func runCmd(cmd *exec.Cmd, dir string, envs []string) error {
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), envs...)
	buffer := bytes.NewBuffer([]byte{})
	cmd.Stdout = buffer
	cmd.Stderr = buffer
	err := cmd.Run()
	fmt.Print(buffer.String())
	if err != nil {
		return err
	}
	return nil
}
