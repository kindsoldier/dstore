/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsfile

import (
    "fmt"
    "path/filepath"
    "os"
)

type Crate struct {
    dataDir     string
    filePath    string
    file        *os.File
}

const (
    RDONLY  int = iota
    WRONLY
)


func OpenCrate(dataDir, filePath string, direction int) (*Crate, error) {
    var err error
    var crate Crate
    crate.dataDir = dataDir
    crate.filePath = filePath

    switch direction {
        case RDONLY:
            fullPath := filepath.Join(crate.dataDir, crate.filePath)
            crate.file, err = os.OpenFile(fullPath, os.O_CREATE|os.O_RDONLY, 0644)
            if err != nil {
                err = fmt.Errorf("file open error: %s", err)
                return &crate, err
            }

        default:
            fullPath := filepath.Join(crate.dataDir, crate.filePath)
            err = os.MkdirAll(filepath.Dir(fullPath), 0755)
            if err != nil {
                err = fmt.Errorf("file mkdir error: %s", err)
                return &crate, err
            }
            crate.file, err = os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY, 0644)
            if err != nil {
                err = fmt.Errorf("file open error: %s", err)
                return &crate, err
            }
    }
    return &crate, err
}

func (crate *Crate) Write(data []byte) (int, error) {
    var err error
    var written int
    written, err = crate.file.Write(data)
    if err != nil {
        err = fmt.Errorf("file write error: %s", err)
        return written, err
    }
    return written, err
}

func (crate *Crate)  Read(data []byte) (int, error) {
    var err error
    var read int
    read, err = crate.file.Read(data)
    if err != nil {
        err = fmt.Errorf("file read error: %s", err)
        return read, err
    }
    return read, err
}

func (crate *Crate) Close() error {
    var err error
    if crate.file != nil {
        crate.file.Close()
    }
    return err
}

func (crate *Crate) Clean() error {
    var err error
    fullPath := filepath.Join(crate.dataDir, crate.filePath)
    os.Remove(fullPath)
    return err
}
