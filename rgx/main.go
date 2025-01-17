package main

import (
	"rgx/cmd/rgx"
	"rgx/common/utils"
)

func main() {
	utils.Config = utils.ReadConfig()
	utils.Configure()
	rgx.Execute()
}
