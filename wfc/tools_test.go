package wfc

import (
	"testing"
)

func Test_replacePartOfString(t *testing.T) {
	type args struct {
		s   string
		new string
		i   int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// Error
		{"Err 1", args{"    ", "X", -1}, "    "},
		{"Err 2", args{"    ", "X", 5}, "    "},

		// No error
		{"XX      ", args{"        ", "XX", 0}, "XX      "},
		{"  XX    ", args{"        ", "XX", 1}, "  XX    "},
		{"    XX  ", args{"        ", "XX", 2}, "    XX  "},
		{"      XX", args{"        ", "XX", 3}, "      XX"},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := replaceCharInString(tt.args.s, tt.args.new, tt.args.i); got != tt.want {
					t.Errorf("replaceCharInString() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
