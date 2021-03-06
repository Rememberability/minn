package main

import (
	"bytes"
	"strings"
	"fmt"
	"testing"
)

func TestHelpOption(t *testing.T) {

	cmd := cmdNewContext(dbOpenTemp())
	outbuf := bytes.NewBuffer([]byte{})
	errbuf := bytes.NewBuffer([]byte{})

	for _, command := range cmdList {

		cmd.Stdout = outbuf
		cmd.Stderr = errbuf
		
		for _, option := range []string{ "-h", "--help" } {

			outbuf.Reset()
			errbuf.Reset()

			status := cmd.Run(command.name, option)
			if status != 0 {
				t.Errorf("Non-zero exit status when passing %s to %s subcommand.\n", option, command.name) 
			}

			if errbuf.Len() != 0 {
				t.Errorf("Subcommand %s wrote to stderr when passed %s option.\n", command.name, option)
			}

			if !strings.HasPrefix(outbuf.String(), fmt.Sprintf("Usage: %s %s", mainProgramName, command.name)) {
				t.Errorf("The output of %s %s did not have the expected initial output.\n", command.name, option)
			}

			if !strings.Contains(outbuf.String(), "\n") {
				t.Errorf("The output of %s %s does not contain a newline.\n", command.name, option)
			}

			if !strings.HasSuffix(outbuf.String(), "\n") {
				t.Errorf("The output of %s %s does not end of newline.\n", command.name, option)
			}

			lines := strings.Split(outbuf.String(), "\n")
			hasDescriptionHeader := false
			hasOptionsHeader := false
			for _, line := range lines {
				if line == "Description:" {
					hasDescriptionHeader = true
				}
				if line == "Options:" {
					hasOptionsHeader = true
				}
			}
			if !hasDescriptionHeader {
				t.Errorf("The output of %s %s does not have a Description: header.\n", command.name, option)
			}
			if !hasOptionsHeader {
				t.Errorf("The output of %s %s does not have an Options: header.\n", command.name, option)
			}
		}
	}
}

