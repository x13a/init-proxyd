// +build !darwin

package sockets

import "github.com/x13a/init-proxyd/sockets/systemd"

func Get(_ string) ([]int, error) {
	return systemd.Sockets()
}
