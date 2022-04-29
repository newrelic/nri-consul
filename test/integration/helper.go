//go:build integration
// +build integration

package tests

import (
	"bytes"
	"fmt"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/xeipuuv/gojsonschema"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func dockerCompose(vars, cmdLine, containers []string, detached bool) (string, string, error) {
	if detached {
		cmdLine = append(cmdLine, "-d")
	}
	for i := range vars {
		cmdLine = append(cmdLine, "-e")
		cmdLine = append(cmdLine, vars[i])
	}
	cmdLine = append(cmdLine, containers...)
	fmt.Printf("executing: docker-compose %s\n", strings.Join(cmdLine, " "))
	cmd := exec.Command("docker-compose", cmdLine...)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	stdout := outbuf.String()
	stderr := errbuf.String()
	return stdout, stderr, err
}

func dockerComposeUp(vars, containers []string) (string, string, error) {
	return dockerCompose(vars, []string{"up"}, containers, true)
}

func dockerComposeRun(vars []string, container string) (string, string, error) {
	return dockerCompose(vars, []string{"run", "--rm", "--name", container}, []string{container}, false)
}

func validateJSONSchema(fileName string, input string) error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	schemaURI := fmt.Sprintf("file://%s", filepath.Join(pwd, "testdata", fileName))
	log.Info("loading schema from %s", schemaURI)
	schemaLoader := gojsonschema.NewReferenceLoader(schemaURI)
	documentLoader := gojsonschema.NewStringLoader(input)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("Error loading JSON schema, error: %v", err)
	}

	if result.Valid() {
		fmt.Println("no-error?")
		return nil
	}
	fmt.Printf("Errors for JSON schema: '%s'\n", schemaURI)
	for _, desc := range result.Errors() {
		fmt.Printf("\t- %s\n", desc)
	}
	fmt.Printf("\n")
	return fmt.Errorf("the output of the integration doesn't have expected JSON format")
}
