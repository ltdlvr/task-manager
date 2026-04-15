package config

import (
	"os"
)

type Rest struct {
	srvHost string
	srvPort string
	dbURL   string
}

func NewRest() *Rest {
	conf := &Rest{
		dbURL: getDatabaseURL(),
	}

	host := os.Getenv("REST_SERVER_HOST")
	if host == "" {
		fatal("REST server host", "REST_SERVER_HOST", host)
	}
	conf.srvHost = host

	port := os.Getenv("REST_SERVER_PORT")
	if !isValidPort(port) {
		fatal("REST server port", "REST_SERVER_PORT", port)
	}
	conf.srvPort = port

	return conf
}

func (c *Rest) ServerHost() string {
	return c.srvHost
}

func (c *Rest) ServerPort() string {
	return c.srvPort
}

func (c *Rest) DatabaseURL() string {
	return c.dbURL
}
