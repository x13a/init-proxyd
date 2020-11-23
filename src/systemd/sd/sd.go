package sd

import (
	"os"
	"strconv"
)

const (
	envListenPid     = "LISTEN_PID"
	envListenFds     = "LISTEN_FDS"
	envListenFdnames = "LISTEN_FDNAMES"
)

func ListenFds(unsetEnvironment bool) int {
	if unsetEnvironment {
		defer func() {
			os.Unsetenv(envListenPid)
			os.Unsetenv(envListenFds)
			os.Unsetenv(envListenFdnames)
		}()
	}
	if pid, err := strconv.Atoi(os.Getenv(envListenPid)); err != nil ||
		pid != os.Getpid() {

		return 0
	}
	res, _ := strconv.Atoi(os.Getenv(envListenFds))
	return res
}
