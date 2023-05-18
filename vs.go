package tsk

/*
#cgo LDFLAGS: -ltsk
#include <tsk/libtsk.h>
*/
import "C"

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

func (vs *VolumeSystem) Offset() int {
    return int(vs.offset)
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

func (part *Partition) Start() int {
    return int(part.start)
}

func (part *Partition) Len() int {
    return int(part.len)
}

func (part *Partition) TableNumber() int {
    return int(part.table_num)
}

func (part *Partition) SlotNumber() int {
    return int(part.slot_num)
}

