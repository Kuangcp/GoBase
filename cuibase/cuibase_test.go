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

func TestPrintParam(t *testing.T) {
	type args struct {
		verb     string
		param    string
		comment  string
		verbLen  int
		paramLen int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "one",
			args: args{
				verb:     "-h|h",
				param:    "",
				comment:  "help info",
				verbLen:  -3,
				paramLen: -12,
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			format := BuildFormat(tt.args.verbLen, tt.args.paramLen)
			PrintParam(format, tt.args.verb, tt.args.param, tt.args.comment)
		})
	}
}

func TestPrintParams(t *testing.T) {
	type args struct {
		params []ParamInfo
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "one",
			args: args{[]ParamInfo{{Verb: "-h", Param: "", Comment: "help"}}},
		},
	}
	for _, tt := range tests {
		format := BuildFormat(-2, -6)
		t.Run(tt.name, func(t *testing.T) {
			PrintParams(format, tt.args.params)
		})
	}
}

func Test_runAction(t *testing.T) {
	type args struct {
		params        []string
		actions       map[string]func(params []string)
		defaultAction func(params []string)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "one",
			args: args{
				params:        []string{"run.go", "-h"},
				defaultAction: func(params []string) { print("default") },
				actions:       map[string]func(params []string){"-h": func(params []string) { println("help info") }}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runAction(tt.args.params, tt.args.actions, tt.args.defaultAction)
		})
	}
}
