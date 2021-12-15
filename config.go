package ssr

type Config struct {
	HTTP struct {
		Addr string
	}

	DB struct {
		Host string
		Port int
		User string
		Password string
		Name string
	}
}
