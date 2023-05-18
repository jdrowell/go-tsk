package tsk

/*
#cgo LDFLAGS: -ltsk
#include <tsk/libtsk.h>
*/
import "C"

import (
    "fmt"
)

func Version() string {
    cVersion := C.tsk_version_get_str()
    return C.GoString(cVersion)
}

func ErrorNo() int {
    errno := int(C.tsk_error_get_errno())
    return errno
}

func ErrorStr() string {
    errstr := C.GoString(C.tsk_error_get_errstr())
    return errstr
}

func Error() error {
    return fmt.Errorf("%s (%d)", ErrorStr(), ErrorNo())
}
