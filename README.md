
## Overview
## Build
### Build for sqlite
If you want to use sqlite as your store database, you should make it supporting.
go complier option `tag` can finish the work:
```bash
go build -tag sqlite github.com/zktnotify/zktnotify -o zktnotify
```

## Configuration
### Database
zktnotify support sqlite3 and mysql. But sqlite3 need a build tag, see chapter Build.
In config.json, .xserver.database object.

|field|type|value|comment|
|----|----|----|----|
|type|string|`sqlite3`(default),mysql|database type supported|
|host|string||database host name or ip, not support unix|
|user|string||database user name|
|password|string||database user password|
|db_name|string||database instance|
|path|string||path of database file, only used for sqlite|
