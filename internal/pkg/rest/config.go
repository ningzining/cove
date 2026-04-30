package rest

const (
	DevMode  = "dev"
	ProdMode = "prod"
)

type Config struct {
	Name string
	Host string
	Port int
	Mode string
}
