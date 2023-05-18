package tsk

/*
#cgo LDFLAGS: -ltsk
#include <tsk/libtsk.h>
#include "wrapper.h"
*/
import "C"

import (
    "errors"
    "fmt"
    "os"
    "unsafe"
)

type Filesystem C.TSK_FS_INFO
type Directory  C.TSK_FS_DIR
type Filename   C.TSK_FS_NAME
type File       C.TSK_FS_FILE

type WalkCallback func (file *File, path string)
// keep track of ongoing walks
var walkCBs = make(map[int]WalkCallback)

var WalkFlag = map[string]int{
    "None": int(C.TSK_FS_DIR_WALK_FLAG_NONE),
    "Alloc": int(C.TSK_FS_DIR_WALK_FLAG_ALLOC),
    "Unalloc": int(C.TSK_FS_DIR_WALK_FLAG_UNALLOC),
    "Recurse": int(C.TSK_FS_DIR_WALK_FLAG_RECURSE),
    "NoOrphan": int(C.TSK_FS_DIR_WALK_FLAG_NOORPHAN),
}

var FileType = map[int]string{
    int(C.TSK_FS_NAME_TYPE_UNDEF): "UNDEF",
    int(C.TSK_FS_NAME_TYPE_FIFO): "FIFO",
    int(C.TSK_FS_NAME_TYPE_CHR): "CHR",
    int(C.TSK_FS_NAME_TYPE_DIR): "DIR",
    int(C.TSK_FS_NAME_TYPE_BLK): "BLK",
    int(C.TSK_FS_NAME_TYPE_REG): "REG",
    int(C.TSK_FS_NAME_TYPE_LNK): "LNK",
    int(C.TSK_FS_NAME_TYPE_SOCK): "SOCK",
    int(C.TSK_FS_NAME_TYPE_SHAD): "SHAD",
    int(C.TSK_FS_NAME_TYPE_WHT): "WHT",
    int(C.TSK_FS_NAME_TYPE_VIRT): "VIRT",
    int(C.TSK_FS_NAME_TYPE_VIRT_DIR): "VIRT_DIR",
}

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

func (fs *Filesystem) PathToInum(path string) (*Filename, error) {
    var status int
    cStatus := C.ulong(status)
    filename := Filename{}
    res := C.tsk_fs_path2inum((*C.TSK_FS_INFO)(fs), C.CString(path), &cStatus, (*C.TSK_FS_NAME)(&filename))
    if res == -1 {
        return nil, errors.New("System error in PathToInum")
    }
    if res == 1 {
        return nil, nil
    }
    return &filename, nil
}

func (fs *Filesystem) OpenDirectory(aMeta int) *Directory {
    cDir := C.tsk_fs_dir_open_meta((*C.TSK_FS_INFO)(fs), C.ulong(aMeta))

    return (*Directory)(cDir)
}

func (fs *Filesystem) OpenFile(path string) (*File, error) {
    cFile := C.tsk_fs_file_open((*C.TSK_FS_INFO)(fs), nil, C.CString(path))
    if cFile == nil {
        return nil, Error()
    }
    return (*File)(cFile), nil
}

func (fs *Filesystem) Show() {
    fmt.Printf("    filesystem type: %s (%d)\n", fs.TypeName(), fs.Type())
    // get the root directory
    // dir := fs.OpenDirectory(fs.RootInum())
    // if fs == nil {
    //     fmt.Printf("    couldn't open root directory\n")
    //     return
    // }
    // defer dir.Close()
    //
    // dir.Show()
    path := "/Windows/System32/config/SYSTEM"
    file, err := fs.OpenFile(path)
    if err != nil {
        fmt.Printf("ERROR: %v\n", err)
        return
    }
    fmt.Printf("file: %v\n", file)
}

func (dir *Directory) Walk(flags int, cb WalkCallback) error {
    k := int(dir.addr)
    walkCBs[k] = cb
    defer delete(walkCBs, k)
    data := unsafe.Pointer(&k)
    ret := C.tsk_fs_dir_walk(dir.fs_info, dir.addr, (C.TSK_FS_DIR_WALK_FLAG_ENUM)(flags), (*[0]byte)(C.dirwalk_callback), data)

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

func (dir *Directory) Show() error {
    cb := func (file *File, path string) {
        filename := file.Name()
        if filename.TypeStr() == "REG" {
            fmt.Printf("%s%s\n", path, filename.Name())
        }
    }
    err := dir.Walk(WalkFlag["Alloc"] | WalkFlag["Recurse"], cb)
    if err != nil {
        return err
    }
    // for i := 0; i < dir.Size(); i += 1 {
    //     fn := dir.GetName(i)
    //     fmt.Printf("    %v\n", fn.Name())
    // }
    return nil
}

func (file *File) Close() {
    C.tsk_fs_file_close((*C.TSK_FS_FILE)(file))
}

func (file *File) Name() *Filename {
    return (*Filename)(file.name)
}

func (file *File) Copy(dest string) error {
    buf := make([]byte, 65536)
    destf, err := os.Create(dest)
    if err != nil {
        return err
    }
    defer destf.Close()

    offset := 0
    for {
        read := int(C.tsk_fs_file_read((*C.TSK_FS_FILE)(file), C.long(offset), (*C.char)(unsafe.Pointer(&buf[0])), C.ulong(len(buf)), 0))
        if read == -1 {
            return Error()
        }
        if read == 0 {
            // we're done
            break
        }
        _, err := destf.Write(buf[0:read])
        if err != nil {
            return err
        }
        offset += read
    }

    return nil
}

func (fn *Filename) Name() string {
    return C.GoString(fn.name)
}

func (fn *Filename) Type() int {
    return int(fn._type)
}

func (fn *Filename) TypeStr() string {
    return FileType[int(fn.Type())]
}
