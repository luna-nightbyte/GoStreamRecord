package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
)

type Exec struct {
	Command     string
	Args        []string
	stdoutBytes []byte
	stderrBytes []byte
}

// ProcessCommand executes a given command string in the terminal.
// It captures the command's standard output and standard error.
func (c *Exec) Process() error {
	if c.Command == "" {
		return fmt.Errorf("no command configured.")
	}
	cmdName := c.Command
	cmdArgs := c.Args

	// Create a new command to be executed.
	cmd := exec.Command(cmdName, cmdArgs...)
	// Create pipes for stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating stdout pipe: %v\n", err)
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Error creating stderr pipe: %v\n", err)
		return err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		return err
	}

	// Read from the pipes
	c.stdoutBytes, _ = io.ReadAll(stdoutPipe)
	c.stderrBytes, _ = io.ReadAll(stderrPipe)

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Command finished with error: %v\n", err)
	}
	return nil
}

func (c *Exec) Error() error {
	if (string(c.stderrBytes)) != "" {
		return errors.New((string(c.stderrBytes)))
	}
	return nil
}

func (c *Exec) Output() string {
	return string(c.stdoutBytes)
}
func CheckPath(command string) string {
	cmd := exec.Command("sh", "-c", "which "+command)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		log.Printf("Unable to locate '%s' binary installation. Add it to your $PATH. \nOr move/create symbolic link for the binary to a known path. I.E 'sudo cp /home/$USER/.local/bin/%s /usr/local/bin/'", command, command)
		fmt.Printf("Unable to locate '%s' binary installation. Add it to your $PATH. \nOr move/create symbolic link for the binary to a known path. I.E 'sudo cp /home/$USER/.local/bin/%s /usr/local/bin/'", command, command)
		return ""
	}

	pwd := bytes.TrimSpace(stdout.Bytes())
	return string(pwd)

}
