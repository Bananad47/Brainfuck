package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	size       = 30000
	plus       = '+'
	minus      = '-'
	next       = '>'
	previous   = '<'
	loopStart  = '['
	loopEnd    = ']'
	getValue   = ','
	printValue = '.'
)

var indexError = errors.New("array index out of bounds")
var unknownCommand = errors.New("unknown command")

type Program struct {
	arr   [30000]byte
	index int
}

func (p *Program) Plus() {
	p.arr[p.index]++
}

func (p *Program) Minus() {
	p.arr[p.index]--
}

func (p *Program) Next() error {
	if p.index+1 == size {
		return indexError
	}
	p.index++
	return nil
}

func (p *Program) Previous() error {
	if p.index == 0 {
		return indexError
	}
	p.index--
	return nil
}

func (p *Program) Print() error {
	_, err := fmt.Printf("%c", p.arr[p.index])
	return err
}

func (p *Program) Get() error {
	_, err := fmt.Scan(&p.arr[p.index])
	if err != nil {
		return err
	}
	return nil
}

func (p *Program) Run(s, e string, getNextToken func() (byte, bool)) error {
	isLoop := false
	loopIndex := 0
	var b strings.Builder
	for token, ok := getNextToken(); ok; token, ok = getNextToken() {
		if isLoop && token != loopEnd {
			b.WriteByte(token)
		}
		switch token {
		case plus:
			p.Plus()
		case minus:
			p.Minus()
		case next:
			err := p.Next()
			if err != nil {
				return err
			}
		case previous:
			err := p.Previous()
			if err != nil {
				return err
			}
		case getValue:
			err := p.Get()
			if err != nil {
				return err
			}
		case printValue:
			err := p.Print()
			if err != nil {
				return err
			}
		case loopStart:
			isLoop = true
			loopIndex = p.index
		case loopEnd:
			isLoop = false
			_idx := 0
			_str := []byte(b.String())
			b = strings.Builder{}
			fn := func() (byte, bool) {
				if _idx == len(_str) {
					return 0, false
				}
				_idx++
				return _str[_idx-1], true
			}
			for p.arr[loopIndex] != 0 {
				err := p.Run(s, e, fn)
				if err != nil {
					return err
				}
				_idx = 0
			}
		default:
			if token > 20 {
				return fmt.Errorf("%w: %s", unknownCommand, string(token))
			}
		}
	}
	return nil
}

func Shell() {
	p := Program{}
	rd := bufio.NewReader(os.Stdin)
	arr := make([]byte, 1)
	fn := func() (byte, bool) {
		_, err := rd.Read(arr)
		if err != nil && err != io.EOF {
			return 0, false
		}
		return arr[0], true
	}
	for {
		err := p.Run(">> ", "\n", fn)
		println("Error: ", err)
	}
}

func main() {
	var file io.ReadCloser
	if len(os.Args) < 2 {
		Shell()
	} else {
		var err error
		file, err = os.Open(os.Args[1])
		if err != nil {
			fmt.Printf("can't open file %s", os.Args[1])
			return
		}
	}
	rd := bufio.NewReader(file)
	arr := make([]byte, 1)
	fn := func() (byte, bool) {
		_, err := rd.Read(arr)
		if err != nil {
			return 0, false
		}
		return arr[0], true
	}
	p := Program{}
	err := p.Run("", "\n", fn)
	if err != nil {
		println("ERROR:", err)
	}
}
