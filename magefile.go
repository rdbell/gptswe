//go:build mage
// +build mage

package main

import (
	"log"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Build command
func Build() error {
	output := "gptswe"
	if mg.Verbose() {
		log.Println("Building", output)
	}

	return sh.RunV("go", "build", "-o", output, ".")
}

// Test command
func Test() error {
	if mg.Verbose() {
		log.Println("Running tests")
	}

	return sh.RunV("go", "test", "./...")
}

// Publish command
func Publish() error {
	if mg.Verbose() {
		log.Println("Publishing Docker image")
	}

	return sh.RunV("bash", "./publish.sh")
}
