package main

import (
	"github.com/vlatan/vps-setup/internal/setup"
)

func main() {
	s := setup.New()
	defer s.Close()
	s.Run()
}
