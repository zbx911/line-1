package phone

import "fmt"

type Manual struct {
	number string
}

func NewManual() *Manual {
	return &Manual{}
}
func (m Manual) CancelNumber() {
	//	No-op
}

func (m Manual) GetNumber() (string, error) {
	m.number = getInput("phone number:")
	return m.number, nil
}

func (m Manual) GetCode() (string, error) {
	return getInput(fmt.Sprintf("code for %v:", m.number)), nil
}

func (m Manual) BanNumber() {
	fmt.Printf("wrong number: %v", m.number)
}

func getInput(msg string) string {
	fmt.Print(msg)
	var first string
	fmt.Scanln(&first)
	return first
}
