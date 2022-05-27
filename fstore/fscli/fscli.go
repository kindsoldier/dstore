/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */


package main

import (
    "encoding/json"
    "fmt"
    "flag"
    "os"
    "io/fs"
    "path/filepath"
    "errors"

    "ndstore/fstore/fsapi"
    "ndstore/dsrpc"
)

type any = interface{}

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

    LocalFilePath   string
    RemoteFilePath  string
}

func NewUtil() *Util {
    var util Util
    util.Port       = "5002"
    util.Address    = "127.0.0.1"
    util.Message    = "hello"
    return &util
}

const getHelloCmd       string = "hello"
const saveFileCmd       string = "save"
const loadFileCmd       string = "load"
const listFilesCmd      string = "list"
const deleteFileCmd     string = "delete"
const helpCmd           string = "help"


func (util *Util) GetOpt() error {
    var err error

    exeName := filepath.Base(os.Args[0])

    flag.StringVar(&util.Port, "p", util.Port, "port")
    flag.StringVar(&util.Address, "a", util.Address, "address")

    help := func() {
        fmt.Println("")
        fmt.Printf("Usage: %s [option] command [command option]\n", exeName)
        fmt.Printf("\n")
        fmt.Printf("Command list: hello, save, load, list, delete \n")
        fmt.Printf("\n")
        fmt.Printf("Global options:\n")
        flag.PrintDefaults()
        fmt.Printf("\n")
    }
    flag.Usage = help
    flag.Parse()

    args := flag.Args()

    //if len(args) == 0 {
    //    args = append(args, getHelloCmd)
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
        case getHelloCmd:
            flagSet := flag.NewFlagSet(getHelloCmd, flag.ContinueOnError)
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
        case saveFileCmd, loadFileCmd:
            flagSet := flag.NewFlagSet(saveFileCmd, flag.ContinueOnError)
            flagSet.StringVar(&util.LocalFilePath, "f", util.LocalFilePath, "local file name")
            flagSet.StringVar(&util.RemoteFilePath, "f", util.RemoteFilePath, "remote file path")

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
        case listFilesCmd:
            flagSet := flag.NewFlagSet(listFilesCmd, flag.ContinueOnError)
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
        case deleteFileCmd:
            flagSet := flag.NewFlagSet(saveFileCmd, flag.ContinueOnError)
            flagSet.StringVar(&util.RemoteFilePath, "f", util.RemoteFilePath, "remote file path")

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

type Response struct {
    Error   error      `json:"error"`
    Result  any        `json:"result,omitempty"`
}

func NewResponse(result any, err error) *Response {
    return &Response{ Result: result, Error: err }
}

func (util *Util) Exec() error {
    var err error
    err = util.GetOpt()
    if err != nil {
        return err
    }
    util.URI = fmt.Sprintf("%s:%s", util.Address, util.Port)

    resp := NewResponse(nil, nil)

    switch util.SubCmd {
        case getHelloCmd:
            result, err := util.GetHelloCmd()
            resp = NewResponse(result, err)
        case saveFileCmd:
            result, err := util.SaveFileCmd()
            resp = NewResponse(result, err)
        case loadFileCmd:
            result, err := util.LoadFileCmd()
            resp = NewResponse(result, err)
        case listFilesCmd:
            result, err := util.ListFilesCmd()
            resp = NewResponse(result, err)
        case deleteFileCmd:
            result, err := util.DeleteFileCmd()
            resp = NewResponse(result, err)
        default:
    }
    respJSON, _ := json.Marshal(resp)
    fmt.Printf("%s\n", string(respJSON))
    err = nil
    return err
}

func (util *Util) GetHelloCmd() (*fsapi.GetHelloResult, error) {
    var err error

    params := fsapi.NewGetHelloParams()
    params.Message = util.Message
    result := fsapi.NewGetHelloResult()

    err = dsrpc.Exec(util.URI, fsapi.GetHelloMethod, params, result, nil)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) SaveFileCmd() (*fsapi.SaveFileResult, error) {
    var err error

    params := fsapi.NewSaveFileParams()
    params.FilePath  = util.RemoteFilePath

    result := fsapi.NewSaveFileResult()

    localFile, err := os.OpenFile(util.LocalFilePath, os.O_RDONLY, 0)
    defer localFile.Close()
    if err != nil {
        return result, err
    }
    fileInfo, err := localFile.Stat()
    if err != nil {
        return result, err
    }
    fileSize := fileInfo.Size()

    err = dsrpc.Put(util.URI, fsapi.SaveFileMethod, localFile, fileSize, params, result, nil)
    if err != nil {
        return result, err
    }
    return result, err
}

const dirPerm   fs.FileMode = 0755
const filePerm  fs.FileMode = 0644

func (util *Util) LoadFileCmd() (*fsapi.LoadFileResult, error) {
    var err error

    params := fsapi.NewLoadFileParams()
    params.FilePath   = util.RemoteFilePath

    result := fsapi.NewLoadFileResult()

    localFile, err := os.OpenFile(util.LocalFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filePerm)
    defer localFile.Close()
    if err != nil {
        return result, err
    }
    err = dsrpc.Get(util.URI, fsapi.LoadFileMethod, localFile, params, result, nil)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) ListFilesCmd() (*fsapi.ListFilesResult, error) {
    var err error
    params := fsapi.NewListFilesParams()
    result := fsapi.NewListFilesResult()
    err = dsrpc.Exec(util.URI, fsapi.ListFilesMethod, params, result, nil)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) DeleteFileCmd() (*fsapi.DeleteFileResult, error) {
    var err error
    params := fsapi.NewDeleteFileParams()
    params.FilePath   = util.RemoteFilePath

    result := fsapi.NewDeleteFileResult()
    err = dsrpc.Exec(util.URI, fsapi.DeleteFileMethod, params, result, nil)
    if err != nil {
        return result, err
    }
    return result, err
}
