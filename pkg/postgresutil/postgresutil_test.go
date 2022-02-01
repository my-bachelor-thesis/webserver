package postgresutil

import "testing"

func TestGeneratePlaceholders(t *testing.T) {
	want := "$1, $2, $3"
	got := GeneratePlaceholders("id, is_admin, first_name")
	if want != got {
		t.Errorf("want %q, but got %q", want, got)
	}
}

func TestGeneratePlaceholdersAndReplace(t *testing.T) {
	want := "$1, $2, REPLACED, $3"
	got := GeneratePlaceholdersAndReplace("a, b, c, d", map[int]string{2:"REPLACED"})
	if want != got {
		t.Errorf("want %q, but got %q", want, got)
	}
}