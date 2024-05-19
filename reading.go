package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	bulk  string
	array []Value
}

type Resp struct {
	scanner *bufio.Scanner
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{scanner: bufio.NewScanner(rd)}
}

func (r *Resp) readArray(len int) (Value, error) {
	v := Value{}
	v.typ = "array"
	v.array = make([]Value, 0)

	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		v.array = append(v.array, val)
	}

	return v, nil
}

func (r *Resp) readLine() (line string, err error) {
	isOver := r.scanner.Scan()

	if !isOver {
		return "", fmt.Errorf("EOF")
	}

	tok := r.scanner.Text()

	return tok, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"

	line, err := r.readLine()
	if err != nil {
		return v, err
	}

	v.bulk = line

	return v, nil
}

func (r *Resp) readMetadata() (typ byte, length int, err error) {
	line, err := r.readLine()

	if err != nil {
		return byte(0), 0, err
	}

	i64, err := strconv.ParseInt(line[1:], 10, 64)
	if err != nil {
		return byte(0), 0, err
	}

	return line[0], int(i64), nil
}

func (r *Resp) Read() (Value, error) {
	_type, len, err := r.readMetadata()

	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray(len)
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Invalid type: %v", string(_type))
		return Value{}, nil
	}
}
