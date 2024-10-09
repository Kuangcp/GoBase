package tk

import (
	"github.com/kuangcp/gobase/pkg/ctool"
	"testing"
)

import (
	"github.com/expr-lang/expr"
	"modernc.org/tk9.0/b5"
)

import . "modernc.org/tk9.0"

// https://pkg.go.dev/modernc.org/tk9.0
// https://gitlab.com/cznic/tk9.0
func TestHi(t *testing.T) {
	Pack(Button(Txt("Hello"), Command(func() { Destroy(App) })))
	App.Wait()
}

func TestCalc(t *testing.T) {
	out := Label(Height(2), Anchor("e"), Txt("(123+232)/(123-10)"))
	Grid(out, Columnspan(4), Sticky("e"))
	var b *TButtonWidget
	background := White
	primary := b5.Colors{b5.ButtonText: "#fff", b5.ButtonFace: "#0d6efd", b5.ButtonFocus: "#98c1fe"}
	secondary := b5.Colors{b5.ButtonText: "#fff", b5.ButtonFace: "#6c757d", b5.ButtonFocus: "#c0c4c8"}
	success := b5.Colors{b5.ButtonText: "#fff", b5.ButtonFace: "#198754", b5.ButtonFocus: "#9dccb6"}
	//danger := b5.Colors{b5.ButtonText: "#fff", b5.ButtonFace: "#dc3545", b5.ButtonFocus: "#f0a9b0"}
	//warning := b5.Colors{b5.ButtonText: "#000", b5.ButtonFace: "#ffc107", b5.ButtonFocus: "#ecd182"}
	//info := b5.Colors{b5.ButtonText: "#000", b5.ButtonFace: "#0dcaf0", b5.ButtonFocus: "#85d5e5"}

	StyleThemeUse("default")

	primaryStyle := Style(b5.ButtonStyle("primary.TButton", primary, background, false))
	secondStyle := Style(b5.ButtonStyle("secondary.TButton", secondary, background, false))
	successStyle := Style(b5.ButtonStyle("success.TButton", success, background, false))

	number := ctool.NewSet[string]("1", "2", "3", "4", "5", "6", "7", "8", "9", "0")
	oprs := ctool.NewSet[string]("+", "-", "*", "/")
	for i, c := range "C()/789*456-123+0.=" {
		char := string(c)
		var style Opt
		if number.Contains(char) {
			style = primaryStyle
		} else if oprs.Contains(char) {
			style = successStyle
		} else {
			style = secondStyle
		}

		b = TButton(Txt(char), style,
			Command(
				func() {
					switch c {
					case 'C':
						out.Configure(Txt(""))
					case '=':
						x, err := expr.Eval(out.Txt(), nil)
						if err != nil {
							MessageBox(Icon("error"), Msg(err.Error()), Title("Error"))
							x = ""
						}
						out.Configure(Txt(x))
					default:
						out.Configure(Txt(out.Txt() + char))
					}
				},
			),
			Width(-4))
		Grid(b, Row(i/4+1), Column(i%4), Sticky("news"), Ipadx("1.5m"), Ipady("2.6m"))
	}
	Grid(b, Columnspan(2))
	App.Configure(Padx(0), Pady(0)).Wait()
}
