/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package bsfun

import (
    "io"
    "dstore/bstore/bsapi"
    "dstore/dscomm/dsrpc"
    "dstore/dscomm/dsdescr"
    "dstore/dscomm/dserr"
)

func GetStatus(uri string, auth *dsrpc.Auth) error {
    var err error
    params := bsapi.NewGetStatusParams()
    result := bsapi.NewGetStatusResult()
    err = dsrpc.Exec(uri, bsapi.GetStatusMethod, params, result, auth)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func SaveBlock(uri string, auth *dsrpc.Auth, descr *dsdescr.Block, blockReader io.Reader, binSize int64) error {
    var err error
    params := bsapi.NewSaveBlockParams()
    params.FileId       = descr.FileId
    params.BatchId      = descr.BatchId
    params.BlockType    = descr.BlockType
    params.BlockId      = descr.BlockId

    params.BlockSize    = descr.BlockSize
    //params.DataSize     = descr.DataSize
    //params.HashAlg      = descr.HashAlg
    //params.HashInit     = descr.HashInit
    //params.HashSum      = descr.HashSum
    result := bsapi.NewSaveBlockResult()

    err = dsrpc.Put(uri, bsapi.SaveBlockMethod, blockReader, binSize, params, result, auth)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func LoadBlock(uri string, auth *dsrpc.Auth, fileId, batchId, blockType, blockId int64, blockWriter io.Writer) error {
    var err error
    params := bsapi.NewLoadBlockParams()
    params.FileId       = fileId
    params.BatchId      = batchId
    params.BlockType    = blockType
    params.BlockId      = blockId

    result := bsapi.NewLoadBlockResult()
    err = dsrpc.Get(uri, bsapi.LoadBlockMethod, blockWriter, params, result, auth)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func ListBlocks(uri string, auth *dsrpc.Auth, fileId int64) ([]*dsdescr.Block, error) {
    var err error
    blockDescrs := make([]*dsdescr.Block, 0)
    params := bsapi.NewListBlocksParams()
    params.FileId = fileId
    result := bsapi.NewListBlocksResult()
    err = dsrpc.Exec(uri, bsapi.ListBlocksMethod, params, result, auth)
    if err != nil {
        return blockDescrs, dserr.Err(err)
    }
    blockDescrs = result.Blocks
    return blockDescrs, dserr.Err(err)
}

func DeleteBlock(uri string, auth *dsrpc.Auth, fileId, batchId, blockType, blockId int64) error {
    var err error
    params := bsapi.NewDeleteBlockParams()
    params.FileId       = fileId
    params.BatchId      = batchId
    params.BlockType    = blockType
    params.BlockId      = blockId

    result := bsapi.NewDeleteBlockResult()
    err = dsrpc.Exec(uri, bsapi.DeleteBlockMethod, params, result, auth)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
