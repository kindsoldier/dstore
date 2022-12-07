# DStore - Distributed Store


This is a project for a decentralized data warehouse with some resistance to failure.

### Status

- The project is implemented in experimental version 1.
- The development of version 2 has been started, the file storage service and the block storage service have been implemented.
- Suspended for now in search of better or optimal solutions for synchronization of shared objects.

### Статус

- Проект реализован в экспериментальной версии 1.
- Начата разработка версии 2, реализованы сервис хранения файлов и сервис хранения блоков.
- Пока приостановлен в поисках лучших или оптимальных решений для синхронизации разделяемых объектов и "сборки мусора".


### In progress in search of an optimal solution:

- Each file, once written, will be distributed to a number of block servers
- Each file after writing will have redundant data, based on which the lost blocks can be recovered


## Some rules and features

### Generic

- Data block services are required for full functionality.
- If there are none, the file will be saved only on file service.

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


## Generic draft

![](/docs/draft01.svg "Draft")


## Comman line samples

All subcommand have appropriate service API.

### File operation

#### Saving file

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

#### File listing, all or with filters

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

#### File stats, all or with filters

```
$ fstorecli fileStats
{
  "error": false,
  "result": {
    "count": 324372,
    "usage": 14747643654
  }
}
```

#### File erase, all or with filters

```
# fstorecli -aLogin user -aPass user eraseFiles -erase -glob '/*/test.txt'
{
  "error": false,
  "result": {
    "files": [
      {
        "filePath": "/b/foo/bar/test.txt",
        "login": "user",
        "fileId": 294367,
        "batchCount": 1,
        "batchSize": 1,
        "blockSize": 16384,
        "dataSize": 289,
        "createdAt": 1658841733,
        "updatedAt": 1658841733
      },
      {
        "filePath": "/b/some/bare/alloc/test.txt",
        "login": "user",
        "fileId": 164256,
        "batchCount": 1,
        "batchSize": 1,
        "blockSize": 16384,
        "dataSize": 289,
        "createdAt": 1658841759,
        "updatedAt": 1658841759
      }
    ]
  }
}

```


#### File listing, stats or erasing with filters


```
$ fstorecli fileStats -help

Usage: fstorecli [global options] fileStats [command options]

The command options: none
  -glob string
        glob pattern
  -patt string
        shell-like pattern
  -regex string
        regexp pattern
```
Where:

* glob: pattern where `/*.c` get all files in any depth with cuffix ".c"
* regex: POSIX regexp for names
* patt: Unix shell-like pattern with '/' as directory separator

The filter combinate with logic AND


```
$ fstorecli -aLogin user -aPass user fileStats -glob '/*.c'
{
  "error": false,
  "result": {
    "count": 11680,
    "usage": 101218570
  }
}

$ fstorecli -aLogin user -aPass user  fileStats -regex '.tar.gz$'
{
  "error": false,
  "result": {
    "count": 32,
    "usage": 228998036
  }
}

```

#### Downloading file

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

#### Deleting file

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

##### User listing

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
#### Creating a user
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

#### Deleting a user itself

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
