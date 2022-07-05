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

    "ndstore/bstore/bsapi"
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
    ALogin      string
    APass       string
    Port        string
    Address     string
    Message     string
    URI         string
    SubCmd      string

    Login       string
    Pass        string
    Role        string
    State       string

    FileId      int64
    BatchId     int64
    BlockId     int64
    BlockType   string
    BlockVer    int64

    FilePath    string
}

func NewUtil() *Util {
    var util Util
    util.Port       = "5101"
    util.Address    = "127.0.0.1"
    util.Message    = "hello"
    util.ALogin     = "admin"
    util.APass      = "admin"
    return &util
}

const getHelloCmd       string = "getHello"
const saveBlockCmd      string = "saveBlock"
const loadBlocksCmd     string = "loadBlock"
const listBlocksCmd     string = "listBlocks"
const deleteBlockCmd    string = "deleteBlock"
const blockExistsCmd    string = "blockExists"
const checkBlockCmd     string = "checkBlock"

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
    flag.StringVar(&util.ALogin, "aLogin", util.ALogin, "access login")
    flag.StringVar(&util.APass, "aPass", util.APass, "access password")

    help := func() {
        fmt.Println("")
        fmt.Printf("Usage: %s [option] command [command option]\n", exeName)
        fmt.Printf("\n")
        fmt.Printf("Command list: hello, saveBlock, loadBlock, listBlocks, deleteBlock, blockExists, checkBlock\n")
        fmt.Printf("              addUser, checkUser, updateUser, listUsers, deleteUser \n")

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
            flagSet := flag.NewFlagSet(getHelloCmd, flag.ExitOnError)
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
        case saveBlockCmd, loadBlocksCmd, deleteBlockCmd, blockExistsCmd, checkBlockCmd:
            flagSet := flag.NewFlagSet(saveBlockCmd, flag.ExitOnError)
            flagSet.Int64Var(&util.FileId, "fileId", util.FileId, "file id")
            flagSet.Int64Var(&util.BatchId, "batchId", util.BatchId, "batch id")
            flagSet.Int64Var(&util.BlockId, "blockId", util.BlockId, "block id")
            flagSet.StringVar(&util.BlockType, "blockType", util.BlockType, "block type")
            flagSet.Int64Var(&util.BlockVer, "blockVer", util.BlockVer, "block version")

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
            flagSet := flag.NewFlagSet(saveBlockCmd, flag.ExitOnError)
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
        case addUserCmd, updateUserCmd:
            flagSet := flag.NewFlagSet(addUserCmd, flag.ExitOnError)
            flagSet.StringVar(&util.Login, "login", util.Login, "login")
            flagSet.StringVar(&util.Pass, "pass", util.Pass, "pass")
            flagSet.StringVar(&util.Role, "role", util.Role, "role")
            flagSet.StringVar(&util.State, "state", util.State, "state")

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
        case checkUserCmd:
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

    auth := dsrpc.CreateAuth([]byte(util.ALogin), []byte(util.APass))

    resp := NewResponse(nil, nil)
    var result interface{}
    switch util.SubCmd {
        case getHelloCmd:
            result, err = util.GetHelloCmd(auth)
        case saveBlockCmd:
            result, err = util.SaveBlockCmd(auth)
        case loadBlocksCmd:
            result, err = util.LoadBlockCmd(auth)
        case listBlocksCmd:
            result, err = util.ListBlocksCmd(auth)
        case deleteBlockCmd:
            result, err = util.DeleteBlockCmd(auth)
        case blockExistsCmd:
            result, err = util.BlockExistsCmd(auth)
        case checkBlockCmd:
            result, err = util.CheckBlockCmd(auth)

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
    respJSON, err := json.MarshalIndent(resp, "", "  ")
    fmt.Printf("%s\n", string(respJSON))
    err = nil
    return err
}

func (util *Util) GetHelloCmd(auth *dsrpc.Auth) (*bsapi.GetHelloResult, error) {
    var err error

    params := bsapi.NewGetHelloParams()
    params.Message = util.Message
    result := bsapi.NewGetHelloResult()

    err = dsrpc.Exec(util.URI, bsapi.GetHelloMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) SaveBlockCmd(auth *dsrpc.Auth) (*bsapi.SaveBlockResult, error) {
    var err error

    params := bsapi.NewSaveBlockParams()
    params.FileId   = util.FileId
    params.BatchId  = util.BatchId
    params.BlockId  = util.BlockId
    params.BlockType = util.BlockType

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
    fileSize := fileInfo.Size()

    params.DataSize = fileSize
    err = dsrpc.Put(util.URI, bsapi.SaveBlockMethod, blockFile, fileSize, params, result, auth)
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
    params.FileId   = util.FileId
    params.BatchId  = util.BatchId
    params.BlockId  = util.BlockId
    params.BlockType = util.BlockType
    params.BlockVer = util.BlockVer

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
    params.FileId   = util.FileId
    params.BatchId  = util.BatchId
    params.BlockId  = util.BlockId
    params.BlockType = util.BlockType
    params.BlockVer = util.BlockVer

    result := bsapi.NewDeleteBlockResult()
    err = dsrpc.Exec(util.URI, bsapi.DeleteBlockMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) BlockExistsCmd(auth *dsrpc.Auth) (*bsapi.BlockExistsResult, error) {
    var err error
    params := bsapi.NewBlockExistsParams()
    params.FileId   = util.FileId
    params.BatchId  = util.BatchId
    params.BlockId  = util.BlockId
    params.BlockType = util.BlockType
    params.BlockVer = util.BlockVer

    result := bsapi.NewBlockExistsResult()
    err = dsrpc.Exec(util.URI, bsapi.BlockExistsMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) CheckBlockCmd(auth *dsrpc.Auth) (*bsapi.CheckBlockResult, error) {
    var err error
    params := bsapi.NewCheckBlockParams()
    params.FileId   = util.FileId
    params.BatchId  = util.BatchId
    params.BlockId  = util.BlockId
    params.BlockType = util.BlockType
    params.BlockVer = util.BlockVer

    result := bsapi.NewCheckBlockResult()
    err = dsrpc.Exec(util.URI, bsapi.CheckBlockMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}


func (util *Util) AddUserCmd(auth *dsrpc.Auth) (*bsapi.AddUserResult, error) {
    var err error
    params := bsapi.NewAddUserParams()
    params.Login = util.Login
    params.Pass = util.Pass
    params.Role = util.Role
    params.State = util.State
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
    params.Login = util.Login
    params.Pass = util.Pass
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
    params.Login = util.Login
    params.Pass = util.Pass
    params.Role = util.Role
    params.State = util.State
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
    params.Login = util.Login
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
