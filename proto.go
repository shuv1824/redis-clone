package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/tidwall/resp"
)

const (
	CommandSET = "SET"
)

type Command interface {
}

type SetCommand struct {
	key, value string
}

func ParseCommand(raw string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Read %s\n", v.Type())
		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				// fmt.Printf("  #%d %s, value: '%s'\n", i, v.Type(), v)
				switch value.String() {
				case CommandSET:
					fmt.Println(len(v.Array()))
					if len(v.Array()) != 3 {
						return "", fmt.Errorf("invalid number of variables for SET command")
					}
					cmd := SetCommand{
						key:   v.Array()[1].String(),
						value: v.Array()[2].String(),
					}

					return cmd, nil
				}
			}
		}
	}
	return "", nil
}
