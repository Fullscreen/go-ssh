package ssh

import (
	"fmt"
	"strings"
	"testing"
)

var (
	defaultSSHOptions = strings.Join(New("ohai").FormattedOptions(), " ")

	sshCommandTests = []struct {
		session  *Session
		expected string
	}{
		{New("foo"), "foo -t"},
		{New("foo@bar"), "foo@bar -t"},
	}
)

func TestSSHCommandCompilation(t *testing.T) {
	for _, tt := range sshCommandTests {
		actual := strings.Join(tt.session.Command(), " ")
		if actual != fmt.Sprintf("%s %s", tt.expected, defaultSSHOptions) {
			t.Errorf("expected \"%s\", actual \"%s\"", tt.expected, actual)
		}
	}
}
