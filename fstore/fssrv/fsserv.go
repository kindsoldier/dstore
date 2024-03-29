/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package main

import (
    "flag"
    "fmt"
    "io/fs"
    "os"
    "os/signal"
    "os/user"
    "path/filepath"
    "strconv"
    "syscall"
    "io"
    "time"

    "dstore/fstore/fsapi"
    "dstore/fstore/fssrv/fscont"
    "dstore/fstore/fssrv/fsreg"
    "dstore/fstore/fssrv/fstore"

    "dstore/dscomm/dskvdb"
    "dstore/dscomm/dslog"
    "dstore/dscomm/dsrpc"
    "dstore/dscomm/dserr"
    "dstore/dscomm/dsalloc"
    "dstore/dscomm/dsinter"
)

const successExit   int = 0
const errorExit     int = 1

func main() {
    var err error
    server := NewServer()

    dserr.SetDevelMode(false)
    dserr.SetDebugMode(false)

    err = server.Execute()
    if err != nil {
        dslog.LogError("config error:", err)
        os.Exit(errorExit)
    }
}

type Server struct {
    Params  *Config
    Backgr  bool
    fileIdAlloc dsinter.Alloc
    serv    *dsrpc.Service
}

func (server *Server) Execute() error {
    var err error

    err = server.ReadConf()
    if err != nil {
        return err
    }
    err = server.GetOptions()
    if err != nil {
        return err
    }

    err = server.PrepareEnv()
    if err != nil {
        return err
    }
    if server.Backgr {
        err = server.ForkCmd()
        if err != nil {
            return err
        }
        err = server.CloseIO()
        if err != nil {
            return err
        }

        err = server.ChangeUid()
        if err != nil {
            return err
        }
    }

    err = server.RedirLog()
    if err != nil {
        return err
    }
    err = server.SavePid()
    if err != nil {
        return err
    }
    err = server.SetSHandler()
    if err != nil {
        return err
    }

    err = server.RunService()
    if err != nil {
        return err
    }
    return err
}


func NewServer() *Server {
    var server Server
    server.Params = NewConfig()
    server.Backgr = false
    return &server
}

func (server *Server) ReadConf() error {
    var err error
    err = server.Params.Read()
    if err != nil {
        return err
    }
    return err
}

func (server *Server) GetOptions() error {
    var err error
    exeName := filepath.Base(os.Args[0])

    flag.StringVar(&server.Params.RunDir, "runDir", server.Params.RunDir, "run direcory")
    flag.StringVar(&server.Params.LogDir, "logDir", server.Params.LogDir, "log direcory")
    flag.StringVar(&server.Params.DataDir, "dataDir", server.Params.DataDir, "data directory")

    flag.StringVar(&server.Params.Port, "port", server.Params.Port, "listen port")
    flag.BoolVar(&server.Backgr, "daemon", server.Backgr, "run as daemon")

    help := func() {
        fmt.Println("")
        fmt.Printf("usage: %s [option]\n", exeName)
        fmt.Println("")
        fmt.Println("options:")
        flag.PrintDefaults()
        fmt.Println("")
    }
    flag.Usage = help
    flag.Parse()

    return err
}


func (server *Server) ForkCmd() error {
    const successExit int = 0
    var keyEnv string = "IMX0LTSELMRF8K"
    var err error


    _, isChild := os.LookupEnv(keyEnv)
    switch  {
        case !isChild:
            os.Setenv(keyEnv, "TRUE")

            procAttr := syscall.ProcAttr{}
            cwd, err := os.Getwd()
            if err != nil {
                    return err
            }
            var sysFiles = make([]uintptr, 3)
            sysFiles[0] = uintptr(syscall.Stdin)
            sysFiles[1] = uintptr(syscall.Stdout)
            sysFiles[2] = uintptr(syscall.Stderr)

            procAttr.Files = sysFiles
            procAttr.Env = os.Environ()
            procAttr.Dir = cwd

            _, err = syscall.ForkExec(os.Args[0], os.Args, &procAttr)
            if err != nil {
                return err
            }
            os.Exit(successExit)
        case isChild:
            _, err = syscall.Setsid()
            if err != nil {
                    return err
            }
    }
    os.Unsetenv(keyEnv)
    return err
}

