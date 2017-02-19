package helium

import (
	"fmt"
	"testing"
)

func TestHelium(t *testing.T) {
	bfi, err := BufferedFileInput("test.txt")

	if err != nil {
		t.Error(err)
		return
	}

	ctx := NewContext(bfi)

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
			Rune(EOF))

	result := ctx.parse(root)

	if !result {
		t.Error(ctx.Error())
	}
}