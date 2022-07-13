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
    dataBase    string
    filePath    string
}

func NewCrate(dataBase, filePath string) (*Crate, error) {
    var err error
    var tank Crate
    tank.dataBase = dataBase
    tank.filePath = filePath
    return &tank, err
}

func (tank *Crate) Write(data []byte) (int, error) {
    var err error
    var written int

    fullPath := filepath.Join(tank.dataBase, tank.filePath)
    err = os.MkdirAll(filepath.Dir(fullPath), 0755)
    if err != nil {
        err = fmt.Errorf("file mkdir error: %s", err)
        return written, err
    }
    file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY, 0655)
    defer file.Close()
    if err != nil {
        err = fmt.Errorf("file open error: %s", err)
        return written, err
    }
    written, err = file.Write(data)
    if err != nil {
        err = fmt.Errorf("file write error: %s", err)
        return written, err
    }
    return written, err
}

func (tank *Crate)  Read(data []byte) (int, error) {
    var err error
    var read int

    fullPath := filepath.Join(tank.dataBase, tank.filePath)
    file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_RDONLY, 0655)
    defer file.Close()
    if err != nil {
        err = fmt.Errorf("file open error: %s", err)
        return read, err
    }
    read, err = file.Read(data)
    if err != nil {
        err = fmt.Errorf("file read error: %s", err)
        return read, err
    }
    return read, err
}

func (tank *Crate) Clean() error {
    var err error
    fullPath := filepath.Join(tank.dataBase, tank.filePath)
    os.Remove(fullPath)
    return err
}
