/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */


package bsfunc

import (
    "io"
    "ndstore/bstore/bsapi"
    "ndstore/dsrpc"
    "ndstore/dscom"
    "ndstore/dserr"
)

const HelloMessage string = "hello"

func GetHello(uri string, auth *dsrpc.Auth) error {
    var err error

    params := bsapi.NewGetHelloParams()
    params.Message = HelloMessage
    result := bsapi.NewGetHelloResult()

    err = dsrpc.Exec(uri, bsapi.GetHelloMethod, params, result, auth)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func SaveBlock(uri string, auth *dsrpc.Auth, descr *dscom.BlockDescr, blockReader io.Reader, binSize int64) error {
    var err error
    params := bsapi.NewSaveBlockParams()
    params.FileId       = descr.FileId
    params.BatchId      = descr.BatchId
    params.BlockId      = descr.BlockId
    params.BlockType    = descr.BlockType
    params.BlockVer     = descr.BlockVer

    params.BlockSize    = descr.BlockSize
    params.DataSize     = descr.DataSize
    params.HashAlg      = descr.HashAlg
    params.HashInit     = descr.HashInit
    params.HashSum      = descr.HashSum
    result := bsapi.NewSaveBlockResult()

    err = dsrpc.Put(uri, bsapi.SaveBlockMethod, blockReader, binSize, params, result, auth)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}


func LoadBlock(uri string, auth *dsrpc.Auth, fileId, batchId, blockId int64,
                                                blockWriter io.Writer, blockType string, blockVer int64) error {
    var err error
    params := bsapi.NewLoadBlockParams()
    params.FileId       = fileId
    params.BatchId      = batchId
    params.BlockId      = blockId
    params.BlockType    = blockType
    params.BlockVer     = blockVer

    result := bsapi.NewLoadBlockResult()
    err = dsrpc.Get(uri, bsapi.LoadBlockMethod, blockWriter, params, result, auth)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func ListBlocks(uri string, auth *dsrpc.Auth) ([]*dscom.BlockDescr, error) {
    var err error
    blockDescrs := make([]*dscom.BlockDescr, 0)
    params := bsapi.NewListBlocksParams()
    result := bsapi.NewListBlocksResult()
    err = dsrpc.Exec(uri, bsapi.ListBlocksMethod, params, result, auth)
    if err != nil {
        return blockDescrs, dserr.Err(err)
    }
    blockDescrs = result.Blocks
    return blockDescrs, dserr.Err(err)
}

func LinkBlock(uri string, auth *dsrpc.Auth, fileId, batchId, blockId int64, blockType string, oldBlockVer, newBlockVer int64) error {
    var err error
    params := bsapi.NewLinkBlockParams()
    params.FileId       = fileId
    params.BatchId      = batchId
    params.BlockId      = blockId
    params.BlockType    = blockType
    params.OldBlockVer     = oldBlockVer
    params.NewBlockVer     = newBlockVer

    result := bsapi.NewLinkBlockResult()
    err = dsrpc.Exec(uri, bsapi.LinkBlockMethod, params, result, auth)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}


func DeleteBlock(uri string, auth *dsrpc.Auth, fileId, batchId, blockId int64, blockType string, blockVer int64) error {
    var err error
    params := bsapi.NewDeleteBlockParams()
    params.FileId       = fileId
    params.BatchId      = batchId
    params.BlockId      = blockId
    params.BlockType    = blockType
    params.BlockVer     = blockVer

    result := bsapi.NewDeleteBlockResult()
    err = dsrpc.Exec(uri, bsapi.DeleteBlockMethod, params, result, auth)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func BlockExists(uri string, auth *dsrpc.Auth, fileId, batchId, blockId int64, blockType string, blockVer int64) (bool, error) {
    var err error
    var exists bool
    params := bsapi.NewBlockExistsParams()
    params.FileId       = fileId
    params.BatchId      = batchId
    params.BlockId      = blockId
    params.BlockType    = blockType
    params.BlockVer     = blockVer

    result := bsapi.NewBlockExistsResult()
    err = dsrpc.Exec(uri, bsapi.BlockExistsMethod, params, result, auth)
    exists = result.Exists
    if err != nil {
        return exists, dserr.Err(err)
    }
    return exists, dserr.Err(err)
}

func CheckBlock(uri string, auth *dsrpc.Auth, fileId, batchId, blockId int64, blockType string, blockVer int64) (bool, error) {
    var err error
    var correct bool
    params := bsapi.NewCheckBlockParams()
    params.FileId       = fileId
    params.BatchId      = batchId
    params.BlockId      = blockId
    params.BlockType    = blockType
    params.BlockVer     = blockVer

    result := bsapi.NewCheckBlockResult()
    err = dsrpc.Exec(uri, bsapi.CheckBlockMethod, params, result, auth)
    correct = result.Correct
    if err != nil {
        return correct, dserr.Err(err)
    }
    return correct, dserr.Err(err)
}

func EraseAll(uri string, auth *dsrpc.Auth) error {
    var err error
    params := bsapi.NewEraseAllParams()
    result := bsapi.NewEraseAllResult()
    err = dsrpc.Exec(uri, bsapi.EraseAllMethod, params, result, auth)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
