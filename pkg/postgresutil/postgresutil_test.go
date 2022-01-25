package postgresutil

import "testing"

func TestGeneratePlaceholder(t *testing.T) {
	want := "$1, $2, $3"
	got := GeneratePlaceholder("id, is_admin, first_name")
	if want != got {
		t.Errorf("want %q, but got %q", want, got)
	}
}