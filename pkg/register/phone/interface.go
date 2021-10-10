package phone

type Service interface {
	GetNumber() (string, error)
	GetCode() (string, error)
	BanNumber()
	CancelNumber()
}
