// +build !darwin

package main

import "github.com/x13a/init-proxyd/systemd"

func Sockets(_ string) ([]int, error) {
	return systemd.Sockets()
}
