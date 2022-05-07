/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */


package main

import (
    "fmt"
    "flag"
    "os"
    "path/filepath"

    "dcstore/dcnode/ndapi"
    "dcstore/dcrpc"
)

func main() {
    var err error
    util := NewUtil()
    err = util.Exec()
    if err != nil {
        fmt.Printf("exec error: %s\n", err)
    }
}

type Util struct {
    Port    string
    Address string
    Message string
}

func NewUtil() *Util {
    var util Util
    util.Port = "5001"
    util.Address = "127.0.0.1"
    util.Message = "hello world!"
    return &util
}

func (util *Util) Exec() error {
    var err error
    err = util.GetOpt()
    if err != nil {
        return err
    }
    err = util.Run()
    if err != nil {
        return err
    }
    return err
}

const helloCmd string = "hello"

func (util *Util) GetOpt() error {
    var err error

    exeName := filepath.Base(os.Args[0])

    flag.StringVar(&util.Port, "p", util.Port, "port")
    flag.StringVar(&util.Address, "a", util.Address, "address")

    help := func() {
        fmt.Println("")
        fmt.Printf("Usage: %s [option] command [command option]\n", exeName)
        fmt.Printf("\n")
        fmt.Printf("Command list: hello\n")
        fmt.Printf("\n")
        fmt.Printf("Global options:\n")
        flag.PrintDefaults()
        fmt.Printf("\n")
    }
    flag.Usage = help
    flag.Parse()

    args := flag.Args()

    if len(args) == 0 {
        args = append(args, helloCmd)
    }

    var subCommand string
    var subArgs []string
    if len(args) > 0 {
        subCommand = args[0]
        subArgs = args[1:]
    }
    switch subCommand {
        case helloCmd:
            flagSet := flag.NewFlagSet(helloCmd, flag.ContinueOnError)
            flagSet.StringVar(&util.Message, "m", util.Message, "message")
            flagSet.Usage = func() {
                fmt.Printf("\n")
                fmt.Printf("Usage: %s [global options] %s [command options]\n", exeName, subCommand)
                fmt.Printf("\n")
                fmt.Printf("The command options:\n")
                flagSet.PrintDefaults()
                fmt.Printf("\n")
            }
            flagSet.Parse(subArgs)
        default:
            help()
    }
    return err
}

func (util *Util) Run() error {
    var err error
    uri := fmt.Sprintf("%s:%s", util.Address, util.Port)

    params := ndapi.NewHelloParams()
    params.Message = util.Message
    result := ndapi.NewHelloResult()

    err = dcrpc.Exec(uri, ndapi.HelloMethod, params, result, nil)
    if err != nil {
        return err
    }
    fmt.Printf("result: %s\n", result.Message)
    return err
}
