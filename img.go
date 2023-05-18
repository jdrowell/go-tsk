package tsk

/*
#cgo LDFLAGS: -ltsk
#include <tsk/libtsk.h>
*/
import "C"

import (
    "fmt"
)

type DiskImage C.TSK_IMG_INFO

func OpenImage(filename string) (*DiskImage, error) {
    //tsk_img_open, err := C.tsk_img_open(C.CString(filename), C.TSK_IMG_TYPE_ENUM.TSK_IMG_TYPE_DETECT, 0)
    cImg, err := C.tsk_img_open_sing(C.CString(filename), C.TSK_IMG_TYPE_DETECT, 0)

    return (*DiskImage)(cImg), err
}

func (di *DiskImage) Close() {
    C.tsk_img_close((*C.TSK_IMG_INFO)(di))
}

func (di *DiskImage) Size() int {
    return int(di.size)
}

func (di *DiskImage) Type() int {
    return int(di.itype)
}

func (di *DiskImage) TypeDescr() string {
    return C.GoString(C.tsk_img_type_todesc(di.itype))
}

func (di *DiskImage) TypeName() string {
    return C.GoString(C.tsk_img_type_toname(di.itype))
}

func (di *DiskImage) Show() {
    fmt.Printf("img size: %v\n", di.Size())
    fmt.Printf("img type: %s '%s' (%d)\n", di.TypeName(), di.TypeDescr(), di.Type())

    vs, err := di.OpenVolumeSystem()
    if err != nil {
        panic(err)
    }
    if vs == nil {
        panic("vs is nil")
    }
    defer vs.Close()

    fmt.Printf("vs partitions: %v, block size: %d\n", vs.PartCount(), vs.BlockSize())
    fmt.Printf("vs type: %s '%s' (%d)\n", vs.TypeName(), vs.TypeDescr(), vs.Type())
    for i := 0; i < int(vs.part_count); i += 1 {
        part := vs.GetPartition(i)
        part.Show(i)
    }
}

