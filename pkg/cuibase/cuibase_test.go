package cuibase

import (
	"testing"
)

func TestAssertCount(t *testing.T) {
	list := make([]string, 1)
	flag := enoughCount(list, 2)
	if flag {
		t.Fail()
	}
	t.Log(flag)

	list = make([]string, 2)
	flag = enoughCount(list, 1)
	if !flag {
		t.Fail()
	}
}

func Test_enoughCount(t *testing.T) {
	type args struct {
		param []string
		count int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "use ",
			args: args{
				param: []string{"file", "param1"},
				count: 1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := enoughCount(tt.args.param, tt.args.count); got != tt.want {
				t.Errorf("enoughCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHelpInfo_PrintHelp(t *testing.T) {
	info := HelpInfo{
		Description:   "Translation between Chinese and English By Baidu API",
		Version:       "1.0.1",
		SingleFlagLen: -2,
		DoubleFlagLen: -8,
		ValueLen:      -9,
		Flags: []ParamVO{
			{Short: "-s", Value: "<value>", Comment: "use jkjfdk jk kj "},
			{Long: "--s", Value: "<value>", Comment: "use jkjfdk jk kj "},
		},
		Options: []ParamVO{{Short: "-s", Long: "--s", Value: "<value>", Comment: "use jkjfdk jk kj "}},
	}
	info.PrintHelp()
}

func TestColor(t *testing.T) {
	//PrintWithColorful()
	print(Red.Print("ddd"))
}
