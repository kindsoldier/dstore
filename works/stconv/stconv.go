package dsstconv

import (
    "fmt"
    "reflect"
    "strconv"
)


func Pack(object interface{}) (map[string]string, error) {
    var err error
    resMap := make(map[string]string)
    rValue := reflect.ValueOf(object)
    rType := reflect.Indirect(rValue).Type()

    if rType.Kind() == reflect.Struct {
        elem := rValue.Elem()
        for i:=0; i < elem.NumField(); i++ {

            vName     := elem.Type().Field(i).Name
            jsonTag, ok := rType.Field(i).Tag.Lookup("json")
            if ok {
                vName = jsonTag
            }
            vType     := elem.Type().Field(i).Type
            vField    := elem.Field(i)

            var sValue string

            switch vType.Kind() {
                case reflect.String:
                    sValue = vField.String()
                case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
                    sValue = strconv.FormatInt(vField.Int(), 10)
                case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
                    sValue = strconv.FormatUint(vField.Uint(), 10)
                case reflect.Float32:
                    sValue = strconv.FormatFloat(vField.Float(), 'E', -1, 32)
                case reflect.Float64:
                    sValue = strconv.FormatFloat(vField.Float(), 'E', -1, 64)
                case reflect.Bool:
                    sValue = strconv.FormatBool(vField.Bool())
            }
            resMap[fmt.Sprintf("%s.%s", vName, vType)] = sValue
        }
    }
    return resMap, err
}

func Unpack(resMap map[string]string, object interface{}) error {
    var err error
    refValue := reflect.ValueOf(object)
    refType := reflect.Indirect(refValue).Type()

    if refType.Kind() == reflect.Struct {
        elem := refValue.Elem()

        for i:=0; i < elem.NumField(); i++ {

            vName     := elem.Type().Field(i).Name
            jsonTag, ok := refType.Field(i).Tag.Lookup("json")
            if ok {
                vName = jsonTag
            }
            tField  := elem.Type().Field(i)
            vType   := elem.Type().Field(i).Type
            vField  := refValue.Elem().Field(i)
            vKey := fmt.Sprintf("%s.%s", vName, vType)

            switch vType.Kind() {
                case reflect.String:
                    sValue, ok := resMap[vKey]
                    if ok {
                        elem.Field(i).SetString(sValue)
                    }
                case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
                    sValue, ok := resMap[vKey]
                    if ok {
                        bits := tField.Type.Bits()
                        iVal, err := strconv.ParseInt(sValue, 10, bits)
                        if err == nil {
                            vField.SetInt(int64(iVal))
                        }
                    }
                case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
                    sValue, ok := resMap[vKey]
                    if ok {
                        bits := tField.Type.Bits()
                        iVal, err := strconv.ParseUint(sValue, 10, bits)
                        if err == nil {
                            vField.SetUint(uint64(iVal))
                        }
                    }
                case reflect.Float32, reflect.Float64:
                    sValue, ok := resMap[vKey]
                    if ok {
                        bits := tField.Type.Bits()
                        fVal, err := strconv.ParseFloat(sValue, bits)
                        if err == nil {
                            vField.SetFloat(float64(fVal))
                        }
                    }
                case reflect.Bool:
                    sValue, ok := resMap[vKey]
                    if ok {
                        bVal, err := strconv.ParseBool(sValue)
                        if err == nil {
                            vField.SetBool(bool(bVal))
                        }
                    }
            }
        }
    }
    return err
}
