package main

import (
	"os"
	"soloos/common/util"
)

func main() {
	var (
		env     Env
		options Options
		err     error
	)

	optionsFile := os.Args[1]

	options, err = LoadOptionsFile(optionsFile)
	util.AssertErrIsNil(err)

	util.AssertErrIsNil(env.Init(options))
	util.AssertErrIsNil(env.Serve())
}