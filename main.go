package ssh

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// OptionDefaults are default ssh options
var OptionDefaults = []Option{
	Option{Name: "ServerAliveInterval", Value: "10"},
	Option{Name: "ConnectTimeout", Value: "10"},
	Option{Name: "LogLevel", Value: "Error"},
	Option{Name: "StrictHostKeyChecking", Value: "no"},
	Option{Name: "UserKnownHostsFile", Value: "/dev/null"},
}

// Session represents an ssh session
type Session struct {
	User    string
	Host    string
	Port    int
	Options []Option

	Stdin          io.Reader
	Stdout, Stderr io.Writer
}

// Option represents and ssh option
type Option struct {
	Name  string
	Value string
}

// New returns a new SSH Session instance
func New(host string) *Session {
	var user string
	parts := strings.SplitN(host, "@", 2)
	if len(parts) > 1 {
		user = parts[0]
		host = parts[1]
	} else {
		host = parts[0]
	}

	return &Session{
		User:    user,
		Host:    host,
		Port:    22,
		Options: OptionDefaults,

		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Connect initiates an interactive SSH session
func (s *Session) Connect() error {
	return s.Run([]string{})
}

// Run runs an SSH session
func (s *Session) Run(args []string) error {
	sshArgs := append(s.Command(), s.FormattedOptions()...)
	cmd := exec.Command("ssh", append(sshArgs, args...)...)
	cmd.Stdin = s.Stdin
	cmd.Stdout = s.Stdout
	cmd.Stderr = s.Stderr

	cmd.Start()
	return cmd.Wait()
}

// Command returns the command arguments for the ssh command
func (s *Session) Command() []string {
	var command []string
	host := s.Host
	if s.User != "" {
		host = fmt.Sprintf("%s@%s", s.User, host)
	}
	port := strconv.Itoa(s.Port)
	// add tty flag
	command = append(command, host, "-p", port, "-t")
	// append options
	for _, o := range s.Options {
		command = append(command, "-o", fmt.Sprintf("%s=%s", o.Name, o.Value))
	}
	return command
}

// FormattedOptions returns a list of SSH options
func (s *Session) FormattedOptions() []string {
	var options []string
	for _, opt := range s.Options {
		options = append(options, "-o", fmt.Sprintf("%s=%s", opt.Name, opt.Value))
	}
	return options
}

// Run runs an ssh command
func Run(sshOptions []string, command string, session *Session) error {
	// add default ssh configuration
	sshOptions = append(sshOptions, []string{"-o", "LogLevel=ERROR", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "-o", "ConnectTimeout=10", "-t", command}...)

	cmd := exec.Command("ssh", sshOptions...)
	cmd.Stdin = session.Stdin
	cmd.Stdout = session.Stdout
	cmd.Stderr = session.Stderr

	cmd.Start()
	return cmd.Wait()
}
