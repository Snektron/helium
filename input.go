package helium

import (
	"os"
	"bufio"
)

const (
	EOF = -1
	EOL = '\n'
	TAB = '\t'
	RETURN = '\r'
	FEED = '\f'
	VTAB = '\v'
	SPACE = ' '
)

type Input interface {
	Get(uint) rune
}

type BufferedInput struct {
	buffer []rune
}

func BufferedFileInput(filename string) (Input, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	buf := make([]rune, 0)

	for {
		r, _, e := reader.ReadRune()
		if e != nil {
			break 
		}

		buf = append(buf, r)
	}

	return &BufferedInput{buf}, nil
}

func (bi *BufferedInput) Get(pos uint) rune {
	if pos >= uint(len(bi.buffer)) {
		return EOF
	}

	return bi.buffer[pos]
}