/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */


package main

import (
    "fmt"
    "flag"
    "os"
    "io/fs"
    "path/filepath"
    "errors"

    "dcstore/dcnode/ndapi"
    "dcstore/dcrpc"
)

func main() {
    var err error
    util := NewUtil()
    err = util.Exec()
    if err != nil {
        fmt.Printf("Exec error: %s\n", err)
    }
}

type Util struct {
    Port        string
    Address     string
    Message     string
    URI         string
    SubCmd      string

    ClusterId   int64
    FileId      int64
    BatchId     int64
    BlockId     int64

    FilePath    string
}

func NewUtil() *Util {
    var util Util
    util.Port = "5001"
    util.Address = "127.0.0.1"
    util.Message = "hello world!"
    return &util
}

const helpCmd   string = "help"
const helloCmd  string = "hello"
const saveCmd   string = "save"
const loadCmd   string = "load"

func (util *Util) GetOpt() error {
    var err error

    exeName := filepath.Base(os.Args[0])

    flag.StringVar(&util.Port, "p", util.Port, "port")
    flag.StringVar(&util.Address, "a", util.Address, "address")

    help := func() {
        fmt.Println("")
        fmt.Printf("Usage: %s [option] command [command option]\n", exeName)
        fmt.Printf("\n")
        fmt.Printf("Command list: hello, save, load, list \n")
        fmt.Printf("\n")
        fmt.Printf("Global options:\n")
        flag.PrintDefaults()
        fmt.Printf("\n")
    }
    flag.Usage = help
    flag.Parse()

    args := flag.Args()

    //if len(args) == 0 {
    //    args = append(args, helloCmd)
    //}

    var subCmd string
    var subArgs []string
    if len(args) > 0 {
        subCmd = args[0]
        subArgs = args[1:]
    }
    switch subCmd {
        case helpCmd:
            help()
            return errors.New("unknown command")
        case helloCmd:
            flagSet := flag.NewFlagSet(helloCmd, flag.ContinueOnError)
            flagSet.StringVar(&util.Message, "m", util.Message, "message")
            flagSet.Usage = func() {
                fmt.Printf("\n")
                fmt.Printf("Usage: %s [global options] %s [command options]\n", exeName, subCmd)
                fmt.Printf("\n")
                fmt.Printf("The command options:\n")
                flagSet.PrintDefaults()
                fmt.Printf("\n")
            }
            flagSet.Parse(subArgs)
            util.SubCmd = subCmd
        case saveCmd, loadCmd:
            flagSet := flag.NewFlagSet(saveCmd, flag.ContinueOnError)
            flagSet.Int64Var(&util.ClusterId, "c", util.ClusterId, "cluster id")
            flagSet.Int64Var(&util.FileId, "f", util.FileId, "file id")
            flagSet.Int64Var(&util.BatchId, "ba", util.BatchId, "batch id")
            flagSet.Int64Var(&util.BlockId, "bl", util.BlockId, "block id")
            flagSet.StringVar(&util.FilePath, "n", util.FilePath, "block file name")

            flagSet.Usage = func() {
                fmt.Printf("\n")
                fmt.Printf("Usage: %s [global options] %s [command options]\n", exeName, subCmd)
                fmt.Printf("\n")
                fmt.Printf("The command options:\n")
                flagSet.PrintDefaults()
                fmt.Printf("\n")
            }
            flagSet.Parse(subArgs)
            util.SubCmd = subCmd
        default:
            help()
            return errors.New("unknown command")
    }
    return err
}


func (util *Util) Exec() error {
    var err error
    err = util.GetOpt()
    if err != nil {
        return err
    }
    util.URI = fmt.Sprintf("%s:%s", util.Address, util.Port)

    switch util.SubCmd {
        case helloCmd:
            err = util.HelloCmd()
        case saveCmd:
            err = util.SaveCmd()
        case loadCmd:
            err = util.LoadCmd()
        default:
            fmt.Printf("%s\n", util.SubCmd)
    }
    return err
}

func (util *Util) HelloCmd() error {
    var err error

    params := ndapi.NewHelloParams()
    params.Message = util.Message
    result := ndapi.NewHelloResult()

    err = dcrpc.Exec(util.URI, ndapi.HelloMethod, params, result, nil)
    if err != nil {
        return err
    }
    fmt.Printf("result: %s\n", result.Message)
    return err
}

func (util *Util) SaveCmd() error {
    var err error

    params := ndapi.NewSaveParams()
    params.ClusterId    = util.ClusterId
    params.FileId       = util.FileId
    params.BatchId      = util.BatchId
    params.BlockId      = util.BlockId

    result := ndapi.NewSaveResult()

    blockFile, err := os.OpenFile(util.FilePath, os.O_RDONLY, 0)
    defer blockFile.Close()
    if err != nil {
        return err
    }
    fileInfo, err := blockFile.Stat()
    if err != nil {
        return err
    }
    fileSize := fileInfo.Size()

    err = dcrpc.Put(util.URI, ndapi.SaveMethod, blockFile, fileSize, params, result, nil)
    if err != nil {
        return err
    }
    fmt.Printf("result: ok\n")
    return err
}

const dirPerm   fs.FileMode = 0755
const filePerm  fs.FileMode = 0644

func (util *Util) LoadCmd() error {
    var err error

    params := ndapi.NewLoadParams()
    params.ClusterId    = util.ClusterId
    params.FileId       = util.FileId
    params.BatchId      = util.BatchId
    params.BlockId      = util.BlockId

    result := ndapi.NewLoadResult()

    blockFile, err := os.OpenFile(util.FilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filePerm)
    defer blockFile.Close()
    if err != nil {
        return err
    }
    err = dcrpc.Get(util.URI, ndapi.LoadMethod, blockFile, params, result, nil)
    if err != nil {
        return err
    }
    fmt.Printf("result: ok\n")
    return err
}
