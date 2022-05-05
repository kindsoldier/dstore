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
    "time"
    "io"

    "dcstore/dcnode/config"
    "dcstore/dclog"
)

const successExit   int = 0
const errorExit     int = 1

type Server struct {
    config      *config.Config
    background  bool
    logFile     *os.File
}


func main() {
    var err error
    server := NewServer()

    err = server.Execute()
    if err != nil {
        dclog.LogError("config error:", err)
        os.Exit(errorExit)
    }
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

    if server.background {
        err = server.ForkCmd()
        if err != nil {
            return err
        }
    }
    err = server.CloseIO()
    if err != nil {
        return err
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
    server.config = config.NewConfig()
    server.background = false
    return &server
}

func (server *Server) ReadConf() error {
    var err error
    err = server.config.Read()
    if err != nil {
        return err
    }
    return err
}

func (server *Server) GetOptions() error {
    var err error
    exeName := filepath.Base(os.Args[0])

    flag.StringVar(&server.config.Port, "p", server.config.Port, "listen port")
    flag.BoolVar(&server.background, "d", server.background, "run as daemon")

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

func (server *Server) SavePid() error {
    var err error

    const runDirPerm fs.FileMode = 0755
    const pidFilePerm fs.FileMode = 0644

    pidFile := filepath.Join(server.config.RunDir, server.config.PidName)
    runDir := server.config.RunDir

    err = os.MkdirAll(runDir, runDirPerm)
    if err != nil {
            return err
    }
    err = os.Chmod(runDir, runDirPerm)
    if err != nil {
            return err
    }

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

    const logDirPerm fs.FileMode = 0755
    const logFilePerm fs.FileMode = 0644

    logFile := filepath.Join(server.config.LogDir, server.config.LogName)
    logDir := server.config.LogDir

    err = os.MkdirAll(logDir, logDirPerm)
    if err != nil {
            return err
    }
    err = os.Chmod(logDir, logDirPerm)
    if err != nil {
            return err
    }
    openMode := os.O_WRONLY | os.O_CREATE | os.O_APPEND
    file, err := os.OpenFile(logFile, openMode, logFilePerm)
    if err != nil {
            return err
    }
    server.logFile = file
    writer := io.MultiWriter(os.Stdout, file)
    dclog.SetOutput(writer)

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
                        case server.background:
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

func (server *Server) RunService() error {
    var err error

    for {
        dclog.LogDebug("run")
        time.Sleep(1 * time.Second)
    }
    return err
}

func (server *Server) StopAll() error {
    var err error
    dclog.LogInfo("stop processes")
    if server.logFile != nil {
        server.logFile.Close()
    }
    return err
}
