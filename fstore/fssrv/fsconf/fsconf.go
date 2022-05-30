package fdconf

import (
    "fmt"
    "path/filepath"
    "os"
    "io/fs"
    "github.com/go-yaml/yaml"
)


const configName string = "fstore.conf"

type Config struct {
    Port        string      `json:"port"    yaml:"port"`
    ConfDir     string      `json:"confDir" yaml:"confDir"`
    DataDir     string      `json:"dataDir" yaml:"dataDir"`
    LogDir      string      `json:"logDir"  yaml:"logDir"`
    RunDir      string      `json:"runDir"  yaml:"runDir"`

    AccName     string      `json:"-"       yaml:"-"`
    MsgName     string      `json:"-"       yaml:"-"`
    PidName     string      `json:"-"       yaml:"-"`

    DbName      string      `json:"dbName"  yaml:"dbName"`
    DbHost      string      `json:"dbHost"  yaml:"dbHost"`
    DbUser      string      `json:"dbUser"  yaml:"dbUser"`
    DbPass      string      `json:"dbPass"  yaml:"dbPass"`

    FilePerm    fs.FileMode `json:"-"       yaml:"-"`
    DirPerm     fs.FileMode `json:"-"       yaml:"-"`
}

func NewConfig() *Config {
    var config Config
    config.RunDir   = "/home/ziggi/ndstore/fstore/run"
    config.LogDir   = "/home/ziggi/ndstore/fstore/log"
    config.DataDir  = "/home/ziggi/ndstore/fstore/data"
    config.ConfDir  = "/home/ziggi/ndstore/fstore/"
    config.Port     = "5200"

    config.PidName  = "fstore.pid"
    config.MsgName  = "message.log"
    config.AccName  = "access.log"

    config.FilePerm = 0655
    config.DirPerm  = 0755

    return &config
}

func (conf *Config) Read() error {
    var err error
    filename := filepath.Join(conf.ConfDir, configName)
    confData, err := os.ReadFile(filename)
    err = yaml.Unmarshal(confData, conf)
    if err != nil {
        return err
    }
    return err
}

func (conf *Config) GetDBPath() string {
    return fmt.Sprintf("postgres://%s:%s@%s/$s",
                    conf.DbUser, conf.DbPass, conf.DbHost, conf.DbName)
}
