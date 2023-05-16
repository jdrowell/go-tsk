package main

/*
#cgo LDFLAGS: -ltsk
#include <tsk/libtsk.h>
*/
import "C"

import (
    "fmt"
)

type VolumeSystem  C.TSK_VS_INFO
type Partition     C.TSK_VS_PART_INFO

func (di *DiskImage) OpenVolumeSystem() (*VolumeSystem, error) {
    cVs, err := C.tsk_vs_open((*C.TSK_IMG_INFO)(di), 0, C.TSK_VS_TYPE_DETECT)

    return (*VolumeSystem)(cVs), err
}

func (vs *VolumeSystem) Close() {
    C.tsk_vs_close((*C.TSK_VS_INFO)(vs))
}

func (vs *VolumeSystem) BlockSize() int {
    return int(vs.block_size)
}

func (vs *VolumeSystem) PartCount() int {
    return int(vs.part_count)
}

func (vs *VolumeSystem) Type() int {
    return int(vs.vstype)
}

func (vs *VolumeSystem) TypeDescr() string {
    return C.GoString(C.tsk_vs_type_todesc(vs.vstype))
}

func (vs *VolumeSystem) TypeName() string {
    return C.GoString(C.tsk_vs_type_toname(vs.vstype))
}

func (vs *VolumeSystem) GetPartition(n int) *Partition {
    cPi := C.tsk_vs_part_get((*C.TSK_VS_INFO)(vs), C.uint(n))
    return (*Partition)(cPi)
}

func (part *Partition) Descr() string {
    return C.GoString(part.desc)
}

func (part *Partition) Len() int {
    return int(part.len)
}

func (part *Partition) SlotNumber() int {
    return int(part.slot_num)
}

func (part *Partition) Show(volNo int) {
    vs := (*VolumeSystem)(part.vs)

    fmt.Printf("  p[%d][%d]: %s (%d bytes)\n", volNo, part.SlotNumber(), part.Descr(), part.Len() * vs.BlockSize())
    fs := part.OpenFilesystem()
    if fs == nil {
        fmt.Printf("    unknown filesystem\n")
        return
    }
    defer fs.Close()

    fmt.Printf("    filesystem type: %s (%d)\n", fs.TypeName(), fs.Type())
    dir := fs.OpenDirectory(fs.RootInum())
    if fs == nil {
        fmt.Printf("    couldn't open root directory\n")
        return
    }
    defer dir.Close()

    dir.Show()
}

