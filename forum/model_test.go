package forum

import (
	"testing"
)

var calculateCreationDiffProvider = []struct {
	in  string
	out string
}{
	{"%a", "[%a]"},
	{"%-a", "[%-a]"},
	{"%+a", "[%+a]"},
	{"%#a", "[%#a]"},
	{"% a", "[% a]"},
	{"%0a", "[%0a]"},
	{"%1.2a", "[%1.2a]"},
	{"%-1.2a", "[%-1.2a]"},
	{"%+1.2a", "[%+1.2a]"},
	{"%-+1.2a", "[%+-1.2a]"},
	{"%-+1.2abc", "[%+-1.2a]bc"},
	{"%-1.2abc", "[%-1.2a]bc"},
}

func TestCalculateCreationDiff(t *testing.T) {
	var calculateCreationDiffProvider calculateCreationDiffProvider
	for i, tt := range calculateCreationDiffProvider {
		s := Sprintf(tt.in, &flagprinter)
		if s != tt.out {
			t.Errorf("%d. Sprintf(%q, &flagprinter) => %q, want %q", i, tt.in, s, tt.out)
		}
	}
}


