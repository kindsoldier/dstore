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
    "path/filepath"
    "strconv"
    "syscall"
    "io"

    "dcstore/dcnode/ndapi"
    "dcstore/dclog"
    "dcstore/dcrpc"
    "dcstore/dcnode/ndsrv/ndconf"
    "dcstore/dcnode/ndsrv/ndcontr"
    "dcstore/dcnode/ndsrv/ndstore"
)

const successExit   int = 0
const errorExit     int = 1


func main() {
    var err error
    server := NewServer()

    err = server.Execute()
    if err != nil {
        dclog.LogError("config error:", err)
        os.Exit(errorExit)
    }
}

type Server struct {
    Params  *ndconf.Config
    Backgr  bool
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
    server.Params = ndconf.NewConfig()
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

    flag.StringVar(&server.Params.RunDir, "run", server.Params.RunDir, "run direcory")
    flag.StringVar(&server.Params.LogDir, "log", server.Params.LogDir, "log direcory")
    flag.StringVar(&server.Params.DataDir, "data", server.Params.DataDir, "data directory")

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
    dclog.SetOutput(logWriter)
    dcrpc.SetMessageWriter(logWriter)

    accFileName := filepath.Join(server.Params.LogDir, server.Params.AccName)
    accFile, err := os.OpenFile(accFileName, logOpenMode, logFilePerm)
    if err != nil {
            return err
    }

    accWriter := io.MultiWriter(os.Stdout, accFile)
    dcrpc.SetAccessWriter(accWriter)
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
            dclog.LogInfo("signal handler start")
            sig := <-sigs
            dclog.LogInfo("received signal", sig.String())

            switch sig {
                case syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP:
                    dclog.LogInfo("exit by signal", sig.String())
                    server.StopAll()
                    os.Exit(successExit)

                case syscall.SIGHUP:
                    switch {
                        case server.Backgr:
                            dclog.LogInfo("restart server")

                            err = server.StopAll()
                            if err != nil {
                                dclog.LogError("stop all error:", err)
                            }
                            err = server.ForkCmd()
                            if err != nil {
                                dclog.LogError("fork error:", err)
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

func (server *Server) StopAll() error {
    var err error
    dclog.LogInfo("stop processes")
    //if server.logFile != nil {
    //    server.logFile.Close()
    //}
    return err
}


func (server *Server) RunService() error {
    var err error

    serv := dcrpc.NewService()

    contr := ndcontr.NewContr()
    dclog.LogDebug("data dir is", server.Params.DataDir)
    store := ndstore.NewStore(server.Params.DataDir)
    store.SetPerm(server.Params.DirPerm, server.Params.FilePerm)
    err = store.OpenReg()
    if err != nil {
        return err
    }
    contr.Store = store

    serv.Handler(ndapi.HelloMethod, contr.HelloHandler)
    serv.Handler(ndapi.SaveMethod, contr.SaveHandler)
    serv.Handler(ndapi.LoadMethod, contr.LoadHandler)
    serv.Handler(ndapi.DeleteMethod, contr.DeleteHandler)
    serv.Handler(ndapi.ListMethod, contr.ListHandler)

    serv.PreMiddleware(dcrpc.LogRequest)
    serv.PostMiddleware(dcrpc.LogResponse)
    serv.PostMiddleware(dcrpc.LogAccess)

    listenParam := fmt.Sprintf(":%s", server.Params.Port)
    err = serv.Listen(listenParam)
    if err != nil {
        return err
    }

    return err
}
