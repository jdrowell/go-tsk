package tsk

/*
#cgo LDFLAGS: -ltsk
#include <tsk/libtsk.h>
*/
import "C"

import (
    "fmt"
)

// Version returns a text representation of the libtsk version.
// It wraps tsk_version_get_str().
func Version() string {
    cVersion := C.tsk_version_get_str()
    return C.GoString(cVersion)
}

// ErrNo returns the libtsk specific errno.
// It wraps tsk_error_get_errno().
func ErrorNo() int {
    errno := int(C.tsk_error_get_errno())
    return errno
}

// ErrorStr return a description of the current error.
// It wraps tsk_error_get_errstr().
func ErrorStr() string {
    errstr := C.GoString(C.tsk_error_get_errstr())
    return errstr
}

// Error returns a string containing the error's description and number (errno).
func Error() error {
    return fmt.Errorf("%s (%d)", ErrorStr(), ErrorNo())
}
