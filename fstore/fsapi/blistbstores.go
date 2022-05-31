
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

import (
    "ndstore/dscom"
)

const ListBStoresMethod string = "listBStores"

type ListBStoresParams struct {
}

type ListBStoresResult struct {
    BStores  []*dscom.BStoreDescr     `json:"bStores,omitempty"`
}

func NewListBStoresResult() *ListBStoresResult {
    return &ListBStoresResult{}
}
func NewListBStoresParams() *ListBStoresParams {
    return &ListBStoresParams{}
}