func (server *Server) ChangeUid() error {
    var err error

    username := server.Params.SrvUser

    currUid := syscall.Getuid()
    if currUid != 0 {
        return err
    }

    userDescr, err := user.Lookup(username)
    if err != nil {
        err = fmt.Errorf("no username %s found, err: %v", username, err)
        return err
    }

    newGid, err := strconv.Atoi(userDescr.Gid)
    if err != nil {
        err = fmt.Errorf("cannot convert gid, err: %v", err)
        return err
    }
    err = syscall.Setgid(newGid)
    if err != nil {
        err = fmt.Errorf("cannot change gid, err: %v", err)
        return err
    }
    currGid := syscall.Getgid()
    if currGid != newGid {
        err = fmt.Errorf("unable to change gid for unknown reason")
        return err
    }

    newUid, err := strconv.Atoi(userDescr.Uid)
    if err != nil {
        err = fmt.Errorf("cannot convert uid, err: %v", err)
        return err
    }

    runDir := server.Params.RunDir
    err = os.Chown(runDir, newUid, newGid)
    if err != nil {
            return err
    }

    logDir := server.Params.LogDir
    err = os.Chown(logDir, newUid, newGid)
    if err != nil {
            return err
    }

    dataDir := server.Params.DataDir
    err = os.Chown(dataDir, newUid, newGid)
    if err != nil {
            return err
    }

    err = syscall.Setuid(newUid)
    if err != nil {
        err = fmt.Errorf("cannot change uid, err: %v", err)
        return err
    }

    err = syscall.Seteuid(newUid)
    if err != nil {
        err = fmt.Errorf("cannot change euid, err: %v", err)
        return err
    }


    currUid = syscall.Getuid()
    if currUid != newUid {
        err = fmt.Errorf("unable to change uid for unknown reason")
        return err
    }
    return err
}



func (server *Server) PrepareEnv() error {
    var err error

    var runDirPerm fs.FileMode = server.Params.DirPerm
    var logDirPerm fs.FileMode = server.Params.DirPerm
    var dataDirPerm fs.FileMode = server.Params.DirPerm

    runDir := server.Params.RunDir
    err = os.MkdirAll(runDir, runDirPerm)
    if err != nil {
            return err
    }
    err = os.Chmod(runDir, runDirPerm)
    if err != nil {
            return err
    }

    logDir := server.Params.LogDir
    err = os.MkdirAll(logDir, logDirPerm)
    if err != nil {
            return err
    }
    err = os.Chmod(logDir, logDirPerm)
    if err != nil {
            return err
    }

    dataDir := server.Params.DataDir
    err = os.MkdirAll(dataDir, dataDirPerm)
    if err != nil {
            return err
    }
    err = os.Chmod(dataDir, dataDirPerm)
    if err != nil {
            return err
    }
    return err
}

func (server *Server) SavePid() error {
    var err error

    var pidFilePerm fs.FileMode = server.Params.DirPerm

    pidFile := filepath.Join(server.Params.RunDir, server.Params.PidName)

    openMode := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
    file, err := os.OpenFile(pidFile, openMode, pidFilePerm)
    if err != nil {
            return err
    }
    defer file.Close()

    pid := os.Getpid()
    _, err = file.WriteString(strconv.Itoa(pid))
    if err != nil {
            return err
    }
    err = os.Chmod(pidFile, pidFilePerm)
    if err != nil {
            return err
    }
    file.Sync()
    return err
}

func (server *Server) RedirLog() error {
    var err error

    var logFilePerm fs.FileMode = server.Params.FilePerm

    logOpenMode := os.O_WRONLY|os.O_CREATE|os.O_APPEND
    msgFileName := filepath.Join(server.Params.LogDir, server.Params.MsgName)
    msgFile, err := os.OpenFile(msgFileName, logOpenMode, logFilePerm)
    if err != nil {
            return err
    }

    logWriter := io.MultiWriter(os.Stdout, msgFile)
    dslog.SetOutput(logWriter)
    dsrpc.SetMessageWriter(logWriter)

    accFileName := filepath.Join(server.Params.LogDir, server.Params.AccName)
    accFile, err := os.OpenFile(accFileName, logOpenMode, logFilePerm)
    if err != nil {
            return err
    }

    accWriter := io.MultiWriter(os.Stdout, accFile)
    dsrpc.SetAccessWriter(accWriter)
    return err
}

func (server *Server) CloseIO() error {
    var err error
    file, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
    if err != nil {
            return err
    }
    err = syscall.Dup2(int(file.Fd()), int(os.Stdin.Fd()))
    if err != nil {
            return err
    }
    err = syscall.Dup2(int(file.Fd()), int(os.Stdout.Fd()))
    if err != nil {
            return err
    }
    err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
    if err != nil {
            return err
    }
    return err
}


func (server *Server) SetSHandler() error {
    var err error
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGSTOP,
                                    syscall.SIGTERM, syscall.SIGQUIT)

    handler := func() {
        var err error
        for {
            dslog.LogInfo("signal handler start")
            sig := <-sigs
            dslog.LogInfo("received signal", sig.String())

            switch sig {
                case syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGQUIT:
                    dslog.LogInfo("exit by signal", sig.String())
                    server.StopAll()
                    os.Exit(successExit)

                case syscall.SIGHUP:
                    switch {
                        case server.Backgr:
                            dslog.LogInfo("restart server")

                            err = server.StopAll()
                            if err != nil {
                                dslog.LogError("stop all error:", err)
                            }
                            dslog.LogInfo("fork")
                            err = server.ForkCmd()
                            if err != nil {
                                dslog.LogError("fork error:", err)
                            }
                        default:
                            server.StopAll()
                            os.Exit(successExit)
                    }
            }
        }
    }
    go handler()
    return err
}

