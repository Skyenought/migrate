## hertz_migrate

### 介绍

基于 go/ast 将目标代码迁移到 cloudwego/hertz.

### 支持进行迁移的框架

- net/http
    部分支持
- gin
    基本为完美支持, 覆盖大范围的常用 api 迁移
- chi 
    部分支持, `func (http.ResponseWriter, *http.Request)` 可以转为 cloudwego/hertz 的 `func(ctx context.Context, c *app.RequestContext)`)
  
### 安装

```bash
go install github.com/hertz-contrib/migrate/cmd/hertz_migrate@latest
```

### 命令行参数

#### --hz-repo 
  别名 `-r`, 用于设定 hertz 的仓库地址, 默认为 `github.com/cloudwego/hertz`
#### --target-dir     
  别名 `-d`, 设定需要迁移的代码目标文件夹
#### --ignore-dir 
  别名 `-I`, 可以设定需要迁移工具忽略的文件夹, 可以声明多个, 适当的文件夹忽略可以提高工具性能.

  例如:
  ```bash
  hertz_migrate -target-dir ./project -ignore-dirs=hz_gen -ignore-dirs=vendor
  ```
    
#### --use-gin
  别名 `-g`, 使迁移工具启用 gin 的迁移程序

#### --use-net-http
  别名 `-n`, 使迁移工具启用 net/http 的迁移程序

#### --use-chi
  别名 `-c`, 使迁移工具启用 chi 的迁移程序, 但注意指定该参数的同时会同时打开 `--use-gin` 参数
