package helium

type Rule func(*Context) bool
type StringConsumer func(string)
type BoolConsumer func(bool)

func Rune(r rune) Rule {
	return func(ctx *Context) bool {
		if ctx.Peek() == r {
			ctx.Consume()
			return true
		}
		return false
	}
}

func AnyRuneOf(runes ...rune) Rule {
	return func(ctx *Context) bool {
		for _, r := range runes {
			if ctx.Peek() == r {
				ctx.Consume()
				return true
			}
		}
		return false
	}
}

func RuneRange(low, high rune) Rule {
	return func(ctx *Context) bool {
		if ctx.Peek() >= low && ctx.Peek() <= high {
			ctx.Consume()
			return true
		}
		return false
	}
}

func String(text string) Rule {
	return func(ctx *Context) bool {
		if ctx.TestString(text) {
			ctx.Skip(len(text))
			return true
		}
		return false
	}
}

func AnyStringOf(strings ...string) Rule {
	return func(ctx *Context) bool {
		for _, text := range strings {
			if ctx.TestString(text) {
				ctx.Skip(len(text))
				return true
			}
		}
		return false
	}
}

func Any() Rule {
	return func(ctx *Context) bool {
		if ctx.Peek() != EOF {
			ctx.Consume()
			return true
		}
		return false
	}
}

func ZeroOrMore(rule Rule) Rule {
	return func(ctx *Context) bool {
		for ctx.parse(rule) {}
		return true
	}
}

func OneOrMore(rule Rule) Rule {
	return func(ctx *Context) bool {
		res := ctx.parse(rule)
		for ctx.parse(rule) {}
		return res
	}
}

func Optional(rule Rule) Rule {
	return func(ctx *Context) bool {
		ctx.parse(rule)
		return true
	}
}

func Sequence(rules ...Rule) Rule {
	return func(ctx *Context) bool {
		for _, rule := range rules {
			if !ctx.parse(rule) {
				return false
			}
		}
		return true
	}
}

func FirstOf(rules ...Rule) Rule {
	return func(ctx *Context) bool {
		for _, rule := range rules {
			if ctx.parse(rule) {
				return true
			}
		}
		return false
	}
}

func Capture(rule Rule, consumer StringConsumer) Rule {
	return func(ctx *Context) bool {
		if ctx.parse(rule) {
			consumer(ctx.capture())
			return true
		}
		return false
	}
}

func Action(rule Rule, consumer BoolConsumer) Rule {
	return func(ctx *Context) bool {
		res := ctx.parse(rule)
		consumer(res)
		return res
	}	
}