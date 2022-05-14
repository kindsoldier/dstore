
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package ndapi

const ListMethod string = "list"

type ListParams struct {
    ClusterId   int64           `json:"clusterId"`
}

type ListResult struct {
}

func NewListResult() *ListResult {
    return &ListResult{}
}
func NewListParams() *ListParams {
    return &ListParams{}
}
