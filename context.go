package helium

import (
	"bytes"
	"fmt"
)

type state struct {
	pos, line uint
}

type capture struct {
	first, last uint
}

type Context struct {
	input Input
	state, err state
	stateStack *Stack
	cap capture
	filter Rule
	filtering bool
}

func NewContext(in Input) *Context {
	return &Context{input: in, stateStack: NewStack()}
}

func (c *Context) SetFilter(filter Rule) {
	c.filter = filter
}

func (c *Context) PushState() {
	c.stateStack.Push(c.state)
}

func (c *Context) PopState() {
	c.state = c.stateStack.Pop().(state)
}

func (c *Context) DiscardState() {
	c.stateStack.Pop()
}

func (c *Context) Parse(rule Rule) bool {
	if !c.filtering && c.filter != nil {
		c.filtering = true
		c.Parse(c.filter)
		c.filtering = false
	}

	if rule == nil {
		panic(fmt.Errorf("Parse: rule is nil!"))
	}

	c.PushState()
	res := rule(c)

	c.cap.first = c.stateStack.Peek().(state).pos
	c.cap.last = c.state.pos

	if res {
		c.DiscardState()
	} else {
		if c.state.pos > c.err.pos {
			c.err = c.state
		}

		c.PopState()
	}
	
	return res
}

func (c *Context) capture() string {
	l := c.cap.last - c.cap.first

	buffer := make([]rune, l)

	for i := uint(0); i < l; i++ {
		buffer[i] = c.input.Get(c.cap.first + i)
	}

	return string(buffer)
}

func (c *Context) Line() uint {
	return c.state.line
}

func (c *Context) Peek() rune {
	return c.input.Get(c.state.pos)
}

func (c *Context) PeekAhead(n int) rune {
	return c.input.Get(c.state.pos + uint(n))
}

func (c *Context) Consume() rune {
	res := c.Peek()
	c.state.pos++

	if res == EOL {
		c.state.line++
	}

	return res
}

func (c *Context) Skip(n int) {
	for i := 0; i < n; i++ {
		c.Consume()
	}
}

func (c *Context) TestString(text string) bool {
	for i := 0; i < len(text); i++ {
		if rune(text[i]) != c.PeekAhead(i) {
			return false
		}
	}
	return true
}

func (c *Context) Error() string {
	if c.input.Get(c.err.pos) == EOF {
		return fmt.Sprintf("%d: Unexpected end of file.", c.Line())
	} else {
		return fmt.Sprintf("%d: Unexpected input near '%s'.", c.Line(), c.ErrLine())
	} 
}

func (c *Context) ErrLine() string {
	var buffer bytes.Buffer

	for i := c.err.pos; c.input.Get(i) != EOL && c.input.Get(i) != EOF; i++ {
		buffer.WriteRune(c.input.Get(i))
	}

	return buffer.String()
}
