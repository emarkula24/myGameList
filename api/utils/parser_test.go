package utils

import "testing"

func TestParseSearchQuery(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseSearchQuery(tt.args.query); got != tt.want {
				t.Errorf("ParseSearchQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