func (server *Server) RunService() error {
    var err error

    filePerm    := server.Params.FilePerm
    dirPerm     := server.Params.DirPerm
    dataDir     := server.Params.DataDir

    //develMode   := false
    //debugMode   := false

    develMode   := server.Params.DevelMode
    debugMode   := server.Params.DebugMode

    dslog.SetDebugMode(debugMode)

    //dserr.SetDevelMode(develMode)
    //dserr.SetDebugMode(debugMode)

    //dsrpc.SetDevelMode(develMode)
    //dsrpc.SetDebugMode(debugMode)

    db, err := dskvdb.OpenDB(dataDir, "storedb")
    if err != nil {
        return err
    }
    reg, err := fsreg.NewReg(db)
    if err != nil {
        return err
    }
    server.fileIdAlloc, err = dsalloc.OpenAlloc(db, []byte("fileIds"))
    if err != nil {
        return err
    }
    go server.fileIdAlloc.Syncer()

    store, err := fstore.NewStore(dataDir, reg, server.fileIdAlloc)
    if err != nil {
        return err
    }

    store.SetFilePerm(filePerm)
    store.SetDirPerm(dirPerm)

    err = store.SeedUsers()
    if err != nil {
        return err
    }
    err = store.SeedBStores()
    if err != nil {
        return err
    }

    contr, err := fscont.NewContr(store)
    if err != nil {
        return err
    }

    dslog.LogInfof("dataDir is %s", server.Params.DataDir)
    dslog.LogInfof("logDir is %s", server.Params.LogDir)
    dslog.LogInfof("runDir is %s", server.Params.RunDir)

    server.serv = dsrpc.NewService()

    if debugMode || develMode {
        server.serv.PreMiddleware(dsrpc.LogRequest)
    }
    server.serv.PreMiddleware(contr.AuthMidware(debugMode))

    server.serv.Handler(fsapi.SaveFileMethod, contr.SaveFileHandler)
    server.serv.Handler(fsapi.LoadFileMethod, contr.LoadFileHandler)
    server.serv.Handler(fsapi.FileStatsMethod, contr.FileStatsHandler)
    server.serv.Handler(fsapi.ListFilesMethod, contr.ListFilesHandler)
    server.serv.Handler(fsapi.DeleteFileMethod, contr.DeleteFileHandler)
    server.serv.Handler(fsapi.EraseFilesMethod, contr.EraseFilesHandler)

    server.serv.Handler(fsapi.AddUserMethod, contr.AddUserHandler)
    server.serv.Handler(fsapi.CheckUserMethod, contr.CheckUserHandler)
    server.serv.Handler(fsapi.UpdateUserMethod, contr.UpdateUserHandler)
    server.serv.Handler(fsapi.ListUsersMethod, contr.ListUsersHandler)
    server.serv.Handler(fsapi.DeleteUserMethod, contr.DeleteUserHandler)

    server.serv.Handler(fsapi.AddBStoreMethod, contr.AddBStoreHandler)
    server.serv.Handler(fsapi.CheckBStoreMethod, contr.CheckBStoreHandler)
    server.serv.Handler(fsapi.UpdateBStoreMethod, contr.UpdateBStoreHandler)
    server.serv.Handler(fsapi.ListBStoresMethod, contr.ListBStoresHandler)
    server.serv.Handler(fsapi.DeleteBStoreMethod, contr.DeleteBStoreHandler)

    server.serv.Handler(fsapi.GetStatusMethod, contr.GetStatusHandler)

    //if debugMode || develMode {
    //    server.serv.PostMiddleware(dsrpc.LogResponse)
    //}


    logAccess := func(context *dsrpc.Context) error {
        var err error
        execTime := time.Since(context.Start())
        login := string(context.AuthIdent())
        dslog.LogInfo(context.RemoteHost(), login, context.Method(),
                            context.ReqRpcSize(), context.ReqBinSize(),
                            context.ResRpcSize(), context.ResBinSize(),
                            execTime)
        return err
    }

    server.serv.PostMiddleware(logAccess)

    listenParam := fmt.Sprintf(":%s", server.Params.Port)
    err = server.serv.Listen(listenParam)
    if err != nil {
        return err
    }
    return err
}

func (server *Server) StopAll() error {
    var err error
    dslog.LogInfo("stop processes")
    if server.fileIdAlloc != nil {
        server.fileIdAlloc.Stop()
    }
    if server.serv != nil {
        server.serv.Stop()
    }
    return err
}
