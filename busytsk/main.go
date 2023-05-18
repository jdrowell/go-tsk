package main

import (
    "fmt"
    "os"

    "github.com/jdrowell/go-tsk"
)

// mmls is a clone of the TSK mmls utility.
func mmls(path string) error {
    // FIXME: handle command line options.
    di, err := tsk.OpenImage(path)
    if err != nil { return err }

    vs, err := di.OpenVolumeSystem()
    if err != nil { return err }
    if vs == nil {
        return fmt.Errorf("Unknown volume system")
    }
    defer vs.Close()

    fmt.Println(vs.TypeDescr())
    fmt.Printf("Offset Sector: %d\n", vs.Offset())
    fmt.Printf("Units are in %d-byte sectors\n\n", vs.BlockSize())

    fmt.Println("      Slot      Start        End          Length       Description")
    for i := 0; i < int(vs.PartCount()); i += 1 {
        part := vs.GetPartition(i)

        tn := part.TableNumber()
        sn := part.SlotNumber()
        ts := "-------"
        if tn != -1 {
            ts = fmt.Sprintf("%03d:%03d", tn, sn)
        }

        ps := part.Start()
        pl := part.Len()
        if pl <= 1 {
            ts = "Meta   "
        }

        fmt.Printf("%03d:  %s   %010d   %010d   %010d   %s\n",
            i,
            ts,
            ps,
            ps + pl - 1,
            pl,
            part.Descr() )
    }

    return nil
}

func dispatch() error {
    if len(os.Args) < 2 {
        return fmt.Errorf("Missing image name")
    }
    path := os.Args[1]
    err := mmls(path)

    return err
}

func main() {
    err := dispatch()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

