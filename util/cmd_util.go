package util

import "os/exec"

var EnableGoTidy bool

func DoGoGet(packageName string) error {
	cmd := exec.Command("go", "-u", packageName)
	return cmd.Run()
}

func DoGoModTidy() error {
	cmd := exec.Command("go", "mod", "tidy")
	return cmd.Run()
}

func DoGoFmt(filepath string) error {
	cmd := exec.Command("gofmt", "-w", filepath)
	return cmd.Run()
}
