package main

import (
    "path/filepath"
    "os"

    "github.com/go-yaml/yaml"
)


const configName string = "dcnode.conf"

type Config struct {
    Port        string      `json:"port"    yaml:"port"`
    ConfDir     string      `json:"confdir" yaml:"confdir"`
    DataDir     string      `json:"datadir" yaml:"datadir"`
    LogDir      string      `json:"logdir"  yaml:"logdir"`
    RunDir      string      `json:"rundir"  yaml:"rundir"`

    AccName     string      `json:"-"       yaml:"-"`
    MsgName     string      `json:"-"       yaml:"-"`
    PidName     string      `json:"-"       yaml:"-"`
}

func NewConfig() *Config {
    var config Config
    config.RunDir   = "/home/ziggi/dcstore/dcnode/run"
    config.LogDir   = "/home/ziggi/dcstore/dcnode/log"
    config.DataDir  = "/home/ziggi/dcstore/dcnode/data"
    config.ConfDir  = "/home/ziggi/dcstore/dcnode/"
    config.Port     = "5001"

    config.PidName  = "dcnode.pid"
    config.MsgName  = "message.log"
    config.AccName  = "access.log"

    return &config
}

func (this *Config) Read() error {
    var err error
    filename := filepath.Join(this.ConfDir, configName)
    confData, err := os.ReadFile(filename)
    err = yaml.Unmarshal(confData, this)
    if err != nil {
        return err
    }
    return err
}
