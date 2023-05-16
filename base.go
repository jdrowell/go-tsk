package main

/*
#cgo LDFLAGS: -ltsk
#include <tsk/libtsk.h>
*/
import "C"

func Version() string {
    cVersion := C.tsk_version_get_str()
    return C.GoString(cVersion)
}
