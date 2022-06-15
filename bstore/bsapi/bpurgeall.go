
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

const PurgeAllMethod string = "purgeAll"

type PurgeAllParams struct {
}

type PurgeAllResult struct {
}

func NewPurgeAllResult() *PurgeAllResult {
    return &PurgeAllResult{}
}
func NewPurgeAllParams() *PurgeAllParams {
    return &PurgeAllParams{}
}
