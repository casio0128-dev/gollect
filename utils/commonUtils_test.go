package utils

import (
	"os"
	"testing"
)

func TestMakePath(t *testing.T) {
	type args struct {
		paths []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "空文字のケース", args: args{paths: []string{""}}, want: ""},
		{name: "最小文字列数のケース", args: args{paths: []string{"test"}}, want: "test"},
		{name: "通常のケース", args: args{paths: []string{"test", "test"}}, want: "test" + string(os.PathSeparator) + "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakePath(tt.args.paths...); got != tt.want {
				t.Errorf("MakePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
