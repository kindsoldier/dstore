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

    "dstore/fstore/fsapi"
    "dstore/dsrpc"
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

    LocalFilePath   string
    RemoteFilePath  string
}

func NewUtil() *Util {
    var util Util
    util.Port       = "5200"
    util.Address    = "127.0.0.1"
    util.Message    = "hello"
    util.aLogin     = "admin"
    util.aPass      = "admin"
    return &util
}

const getStatusCmd       string = "getStatus"
const saveFileCmd       string = "saveFile"
const loadFileCmd       string = "loadFile"
const listFilesCmd      string = "listFiles"
const deleteFileCmd     string = "deleteFile"

const addUserCmd        string = "addUser"
const checkUserCmd      string = "checkUser"
const updateUserCmd     string = "updateUser"
const deleteUserCmd     string = "deleteUser"
const listUsersCmd      string = "listUsers"

const addBStoreCmd        string = "addBStore"
const updateBStoreCmd     string = "updateBStore"
const deleteBStoreCmd     string = "deleteBStore"
const listBStoresCmd      string = "listBStores"


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
        fmt.Printf("    saveFile, loadFile, listFiles, deleteFile \n")
        fmt.Printf("    addUser, checkUser, updateUser, listUsers, deleteUser \n")
//        fmt.Printf("    addBStore, checkBStore, updateBStore, listBStores, deleteBStore \n")

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
        case saveFileCmd, loadFileCmd:
            flagSet := flag.NewFlagSet(saveFileCmd, flag.ExitOnError)
            flagSet.StringVar(&util.LocalFilePath, "local", util.LocalFilePath, "local file name")
            flagSet.StringVar(&util.RemoteFilePath, "remote", util.RemoteFilePath, "remote file path")
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
            flagSet := flag.NewFlagSet(listFilesCmd, flag.ExitOnError)
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
        case deleteFileCmd:
            flagSet := flag.NewFlagSet(saveFileCmd, flag.ExitOnError)
            flagSet.StringVar(&util.RemoteFilePath, "path", util.RemoteFilePath, "remote file path")

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

//        case addBStoreCmd, updateBStoreCmd:
//            flagSet := flag.NewFlagSet(addBStoreCmd, flag.ExitOnError)
//            flagSet.StringVar(&util.bAddress, "address", util.bAddress, "address")
//            flagSet.StringVar(&util.bPort, "port", util.bPort, "port")
//            flagSet.StringVar(&util.Login, "login", util.Login, "login")
//            flagSet.StringVar(&util.Pass, "pass", util.Pass, "pass")
//            flagSet.Usage = func() {
//                fmt.Printf("\n")
//                fmt.Printf("Usage: %s [global options] %s [command options]\n", exeName, subCmd)
//                fmt.Printf("\n")
//                fmt.Printf("The command options:\n")
//                flagSet.PrintDefaults()
//                fmt.Printf("\n")
//            }
//            flagSet.Parse(subArgs)
//            util.SubCmd = subCmd
//        case deleteBStoreCmd:
//            flagSet := flag.NewFlagSet(deleteBStoreCmd, flag.ExitOnError)
//            flagSet.StringVar(&util.bAddress, "address", util.bAddress, "address")
//            flagSet.StringVar(&util.bPort, "port", util.bPort, "port")
//            flagSet.Usage = func() {
//                fmt.Printf("\n")
//                fmt.Printf("Usage: %s [global options] %s [command options]\n", exeName, subCmd)
//                fmt.Printf("\n")
//                fmt.Printf("The command options:\n")
//                flagSet.PrintDefaults()
//                fmt.Printf("\n")
//            }
//            flagSet.Parse(subArgs)
//            util.SubCmd = subCmd
//        case listBStoresCmd:
//            flagSet := flag.NewFlagSet(deleteBStoreCmd, flag.ExitOnError)
//            flagSet.Usage = func() {
//                fmt.Printf("\n")
//                fmt.Printf("Usage: %s [global options] %s [command options]\n", exeName, subCmd)
//                fmt.Printf("\n")
//                fmt.Printf("The command options: none\n")
//                flagSet.PrintDefaults()
//                fmt.Printf("\n")
//            }
//            flagSet.Parse(subArgs)
//            util.SubCmd = subCmd
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

        case saveFileCmd:
            result, err = util.SaveFileCmd(auth)
        case loadFileCmd:
            result, err = util.LoadFileCmd(auth)
        case listFilesCmd:
            result, err = util.ListFilesCmd(auth)
        case deleteFileCmd:
            result, err = util.DeleteFileCmd(auth)

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

//        case addBStoreCmd:
//            result, err = util.AddBStoreCmd(auth)
//        case updateBStoreCmd:
//            result, err = util.UpdateBStoreCmd(auth)
//        case deleteBStoreCmd:
//            result, err = util.DeleteBStoreCmd(auth)
//        case listBStoresCmd:
//            result, err = util.ListBStoresCmd(auth)
        default:
            err = errors.New("unknown cli command")
    }
    resp = NewResponse(result, err)
    respJSON, _ := json.MarshalIndent(resp, "", "  ")
    fmt.Printf("%s\n", string(respJSON))
    err = nil
    return err
}

