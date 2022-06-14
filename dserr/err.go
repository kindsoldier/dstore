package dserr

import (
    "fmt"
    "runtime"
)

var develMode bool = true

func SetDevelMode(mode bool) {
    develMode = mode
} 

func Err(err error) error {
    if err != nil {
        switch {
            case develMode == true:
                pc, filename, line, _ := runtime.Caller(1)
                funcName := runtime.FuncForPC(pc).Name()
                err = fmt.Errorf("\n%s:%d:%s:%s", filename, line, funcName, err.Error())
            default:
                pc, _, line, _ := runtime.Caller(1)
                funcName := runtime.FuncForPC(pc).Name()
                err = fmt.Errorf(" %d:%s:%s ", funcName, line, err.Error())
        }
    }
    return err
}
