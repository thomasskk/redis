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
	num   int
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
	v.array = make([]Value, len)

	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		v.array = append(v.array, val)
	}

	return v, nil
}

func (r *Resp) readLine() (line []byte, err error) {
	r.scanner.Scan()
	tok := r.scanner.Bytes()

	return tok, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"

	line, err := r.readLine()
	if err != nil {
		return v, err
	}

	v.bulk = string(line)

	return v, nil
}

func (r *Resp) readMetadata() (typ byte, len int, err error) {
	line, err := r.readLine()

	if err != nil {
		return byte(0), 0, err
	}

	i64, err := strconv.ParseInt(string(line[1:]), 10, 64)
	if err != nil {
		return typ, 0, err
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
		return r.readBulk(len)
	default:
		fmt.Printf("Invalid type: %v", string(_type))
		return Value{}, nil
	}
}
