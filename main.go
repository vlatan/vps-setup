package main

import (
	"github.com/vlatan/vps-setup/internal/setup"
	"github.com/vlatan/vps-setup/internal/utils"
)

func main() {
	s, err := setup.New()
	if err != nil {
		utils.Exit(err)
	}

	defer s.Close()
	s.Run()
}
