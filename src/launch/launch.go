package launch

/*
#include <errno.h>
#include <stdlib.h>
int launch_activate_socket(const char *name, int **fds, size_t *cnt);
*/
import "C"
import (
	"errors"
	"strconv"
	"unsafe"
)

var (
	ErrNotExist         = errors.New("The socket name specified does not exist in the caller's launchd.plist")
	ErrNotManaged       = errors.New("The calling process is not managed by launchd")
	ErrAlreadyActivated = errors.New("The specified socket has already been activated")
)

func ActivateSocket(name string) ([]int, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var cFds *C.int
	cCnt := C.size_t(0)
	if err := C.launch_activate_socket(cName, &cFds, &cCnt); err != 0 {
		switch err {
		case C.ENOENT:
			return nil, ErrNotExist
		case C.ESRCH:
			return nil, ErrNotManaged
		case C.EALREADY:
			return nil, ErrAlreadyActivated
		default:
			return nil, errors.New("Unknown error: " + strconv.Itoa(int(err)))
		}
	}
	ptr := unsafe.Pointer(cFds)
	defer C.free(ptr)
	cnt := int(cCnt)
	fds := (*[1 << 30]C.int)(ptr)[:cnt:cnt]
	res := make([]int, cnt)
	for i := 0; i < cnt; i++ {
		res[i] = int(fds[i])
	}
	return res, nil
}
