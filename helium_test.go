package helium

import (
	_ "fmt"
	"testing"
)

type test func()

func TestHelium(t *testing.T) {
	bfi, err := BufferedFileInput("test.txt")

	if err != nil {
		t.Error(err)
		return
	}

	ctx := NewContext(bfi)

	var a Rule
	var b Rule

	a = Sequence(
		Rune('a'),
		Recursive(&b))

	b = Rune('b')

	/*
	root :=
		Sequence(
			Capture(
				Sequence(
					Optional(
						AnyRuneOf('a', 'b', 'c')),
					Rune('d')),
				func(text string) {
					fmt.Printf("Captured: %q\n", text)
				}),
			Rune(EOF)) */

	result := ctx.Parse(a)

	if !result {
		t.Error(ctx.Error())
	}
}