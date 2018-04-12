package main

import (
	"recho/utils"
	"recho/handlers"
	"recho/validators"
)

func main() {
	s := utils.InitEnv("./routes.toml")
	s.RegisterHandler(&handlers.User{})
	s.RegisterValidator(&validators.User{})
	//utils.Pr(s)
}
