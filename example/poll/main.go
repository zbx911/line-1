package main

import (
	"fmt"
	"github.com/line-api/line"
	"github.com/line-api/model/go/model"
	"time"
)

func main() {
	cl := line.New(line.Proxy(""))
	//err := cl.LoginViaV3Token(
	//	"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI5YmRmMTIwNy0zNGM0LTRhZDAtOWVlNC0wNGY5MWZhZjUyZWEiLCJhdWQiOiJMSU5FIiwiaWF0IjoxNjMzMzE0NzYwLCJleHAiOjE2MzM5MTk1NjAsInNjcCI6IkxJTkVfQ09SRSIsInJ0aWQiOiJlZThhMmMwZC1iNDFkLTQ0NWUtYThmOS0yYjQyYzY5YjEwOWMiLCJyZXhwIjoxNzkwOTk0NzYwLCJ2ZXIiOiIzLjEiLCJhaWQiOiJ1MjI2MDY2ODk5NDM2MDVmZTk5ODg4MWY3ZGE0M2NlMGIiLCJsc2lkIjoiNTU4NmQ2Y2MtNGU1NC00Y2IxLWI1MDktOGI3MWM0ZDEzMTdiIiwiZGlkIjoiNDUwOTdhMmIyNGJiMTFlYzk0N2VhOGExNTkzYzllODgiLCJjdHlwZSI6IkFORFJPSUQiLCJjbW9kZSI6IlBSSU1BUlkiLCJjaWQiOiIwMDAwMDAwMDAwIn0.Mi_a5TdIDfQ7x7S14tge1Y2RS9uOkUhApxnjXkM8UNo",
	//	"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI2NmY3MWYwMS1iZWVlLTRjNzUtODY3OS02ZDc0ZjlhZjY3MmUiLCJhdGkiOiI4ZTIyOGRhYS1iMzQyLTRkNmItYjYyNC02Y2MzYzkxZDBmZTEiLCJhdWQiOiJMSU5FIiwicm90IjoiUk9UQVRFIiwiaWF0IjoxNjMzNTQyMjUyLCJleHAiOjE3OTEyMjIyNTIsInNjcCI6IkxJTkVfQ09SRSIsInZlciI6IjMuMSIsImFpZCI6InUzNDQwZTQxNzRkZDcxNDcwODk2ODRjMWQyZTdmODQyNSIsImxzaWQiOiIyZGFiZmIwNi1kMzk0LTRlYjYtODBmMy05MGI3YzhiZDEyMzYiLCJkaWQiOiJmMmVlODcwMzI2Y2MxMWVjYjgzMmE4YTE1OTNjOWU4OCIsImFwcElkIjoiMDAwMDAwMDAwMCJ9.P3qGN4DaPtBnG2q8xyWjvWN5o7BXAz2SLScb7MngW74",
	//)
	err := cl.LoginViaAuthKey("u07000fb16ec97ac70a3decb5b6cad1f7:hcIyDQFDITd9tDpk7xmk")
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
	for true {
		err := cl.SaveKeeper()
		if err != nil {
			fmt.Printf("%#v\n", err)
		}
		ops, err := cl.FetchLineOperationsTMCP()
		if err != nil {
			fmt.Printf("%#v\n", err)
		}
		for _, op := range ops {
			if op.OpType == model.OpType_RECEIVE_MESSAGE {
				fmt.Printf("%#v\n", op.Message.String())
				switch op.Message.Text {
				case "hi":
					cl.SendMessageCompact(&model.Message{
						To:   op.Message.To,
						From: cl.Profile.Mid,
					})
				}
			}
		}
		time.Sleep(time.Second * 5)
	}
}
