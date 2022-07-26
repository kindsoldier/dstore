/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package main

import (
    "encoding/json"
    "fmt"
    "io/fs"
    "flag"
    "os"
    "path/filepath"
    "errors"

    "dstore/bstore/bsapi"
    "dstore/dscomm/dsrpc"
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
    aLogin      string
    aPass       string

    Port        string
    Address     string
    Message     string
    URI         string
    SubCmd      string

    Login       string
    Pass        string

    bPort       string
    bAddress    string

    FileId      int64
    BatchId     int64
    BlockId     int64
    BlockType   int64

    FilePath   string
}

func NewUtil() *Util {
    var util Util
    util.Port       = "5101"
    util.Address    = "127.0.0.1"
    util.Message    = "hello"
    util.aLogin     = "admin"
    util.aPass      = "admin"
    return &util
}

const getStatusCmd      string = "getStatus"

const saveBlockCmd      string = "saveBlock"
const loadBlockCmd      string = "loadBlock"
const listBlocksCmd     string = "listBlocks"
const deleteBlockCmd    string = "deleteBlock"

const addUserCmd        string = "addUser"
const checkUserCmd      string = "checkUser"
const updateUserCmd     string = "updateUser"
const deleteUserCmd     string = "deleteUser"
const listUsersCmd      string = "listUsers"


const helpCmd           string = "help"


