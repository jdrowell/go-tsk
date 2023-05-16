package main

/*
#cgo LDFLAGS: -ltsk
#include <tsk/libtsk.h>
#include "wrapper.h"
*/
import "C"

import (
    "errors"
    "fmt"
    "unsafe"
)

type Filesystem C.TSK_FS_INFO
type Directory  C.TSK_FS_DIR
type Filename   C.TSK_FS_NAME
type File       C.TSK_FS_FILE

type WalkCallback func (file *File, path string)
var walkCBs = make(map[int]WalkCallback)

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

func my_callback(v string) {
    fmt.Println("hello", v)
}

func (dir *Directory) Walk(cb WalkCallback) error {
    k := int(dir.addr)
    walkCBs[k] = cb
    data := unsafe.Pointer(&k)
    ret := C.tsk_fs_dir_walk(dir.fs_info, dir.addr, C.TSK_FS_DIR_WALK_FLAG_ALLOC | C.TSK_FS_DIR_WALK_FLAG_RECURSE, (*[0]byte)(C.dirwalk_callback), data)

    //if ret != C.TSK_ERR_OK {
    if ret != 0 {
        return fmt.Errorf("%w", errors.New("Could not walk the directory tree"))
    }

    return nil
}

//export go_dirwalk_callback
func go_dirwalk_callback(cFile *C.TSK_FS_FILE, cString *C.char, data unsafe.Pointer) C.int {
    k := (*int)(data)
    cb := walkCBs[*k]
    cb((*File)(cFile), C.GoString(cString))

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
    cb := func (file *File, path string) {
        name := file.name
        if name._type == C.TSK_FS_NAME_TYPE_REG {
            fmt.Printf("%s%s\n", path, C.GoString(file.name.name))
        }
    }
    err := dir.Walk(cb)
    if err != nil {
        return
    }
    // for i := 0; i < dir.Size(); i += 1 {
    //     fn := dir.GetName(i)
    //     fmt.Printf("    %v\n", fn.Name())
    // }
}

func (fn *Filename) Name() string {
    return C.GoString(fn.name)
}

