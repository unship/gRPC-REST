package cert

//go:generate make generate

var (
	Key, _  = Asset("server.key")
	Cert, _ = Asset("server.pem")
)