func (util *Util) GetOpt() error {
    var err error

    exeName := filepath.Base(os.Args[0])

    flag.StringVar(&util.Port, "port", util.Port, "service port")
    flag.StringVar(&util.Address, "address", util.Address, "service address")
    flag.StringVar(&util.aLogin, "aLogin", util.aLogin, "access login")
    flag.StringVar(&util.aPass, "aPass", util.aPass, "access password")

    help := func() {
        fmt.Println("")
        fmt.Printf("Usage: %s [option] command [command option]\n", exeName)
        fmt.Printf("\n")
        fmt.Printf("Command list: help, getStatus, \n")
        fmt.Printf("    saveBlock, loadBlock, listBlocks, deleteBlock \n")
        fmt.Printf("    addUser, checkUser, updateUser, listUsers, deleteUser \n")

        fmt.Printf("\n")
        fmt.Printf("Global options:\n")
        flag.PrintDefaults()
        fmt.Printf("\n")
    }
    flag.Usage = help
    flag.Parse()

    args := flag.Args()

    //if len(args) == 0 {
    //    args = append(args, getStatusCmd)
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
        case getStatusCmd:
            flagSet := flag.NewFlagSet(getStatusCmd, flag.ExitOnError)
            flagSet.Usage = func() {
                fmt.Printf("\n")
                fmt.Printf("Usage: %s [global options] %s [command options]\n", exeName, subCmd)
                fmt.Printf("\n")
                fmt.Printf("The command options: none\n")
                flagSet.PrintDefaults()
                fmt.Printf("\n")
            }
            flagSet.Parse(subArgs)
            util.SubCmd = subCmd

        case saveBlockCmd, loadBlockCmd, deleteBlockCmd:
            flagSet := flag.NewFlagSet(saveBlockCmd, flag.ExitOnError)
            flagSet.Int64Var(&util.FileId, "fileId", util.FileId, "file id")
            flagSet.Int64Var(&util.BatchId, "batchId", util.BatchId, "batch id")
            flagSet.Int64Var(&util.BlockType, "blockType", util.BlockType, "block type")
            flagSet.Int64Var(&util.BlockId, "blockId", util.BlockId, "block id")
            flagSet.StringVar(&util.FilePath, "file", util.FilePath, "block file name")
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
        case listBlocksCmd:
            flagSet := flag.NewFlagSet(listBlocksCmd, flag.ExitOnError)
            flagSet.Int64Var(&util.FileId, "fileId", util.FileId, "file id")

            flagSet.Usage = func() {
                fmt.Printf("\n")
                fmt.Printf("Usage: %s [global options] %s [command options]\n", exeName, subCmd)
                fmt.Printf("\n")
                fmt.Printf("The command options: none\n")
                flagSet.PrintDefaults()
                fmt.Printf("\n")
            }
            flagSet.Parse(subArgs)
            util.SubCmd = subCmd

        case addUserCmd, checkUserCmd, updateUserCmd:
            flagSet := flag.NewFlagSet(addUserCmd, flag.ExitOnError)
            flagSet.StringVar(&util.Login, "login", util.Login, "login")
            flagSet.StringVar(&util.Pass, "pass", util.Pass, "pass")
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
        case deleteUserCmd:
            flagSet := flag.NewFlagSet(deleteUserCmd, flag.ExitOnError)
            flagSet.StringVar(&util.Login, "login", util.Login, "login")
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
        case listUsersCmd:
            flagSet := flag.NewFlagSet(deleteUserCmd, flag.ExitOnError)
            flagSet.Usage = func() {
                fmt.Printf("\n")
                fmt.Printf("Usage: %s [global options] %s [command options]\n", exeName, subCmd)
                fmt.Printf("\n")
                fmt.Printf("The command options: none\n")
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
    Error       bool       `json:"error"`
    ErrorMsg    string     `json:"errorMsg,omitempty"`
    Result      any        `json:"result,omitempty"`
}

func NewResponse(result any, err error) *Response {
    var errMsg string
    var errBool bool
    if err != nil {
        errMsg = err.Error()
        errBool = true
    }
    return &Response{
        Result:     result,
        Error:      errBool,
        ErrorMsg:   errMsg,
    }
}

func (util *Util) Exec() error {
    var err error
    err = util.GetOpt()
    if err != nil {
        return err
    }
    util.URI = fmt.Sprintf("%s:%s", util.Address, util.Port)
    auth := dsrpc.CreateAuth([]byte(util.aLogin), []byte(util.aPass))

    resp := NewResponse(nil, nil)
    var result interface{}

    switch util.SubCmd {
        case getStatusCmd:
            result, err = util.GetStatusCmd(auth)

        case saveBlockCmd:
            result, err = util.SaveBlockCmd(auth)
        case loadBlockCmd:
            result, err = util.LoadBlockCmd(auth)
        case listBlocksCmd:
            result, err = util.ListBlocksCmd(auth)
        case deleteBlockCmd:
            result, err = util.DeleteBlockCmd(auth)

        case addUserCmd:
            result, err = util.AddUserCmd(auth)
        case checkUserCmd:
            result, err = util.CheckUserCmd(auth)
        case updateUserCmd:
            result, err = util.UpdateUserCmd(auth)
        case deleteUserCmd:
            result, err = util.DeleteUserCmd(auth)
        case listUsersCmd:
            result, err = util.ListUsersCmd(auth)

        default:
            err = errors.New("unknown cli command")
    }
    resp = NewResponse(result, err)
    respJSON, _ := json.MarshalIndent(resp, "", "  ")
    fmt.Printf("%s\n", string(respJSON))
    err = nil
    return err
}

func (util *Util) GetStatusCmd(auth *dsrpc.Auth) (*bsapi.GetStatusResult, error) {
    var err error
    params := bsapi.NewGetStatusParams()
    result := bsapi.NewGetStatusResult()
    err = dsrpc.Exec(util.URI, bsapi.GetStatusMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) SaveBlockCmd(auth *dsrpc.Auth) (*bsapi.SaveBlockResult, error) {
    var err error

    result := bsapi.NewSaveBlockResult()

    blockFile, err := os.OpenFile(util.FilePath, os.O_RDONLY, 0)
    defer blockFile.Close()
    if err != nil {
        return result, err
    }
    fileInfo, err := blockFile.Stat()
    if err != nil {
        return result, err
    }
    dataSize := fileInfo.Size()

    params := bsapi.NewSaveBlockParams()
    params.FileId       = util.FileId
    params.BatchId      = util.BatchId
    params.BlockType    = util.BlockType
    params.BlockId      = util.BlockId
    params.BlockSize = dataSize

    err = dsrpc.Put(util.URI, bsapi.SaveBlockMethod, blockFile, dataSize, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

const dirPerm   fs.FileMode = 0755
const filePerm  fs.FileMode = 0644

func (util *Util) LoadBlockCmd(auth *dsrpc.Auth) (*bsapi.LoadBlockResult, error) {
    var err error
    params := bsapi.NewLoadBlockParams()
    params.FileId       = util.FileId
    params.BatchId      = util.BatchId
    params.BlockType    = util.BlockType
    params.BlockId      = util.BlockId

    result := bsapi.NewLoadBlockResult()
    blockFile, err := os.OpenFile(util.FilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filePerm)
    defer blockFile.Close()
    if err != nil {
        return result, err
    }
    err = dsrpc.Get(util.URI, bsapi.LoadBlockMethod, blockFile, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) ListBlocksCmd(auth *dsrpc.Auth) (*bsapi.ListBlocksResult, error) {
    var err error
    params := bsapi.NewListBlocksParams()
    params.FileId       = util.FileId
    result := bsapi.NewListBlocksResult()
    err = dsrpc.Exec(util.URI, bsapi.ListBlocksMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) DeleteBlockCmd(auth *dsrpc.Auth) (*bsapi.DeleteBlockResult, error) {
    var err error
    params := bsapi.NewDeleteBlockParams()
    params.FileId       = util.FileId
    params.BatchId      = util.BatchId
    params.BlockType    = util.BlockType
    params.BlockId      = util.BlockId

    result := bsapi.NewDeleteBlockResult()
    err = dsrpc.Exec(util.URI, bsapi.DeleteBlockMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) AddUserCmd(auth *dsrpc.Auth) (*bsapi.AddUserResult, error) {
    var err error
    params := bsapi.NewAddUserParams()
    params.Login    = util.Login
    params.Pass     = util.Pass
    result := bsapi.NewAddUserResult()
    err = dsrpc.Exec(util.URI, bsapi.AddUserMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) CheckUserCmd(auth *dsrpc.Auth) (*bsapi.CheckUserResult, error) {
    var err error
    params := bsapi.NewCheckUserParams()
    params.Login    = util.Login
    params.Pass     = util.Pass
    result := bsapi.NewCheckUserResult()
    err = dsrpc.Exec(util.URI, bsapi.CheckUserMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) UpdateUserCmd(auth *dsrpc.Auth) (*bsapi.UpdateUserResult, error) {
    var err error
    params := bsapi.NewUpdateUserParams()
    params.Login    = util.Login
    params.Pass     = util.Pass
    result := bsapi.NewUpdateUserResult()
    err = dsrpc.Exec(util.URI, bsapi.UpdateUserMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) DeleteUserCmd(auth *dsrpc.Auth) (*bsapi.DeleteUserResult, error) {
    var err error
    params := bsapi.NewDeleteUserParams()
    params.Login    = util.Login
    result := bsapi.NewDeleteUserResult()
    err = dsrpc.Exec(util.URI, bsapi.DeleteUserMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) ListUsersCmd(auth *dsrpc.Auth) (*bsapi.ListUsersResult, error) {
    var err error
    params := bsapi.NewListUsersParams()
    result := bsapi.NewListUsersResult()
    err = dsrpc.Exec(util.URI, bsapi.ListUsersMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}
