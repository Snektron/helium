package helium

import (
	"bytes"
	"fmt"
)

type capture struct {
	first, last uint
}

type Context struct {
	input Input
	pos, lasterr uint
	posStack *Stack
	cap capture
}

func NewContext(in Input) *Context {
	return &Context{in, 0, 0, NewStack(), capture{0, 0}}
}

func (c *Context) parse(rule Rule) bool {
	c.posStack.Push(c.pos)
	res := rule(c)

	c.cap.first = c.posStack.Peek().(uint)
	c.cap.last = c.pos

	if res {
		c.posStack.Pop()
	} else {
		c.lasterr = c.pos
		c.pos = c.posStack.Pop().(uint)
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
	return 0
}

func (c *Context) Peek() rune {
	return c.input.Get(c.pos)
}

func (c *Context) PeekAhead(n int) rune {
	return c.input.Get(c.pos + uint(n))
}

func (c *Context) Consume() rune {
	res := c.Peek()
	c.pos++
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
	if c.Peek() == EOF {
		return fmt.Sprintf("%d: Unexpected end of file.", c.Line())
	} else {
		return fmt.Sprintf("%d: Unexpected input near '%s'.", c.Line(), c.errLine())
	} 
}

func (c *Context) errLine() string {
	var buffer bytes.Buffer

	for i := c.lasterr; c.input.Get(i) != EOL && c.input.Get(i) != EOF; i++ {
		buffer.WriteRune(c.input.Get(i))
	}

	return buffer.String()
}
