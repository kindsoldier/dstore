
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

import (
    "dstore/dscomm/dsdescr"
)

const AddBStoreMethod string = "addBStore"

type AddBStoreParams struct {
    Address string                  `json:"address" db:"address"`
    Port    string                  `json:"port"    db:"port"`
    Login   string                  `json:"login"   db:"login"`
    Pass    string                  `json:"pass"    db:"pass"`
    State   string                  `json:"state"   db:"state"`
}

type AddBStoreResult struct {
}

func NewAddBStoreResult() *AddBStoreResult {
    return &AddBStoreResult{}
}
func NewAddBStoreParams() *AddBStoreParams {
    return &AddBStoreParams{}
}



const DeleteBStoreMethod string = "deleteBStore"

type DeleteBStoreParams struct {
    Address string                  `json:"address" db:"address"`
    Port    string                  `json:"port"    db:"port"`
}

type DeleteBStoreResult struct {
}

func NewDeleteBStoreResult() *DeleteBStoreResult {
    return &DeleteBStoreResult{}
}
func NewDeleteBStoreParams() *DeleteBStoreParams {
    return &DeleteBStoreParams{}
}


const ListBStoresMethod string = "listBStores"

type ListBStoresParams struct {
    Regular     string              `json:"regular"`
}

type ListBStoresResult struct {
    BStores  []*dsdescr.BStore      `json:"bStores,omitempty"`
}

func NewListBStoresResult() *ListBStoresResult {
    return &ListBStoresResult{}
}
func NewListBStoresParams() *ListBStoresParams {
    return &ListBStoresParams{}
}



const UpdateBStoreMethod string = "updateBStore"

type UpdateBStoreParams struct {
    Address string                  `json:"address" db:"address"`
    Port    string                  `json:"port"    db:"port"`
    Login   string                  `json:"login"   db:"login"`
    Pass    string                  `json:"pass"    db:"pass"`
    State   string                  `json:"state"   db:"state"`
}

type UpdateBStoreResult struct {
}

func NewUpdateBStoreResult() *UpdateBStoreResult {
    return &UpdateBStoreResult{}
}
func NewUpdateBStoreParams() *UpdateBStoreParams {
    return &UpdateBStoreParams{}
}


const CheckBStoreMethod string = "checkBStore"
type CheckBStoreParams struct {
    Address string                  `json:"address" db:"address"`
    Port    string                  `json:"port"    db:"port"`
    Login   string                  `json:"login"`
    Pass    string                  `json:"pass"`
}
type CheckBStoreResult struct {
    Match   bool                    `json:"match"`
}

func NewCheckBStoreResult() *CheckBStoreResult {
    return &CheckBStoreResult{}
}
func NewCheckBStoreParams() *CheckBStoreParams {
    return &CheckBStoreParams{}
}