func (util *Util) GetStatusCmd(auth *dsrpc.Auth) (*fsapi.GetStatusResult, error) {
    var err error
    params := fsapi.NewGetStatusParams()
    result := fsapi.NewGetStatusResult()
    err = dsrpc.Exec(util.URI, fsapi.GetStatusMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) SaveFileCmd(auth *dsrpc.Auth) (*fsapi.SaveFileResult, error) {
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

    err = dsrpc.Put(util.URI, fsapi.SaveFileMethod, localFile, fileSize, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

const dirPerm   fs.FileMode = 0755
const filePerm  fs.FileMode = 0644

func (util *Util) LoadFileCmd(auth *dsrpc.Auth) (*fsapi.LoadFileResult, error) {
    var err error
    params := fsapi.NewLoadFileParams()
    params.FilePath   = util.RemoteFilePath
    result := fsapi.NewLoadFileResult()
    localFile, err := os.OpenFile(util.LocalFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filePerm)
    defer localFile.Close()
    if err != nil {
        return result, err
    }
    err = dsrpc.Get(util.URI, fsapi.LoadFileMethod, localFile, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) ListFilesCmd(auth *dsrpc.Auth) (*fsapi.ListFilesResult, error) {
    var err error
    params := fsapi.NewListFilesParams()
    result := fsapi.NewListFilesResult()
    err = dsrpc.Exec(util.URI, fsapi.ListFilesMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) DeleteFileCmd(auth *dsrpc.Auth) (*fsapi.DeleteFileResult, error) {
    var err error
    params := fsapi.NewDeleteFileParams()
    params.FilePath   = util.RemoteFilePath
    result := fsapi.NewDeleteFileResult()
    err = dsrpc.Exec(util.URI, fsapi.DeleteFileMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) AddUserCmd(auth *dsrpc.Auth) (*fsapi.AddUserResult, error) {
    var err error
    params := fsapi.NewAddUserParams()
    params.Login = util.Login
    params.Pass = util.Pass
    result := fsapi.NewAddUserResult()
    err = dsrpc.Exec(util.URI, fsapi.AddUserMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) CheckUserCmd(auth *dsrpc.Auth) (*fsapi.CheckUserResult, error) {
    var err error
    params := fsapi.NewCheckUserParams()
    params.Login = util.Login
    params.Pass = util.Pass
    result := fsapi.NewCheckUserResult()
    err = dsrpc.Exec(util.URI, fsapi.CheckUserMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) UpdateUserCmd(auth *dsrpc.Auth) (*fsapi.UpdateUserResult, error) {
    var err error
    params := fsapi.NewUpdateUserParams()
    params.Login = util.Login
    params.Pass = util.Pass
    result := fsapi.NewUpdateUserResult()
    err = dsrpc.Exec(util.URI, fsapi.UpdateUserMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) DeleteUserCmd(auth *dsrpc.Auth) (*fsapi.DeleteUserResult, error) {
    var err error
    params := fsapi.NewDeleteUserParams()
    params.Login = util.Login
    result := fsapi.NewDeleteUserResult()
    err = dsrpc.Exec(util.URI, fsapi.DeleteUserMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) ListUsersCmd(auth *dsrpc.Auth) (*fsapi.ListUsersResult, error) {
    var err error
    params := fsapi.NewListUsersParams()
    result := fsapi.NewListUsersResult()
    err = dsrpc.Exec(util.URI, fsapi.ListUsersMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}


//func (util *Util) AddBStoreCmd(auth *dsrpc.Auth) (*fsapi.AddBStoreResult, error) {
//    var err error
//    params := fsapi.NewAddBStoreParams()
//    params.Address = util.bAddress
//    params.Port    = util.bPort
//    params.Login    = util.Login
//    params.Pass     = util.Pass
//    result := fsapi.NewAddBStoreResult()
//    err = dsrpc.Exec(util.URI, fsapi.AddBStoreMethod, params, result, auth)
//    if err != nil {
//        return result, err
//    }
//    return result, err
//}

//func (util *Util) UpdateBStoreCmd(auth *dsrpc.Auth) (*fsapi.UpdateBStoreResult, error) {
//    var err error
//    params := fsapi.NewUpdateBStoreParams()
//    params.Address = util.bAddress
//    params.Port    = util.bPort
//    params.Login    = util.Login
//    params.Pass     = util.Pass
//    result := fsapi.NewUpdateBStoreResult()
//    err = dsrpc.Exec(util.URI, fsapi.UpdateBStoreMethod, params, result, auth)
//    if err != nil {
//        return result, err
//    }
//    return result, err
//}

//func (util *Util) DeleteBStoreCmd(auth *dsrpc.Auth) (*fsapi.DeleteBStoreResult, error) {
//    var err error
//    params := fsapi.NewDeleteBStoreParams()
//    params.Address = util.bAddress
//    params.Port    = util.bPort
//    result := fsapi.NewDeleteBStoreResult()
//    err = dsrpc.Exec(util.URI, fsapi.DeleteBStoreMethod, params, result, auth)
//    if err != nil {
//        return result, err
//    }
//    return result, err
//}

//func (util *Util) ListBStoresCmd(auth *dsrpc.Auth) (*fsapi.ListBStoresResult, error) {
//    var err error
//    params := fsapi.NewListBStoresParams()
//    result := fsapi.NewListBStoresResult()
//    err = dsrpc.Exec(util.URI, fsapi.ListBStoresMethod, params, result, auth)
//    if err != nil {
//        return result, err
//    }
//    return result, err
//}
