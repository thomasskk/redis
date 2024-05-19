package main

import "sync"

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	store.Set(key, value)

	return Value{typ: "string", str: "OK"}
}
