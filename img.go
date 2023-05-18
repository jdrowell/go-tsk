package tsk

/*
#cgo LDFLAGS: -ltsk
#include <tsk/libtsk.h>
*/
import "C"

// DiskImage wraps the TSK_IMG_INFO struct, providing direct access to its
// members, and groups related functions to disk image handling.
type DiskImage C.TSK_IMG_INFO

// OpenImage opens a single disk image for processing.
// It wraps tsk_img_open_sing().
func OpenImage(filename string) (*DiskImage, error) {
    cimg, err := C.tsk_img_open_sing(C.CString(filename), C.TSK_IMG_TYPE_DETECT, 0)

    return (*DiskImage)(cimg), err
}

// Close closes a disk image.
// It wraps tsk_img_close().
func (di *DiskImage) Close() {
    C.tsk_img_close((*C.TSK_IMG_INFO)(di))
}

// Size returns the size of the disk image in bytes.
func (di *DiskImage) Size() int {
    return int(di.size)
}

// Type returns the disk image type.
func (di *DiskImage) Type() int {
    return int(di.itype)
}

// TypeDescr returns a description of the disk image's type.
// It wraps tsk_img_type_todesc().
func (di *DiskImage) TypeDescr() string {
    return C.GoString(C.tsk_img_type_todesc(di.itype))
}

// TypeName returns the name for a disk image's type.
// It wraps tsk_img_type_toname().
func (di *DiskImage) TypeName() string {
    return C.GoString(C.tsk_img_type_toname(di.itype))
}

