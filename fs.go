package main

/*
#cgo LDFLAGS: -ltsk
#include <tsk/libtsk.h>
#include "wrapper.h"
*/
import "C"

import (
    "fmt"
    "unsafe"
)

type Filesystem C.TSK_FS_INFO
type Directory  C.TSK_FS_DIR
type Filename   C.TSK_FS_NAME

func (part *Partition) OpenFilesystem() *Filesystem {
    cFs := C.tsk_fs_open_vol((*C.TSK_VS_PART_INFO)(part), C.TSK_FS_TYPE_DETECT)

    return (*Filesystem)(cFs)
}

func (fs *Filesystem) Close() {
    C.tsk_fs_close((*C.TSK_FS_INFO)(fs))
}

func (fs *Filesystem) Type() int {
    return int(fs.ftype)
}

func (fs *Filesystem) TypeName() string {
    return C.GoString(C.tsk_fs_type_toname(fs.ftype))
}

func (fs *Filesystem) RootInum() int {
    return int(fs.root_inum)
}

func (fs *Filesystem) OpenDirectory(aMeta int) *Directory {
    cDir := C.tsk_fs_dir_open_meta((*C.TSK_FS_INFO)(fs), C.ulong(aMeta))

    return (*Directory)(cDir)
}

//func (fs *Filesystem) DirWalk(aMeta int, aFlags int, aAction *func, aPtr int) {
//func (fs *Filesystem) DirWalk(cb func(fn Filename, path string)) {
func (dir *Directory) Walk() {
    count := 0
    data := unsafe.Pointer(&count)
    //callback := C.callback_func(C.dirwalk_callback)
    // ret := C.tsk_fs_dir_walk((*C.TSK_FS_INFO)(fs), C.ulong(fs.RootInum()), C.TSK_FS_DIR_WALK_FLAG_ALLOC | C.TSK_FS_DIR_WALK_FLAG_RECURSE, (*[0]byte)(C.dirwalk_callback), data)
    ret := C.tsk_fs_dir_walk(dir.fs_info, dir.addr, C.TSK_FS_DIR_WALK_FLAG_ALLOC | C.TSK_FS_DIR_WALK_FLAG_RECURSE, (*[0]byte)(C.dirwalk_callback), data)

    //if ret != C.TSK_ERR_OK {
    if ret != 0 {
        fmt.Println("Error: could not walk directory tree")
        return
    }

    fmt.Printf("Found %d files\n", count)
}

//export go_dirwalk_callback
func go_dirwalk_callback(cFile *C.TSK_FS_FILE, cString *C.char, data unsafe.Pointer) C.int {
    // prog := "-\\|/"
    // count_ptr := (*int)(data)
    // i := (*count_ptr) % 4
    //fmt.Printf("%c%c", prog[i], 8)
    name := cFile.name
    if name._type == C.TSK_FS_NAME_TYPE_REG {
        fmt.Printf("%s%s\n", C.GoString(cString), C.GoString(cFile.name.name))
    }

    return 0;
}

func (dir *Directory) Close() {
    C.tsk_fs_dir_close((*C.TSK_FS_DIR)(dir))
}

func (dir *Directory) Size() int {
    return int(C.tsk_fs_dir_getsize((*C.TSK_FS_DIR)(dir)))
}

func (dir *Directory) GetName(n int) *Filename {
    cFn := C.tsk_fs_dir_get_name((*C.TSK_FS_DIR)(dir), C.ulong(n))
    return (*Filename)(cFn)
}

func (dir *Directory) Show() {
    fmt.Printf("    root dir: %d\n", dir.Size())
    // ((*Filesystem)(dir.fs_info)).DirWalk()
    dir.Walk()
    // for i := 0; i < dir.Size(); i += 1 {
    //     fn := dir.GetName(i)
    //     fmt.Printf("    %v\n", fn.Name())
    // }
}

func (fn *Filename) Name() string {
    return C.GoString(fn.name)
}


/*
func (vs *VolumeSystem) BlockSize() int {
    return int(vs.block_size)
}

func (vs *VolumeSystem) PartCount() int {
    return int(vs.part_count)
}

func (vs *VolumeSystem) TypeDescr() string {
    return C.GoString(C.tsk_vs_type_todesc(vs.vstype))
}

func (vs *VolumeSystem) GetPartition(n int) *PartitionInfo {
    cPi := C.tsk_vs_part_get((*C.TSK_VS_INFO)(vs), C.uint(n))
    return (*PartitionInfo)(cPi)
}

func (pi *PartitionInfo) Descr() string {
    return C.GoString(pi.desc)
}

func (pi *PartitionInfo) Len() int {
    return int(pi.len)
}

func (pi *PartitionInfo) SlotNumber() int {
    return int(pi.slot_num)
}
*/

