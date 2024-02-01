### 替换 tpl 模板

下载并替换 `goctl` 模板

```shell
git clone https://github.com/tt90cc/goctl-template.git && rm -rf ~/.goctl/$(goctl -v|awk '{print $3}')/* && cd goctl-template && mv ./* ~/.goctl/$(goctl -v|awk '{print $3}')
```

### 根据 `DDL` 生成 `MODEL`

1. 修改 ddl `cd ./model && vim ./ddl.sql`
2. 在项目根目录执行 `goctl model mysql ddl -src ./ddl.sql -dir . -c`

##### 复杂查询

```go
squirrel.Or{squirrel.Expr("id=?", cast.ToInt64(req.Name)), squirrel.And{squirrel.Eq{"name": req.Name}}}
// squirrel.Or{squirrel.Eq{"id": cast.ToInt64(req.Name)}, squirrel.And{squirrel.Eq{"name": req.Name}}}

Where("FIND_IN_SET(?, platform_type)", req.PlatformType)
```

### 生成 `api` 或者 `rpc` 代码

1. 生成 api：`goctl api go -api ./main.api -dir .`
2. 生成 rpc：`goctl rpc protoc ./main.proto --go_out=./types --go-grpc_out=./types --zrpc_out=.`

### 配置

```yaml
Log:
  #  Mode: file
  #  Path: ../logs
  Encoding: plain
UcenterRpc:
  Timeout: 10000
  Endpoints:
    - 127.0.0.1:8213
```