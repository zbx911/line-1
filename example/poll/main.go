package main

import (
	"fmt"
	"github.com/line-api/line"
	"time"
)

func main() {
	cl, err := line.New()
	if err != nil {

	}
	err = cl.LoginViaV3Token(
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI0MTVmMGJjZS1hMmFlLTQ0ZWQtYjAyYS00YzA0MDA5NzM2YmYiLCJhdWQiOiJMSU5FIiwiaWF0IjoxNjMyODkxOTQwLCJleHAiOjE2MzM0OTY3NDAsInNjcCI6IkxJTkVfQ09SRSIsInJ0aWQiOiIyMWIwMWU3NC00ZTg0LTRkNGMtYTJhNi1iNDNmMWZlYzIzNDYiLCJyZXhwIjoxNzkwNTcxOTQwLCJ2ZXIiOiIzLjEiLCJhaWQiOiJ1ZDFiMWJmMDJiYzM2MmI1MjY2ZjRlMmFmNWJlYTM2OTAiLCJsc2lkIjoiODExNzk0YzktNDc2My00OGJkLWJhMjgtZTRiN2UyY2MwNjA3IiwiZGlkIjoiYmM0ZTY2NzkyMGUyMTFlYzhmMGNhOGExNTkzYzllODgiLCJjdHlwZSI6IkFORFJPSUQiLCJjbW9kZSI6IlBSSU1BUlkiLCJjaWQiOiIwMDAwMDAwMDAwIn0.OO8MvAb1drVOO16D5cQVgUo7nO16Si5CQp2MXq8jVRQ",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJqdGkiOiIyMWIwMWU3NC00ZTg0LTRkNGMtYTJhNi1iNDNmMWZlYzIzNDYiLCJhdGkiOiI0MTVmMGJjZS1hMmFlLTQ0ZWQtYjAyYS00YzA0MDA5NzM2YmYiLCJhdWQiOiJMSU5FIiwicm90IjoiUk9UQVRFIiwiaWF0IjoxNjMyODkxOTQwLCJleHAiOjE3OTA1NzE5NDAsInNjcCI6IkxJTkVfQ09SRSIsInZlciI6IjMuMSIsImFpZCI6InVkMWIxYmYwMmJjMzYyYjUyNjZmNGUyYWY1YmVhMzY5MCIsImxzaWQiOiI4MTE3OTRjOS00NzYzLTQ4YmQtYmEyOC1lNGI3ZTJjYzA2MDciLCJkaWQiOiJiYzRlNjY3OTIwZTIxMWVjOGYwY2E4YTE1OTNjOWU4OCIsImFwcElkIjoiMDAwMDAwMDAwMCJ9.dezyW2xUvURyN7uAwMWXRTP7HqLNWAoknyNE4kmbn2s",
	)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
	for true {
		err := cl.SaveKeeper()
		if err != nil {
			fmt.Printf("%#v\n", err)
		}
		ops, err := cl.FetchLineOperations()
		if err != nil {
			fmt.Printf("%#v\n", err)
		}
		for _, op := range ops {
			fmt.Printf("%#v\n", op)
		}
		time.Sleep(time.Second * 5)
	}
}
