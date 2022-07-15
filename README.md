# DStore - Distributed Store

This is a project for a decentralized data warehouse.


## Some rules and features

### Users

- The user name and password can be an random set of bytes
- Each user has his own file space
- Only the administrator can add a user
- User can be deleted by the administrator or the user himself
- A user can only be deleted if he has no files
- The user can be disabled

### Files

- Filename can be random and not limited
- But the filename passes through a unix-like normalization routine when queried
to be able to use the pseudo-directory listing
- The file upload can be interrupted, the received amount will be saved
- The listing can be made using a pattern

### In progress in search of an optimal solution:

- Each file, once written, will be distributed to a number of block servers
- Each file after writing will have redundant data, based on which the lost blocks can be recovered

## Generic draft


![](/docs/draft01.svg "Draft")


## Comman line samples

### File operation

#### Save file

```
$ dd if=/dev/urandom of=/tmp/block.bin bs=1M count=1024
1024+0 records in
1024+0 records out
1073741824 bytes transferred in 12.303218 secs (87273253 bytes/sec)

$ fstorecli -aLogin user -aPass user saveFile -local /tmp/block.bin -remote a/b/file1.bin
{
  "error": false,
  "result": {
    "file": {
      "filePath": "/a/b/file1.bin",
      "login": "user",
      "fileId": 1,
      "batchCount": 26,
      "batchSize": 5,
      "blockSize": 8388608,
      "dataSize": 1073741824,
      "createdAt": 1657878624,
      "updatedAt": 1657878624
    }
  }
}
```

#### List saved files
```
$ fstorecli -aLogin user -aPass user listFiles
{
  "error": false,
  "result": {
    "files": [
      {
        "filePath": "/a/b/file1.bin",
        "login": "user",
        "fileId": 1,
        "batchCount": 26,
        "batchSize": 5,
        "blockSize": 8388608,
        "dataSize": 1073741824,
        "createdAt": 1657878624,
        "updatedAt": 1657878624
      }
    ]
  }
}
```

#### Download file

```
$ fstorecli -aLogin user -aPass user loadFile -remote a/b/file1.bin -local /tmp/file1.bin
{
  "error": false,
  "result": {
    "file": {
      "filePath": "/a/b/file1.bin",
      "login": "user",
      "fileId": 1,
      "batchCount": 26,
      "batchSize": 5,
      "blockSize": 8388608,
      "dataSize": 1073741824,
      "createdAt": 1657878624,
      "updatedAt": 1657878624
    }
  }
}

# sha256sum /tmp/*.bin
3c9e98be6bb73cdf4c48f2f59da74033c730b8f5235c2088dc1d961c8979c43f  /tmp/block.bin
3c9e98be6bb73cdf4c48f2f59da74033c730b8f5235c2088dc1d961c8979c43f  /tmp/file1.bin
```

#### Delete file

```
$ fstorecli -aLogin user -aPass user deleteFile -path a/b/file1.bin
{
  "error": false,
  "result": {
    "file": {
      "filePath": "/a/b/file1.bin",
      "login": "user",
      "fileId": 1,
      "batchCount": 26,
      "batchSize": 5,
      "blockSize": 8388608,
      "dataSize": 1073741824,
      "createdAt": 1657878624,
      "updatedAt": 1657878624
    }
  }
}

$ fstorecli -aLogin user -aPass user deleteFile -path a/b/file1.bin
{
  "error": false,
  "result": {
    "file": null
  }
}

$ fstorecli -aLogin user -aPass user listFiles
{
  "error": false,
  "result": {
    "files": []
  }
}
```

### User operation

##### List users

```
$ fstorecli -aLogin admin -aPass admin listUsers
{
  "error": false,
  "result": {
    "users": [
      {
        "login": "admin",
        "pass": "admin",
        "role": "admin",
        "state": "enabled",
        "updatedAt": 1657835826,
        "createdAt": 1657835826
      },
      {
        "login": "user",
        "pass": "user",
        "role": "user",
        "state": "enabled",
        "updatedAt": 1657835826,
        "createdAt": 1657835826
      }
    ]
  }
}
```
#### Create user
```
$ fstorecli -aLogin admin -aPass admin addUser -login user2 -pass 12345
{
  "error": false,
  "result": {}
}

$ fstorecli -aLogin admin -aPass admin listUsers
{
  "error": false,
  "result": {
    "users": [
      ...
      {
        "login": "user2",
        "pass": "12345",
        "role": "user",
        "state": "enabled",
        "updatedAt": 1657879608,
        "createdAt": 1657879608
      }
    ]
  }
}

$ fstorecli -aLogin user2 -aPass 12345 saveFile -local ./configure -remote ./configure
{
  "error": false,
  "result": {
    "file": {
      "filePath": "/configure",
      "login": "user2",
      "fileId": 1,
      "batchCount": 1,
      "batchSize": 5,
      "blockSize": 32768,
      "dataSize": 121507,
      "createdAt": 1657879930,
      "updatedAt": 1657879930
    }
  }
}


```

#### User delete itself

```
$ fstorecli -aLogin admin -aPass admin deleteUser -login user2
{
  "error": true,
  "errorMsg": "user user2 have files",
  "result": {}
}

$ fstorecli -aLogin user2 -aPass 12345 deleteFile -path ./configure
{
  "error": false,
  "result": {
    "file": {
      "filePath": "/configure",
      "login": "user2",
      "fileId": 1,
      "batchCount": 1,
      "batchSize": 5,
      "blockSize": 32768,
      "dataSize": 121507,
      "createdAt": 1657881714,
      "updatedAt": 1657881714
    }
  }
}

$ fstorecli -aLogin user2 -aPass 12345 listFiles
{
  "error": false,
  "result": {
    "files": []
  }
}

# fstorecli -aLogin user2 -aPass 12345 deleteUser -login user2
{
  "error": false,
  "result": {}
}
```

### Rights checking samples but not all

#### Insufficient rights

```
$ fstorecli -aLogin user2 -aPass 12345 deleteUser -login user
{
  "error": true,
  "errorMsg": "user user2 have insufficient rights",
  "result": {}
}
```

#### Wrong login or password

```
$ fstorecli -aLogin user2 -aPass wrongpass deleteFile -path ./configure
{
  "error": true,
  "errorMsg": "auth mismatch",
  "result": {
    "file": null
  }
}

$ fstorecli -aLogin wronglogin -aPass 12345 deleteFile -path ./configure
{
  "error": true,
  "errorMsg": "auth error",
  "result": {
    "file": null
  }
}
```
