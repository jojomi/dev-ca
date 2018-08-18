package main

import "os"

func fail(err error) {
	_, _ = os.Stderr.WriteString(err.Error())
	os.Exit(1)
}

func checkFail(err error) {
	if err == nil {
		return
	}
	fail(err)
}
