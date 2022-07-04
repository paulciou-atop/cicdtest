package version

import "testing"

func TestFormat(t *testing.T) {
	expects := "nconfig version 1.4.0 (2020-12-15)\n"
	if got := Format("1.4.0", "2020-12-15"); got != expects {
		t.Errorf("Format() = %q, wants %q", got, expects)
	}
}

