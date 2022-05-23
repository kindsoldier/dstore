package fdconf

import (
    "path/filepath"
    "os"
    "io/fs"
    "github.com/go-yaml/yaml"
)


const configName string = "fstore.conf"

type Config struct {
    Port        string      `json:"port"    yaml:"port"`
    ConfDir     string      `json:"confdir" yaml:"confdir"`
    DataDir     string      `json:"datadir" yaml:"datadir"`
    LogDir      string      `json:"logdir"  yaml:"logdir"`
    RunDir      string      `json:"rundir"  yaml:"rundir"`

    AccName     string      `json:"-"       yaml:"-"`
    MsgName     string      `json:"-"       yaml:"-"`
    PidName     string      `json:"-"       yaml:"-"`

    FilePerm    fs.FileMode `json:"-"       yaml:"-"`
    DirPerm     fs.FileMode `json:"-"       yaml:"-"`
}

func NewConfig() *Config {
    var config Config
    config.RunDir   = "/home/ziggi/ndstore/fstore/run"
    config.LogDir   = "/home/ziggi/ndstore/fstore/log"
    config.DataDir  = "/home/ziggi/ndstore/fstore/data"
    config.ConfDir  = "/home/ziggi/ndstore/fstore/"
    config.Port     = "5001"

    config.PidName  = "fstore.pid"
    config.MsgName  = "message.log"
    config.AccName  = "access.log"

    config.FilePerm = 0655
    config.DirPerm  = 0755

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
