//go:build mage

package main

import (
	"fmt"
	"os"
	"strings"

	_ "github.com/magefile/mage/mage"
	"github.com/magefile/mage/sh"
)

var Default = Build

func Build() error {
	if err := sh.RunV("go", "mod", "tidy"); err != nil {
		return err
	}

	env := map[string]string{
		"GO111MODULE": "on",
		"CGO_ENABLE":  "0",
	}

	args := []string{"build"}
	args = append(args, "-o", "./thingscript")

	fmt.Println("Build thingscript...")
	err := sh.RunWithV(env, "go", args...)
	if err != nil {
		return err
	}
	fmt.Println("Build done.")
	return nil
}

func Test() error {
	if err := sh.RunV("go", "test", "./...", "-cover", "-coverprofile", "./cover.out"); err != nil {
		return err
	}
	if output, err := sh.Output("go", "tool", "cover", "-func=./cover.out"); err != nil {
		return err
	} else {
		lines := strings.Split(output, "\n")
		fmt.Println(lines[len(lines)-1])
	}
	return nil
}

func Clean() error {
	os.Remove("thingscript")
	os.Remove("cover.out")
	return nil
}
