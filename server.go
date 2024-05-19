package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func connect() (net.Conn, error) {
	fmt.Println("Listening on port :6379")

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		return nil, err
	}

	conn, err := l.Accept()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func server() error {
	conn, err := connect()
	if err != nil {
		return err
	}

	r := bufio.NewReader(conn)

	valueCh := make(chan Value)

	go func() {
		defer close(valueCh)

		for {
			resp := NewResp(r)
			value, err := resp.Read()

			if err != nil {
				fmt.Println("Error reading value: ", err)
				return
			}

			if value.typ != "array" {
				fmt.Println("Invalid request, expected array")
				continue
			}

			if len(value.array) == 0 {
				fmt.Println("Invalid request, expected array length > 0")
				continue
			}

			valueCh <- value
		}
	}()

	for value := range valueCh {
		writer := NewWriter(conn)
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			err := writer.Write(Value{typ: "string", str: ""})
			if err != nil {
				fmt.Println("Error writing value to client: ", err)
			}
			continue
		}

		result := handler(args)
		err := writer.Write(result)
		if err != nil {
			fmt.Println("Error writing value to client: ", err)
		}
	}

	return nil
}
