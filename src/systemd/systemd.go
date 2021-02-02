package systemd

import (
	"runtime"

	"github.com/x13a/init-proxyd/systemd/sd"
)

const listenFdsStart = 3

func Is() bool {
	return runtime.GOOS == "linux"
}

func Sockets() ([]int, error) {
	cnt := sd.ListenFds(true)
	res := make([]int, cnt)
	for i := 0; i < cnt; i++ {
		res[i] = i + listenFdsStart
	}
	return res, nil
}
