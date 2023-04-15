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
		{"X   ", args{"    ", "X", 0}, "X   "},
		{" X  ", args{"    ", "X", 1}, " X  "},
		{"  X ", args{"    ", "X", 2}, "  X "},
		{"   X", args{"    ", "X", 3}, "   X"},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := replacePartOfString(tt.args.s, tt.args.new, tt.args.i); got != tt.want {
					t.Errorf("replacePartOfString() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
