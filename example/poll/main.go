package main

import (
	"fmt"
	"github.com/line-api/line"
)

func main() {
	cl, err := line.New()
	if err != nil {

	}
	for true {
		ops, err := cl.FetchLineOperations()
		if err != nil {
			fmt.Printf("%#v\n", err)
		}
		for _, op := range ops {
			fmt.Printf("%#v\n", op)
		}
	}
}
