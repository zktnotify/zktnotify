
## Overview
## Build
### Build for sqlite
If you want to use sqlite as your store database, you should make it supporting.
go complier option `tag` can finish the work:
```bash
go build -tags sqlite3 -o zktnotify
```

## Configuration
### Database
zktnotify support sqlite3 and mysql. But sqlite3 need a build tag, see chapter Build.
In config.json, .xserver.database object.

|field|type|value|comment|
|----|----|----|----|
|type|string|`sqlite3`(default),mysql|database type supported|
|host|string||database host name or ip, not support unix|
|port|uint32|default:3306|database port|
|user|string||database user name|
|password|string||database user password|
|db_name|string||database instance|
|path|string||path of database file, only used for sqlite|

We can use enviroment variable to replace the Configuration of DB.

|variable|config key|comment|
|----|----|----|
|ZKTNOTIFY_DB_TYPE|type|||
|ZKTNOTIFY_DB_HOST|host|||
|ZKTNOTIFY_DB_PORT|port|||
|ZKTNOTIFY_DB_USER|user|||
|ZKTNOTIFY_DB_NAME|db_name|||
|ZKTNOTIFY_DB_PASSWORD|password|||
