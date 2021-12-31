package main

import (
	"publicEmail/config"
	"publicEmail/internal/smtp_server"
)

func main() {
	config.LoadFile("")

	smtp_server.RunSMTPServer()
}
