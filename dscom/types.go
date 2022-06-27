/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dscom

import (
    "time"
    "strings"
)

type UnixTime int64

func (ut *UnixTime) MarshalJSON() ([]byte, error) {
    var err error
    tsTime := time.Unix(int64(*ut), 0)
    tsString := tsTime.Format(time.RFC3339)
    tsBytes := []byte(`"` + tsString + `"`)
    return tsBytes, err
}

func (ut *UnixTime) UnmarshalJSON(src []byte) error {
    var err error
    tsString := strings.Trim(string(src), `"`)
    ts, err := time.Parse(time.RFC3339, tsString)
    if err != nil {
        return err
    }
    *ut = UnixTime(ts.Unix())
    return err
}
