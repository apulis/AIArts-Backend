# AIArtsBackend
Backend for AIArts.


## 开发说明


### 配置
* 配置文件为config.yaml，在项目根目录
* 配置结构体定义在`configs/config.go`，使用`viper`库读取配置
* `config.yaml`已放入.gitignore，格式按照`config.template.yaml`编写


### 数据库
* 使用mysql，配置写在`config.yaml`
* 开发时自行启动mysql
* 一个实例，定义在`database/db.go`，使用时需引入，无需再创建实例
* orm使用`gorm`库，[参考文档](https://gorm.io/zh_CN/docs/index.html)
* 本地开mysql，可使用local_mysql下的docker-compose


### 日志
* 日志配置写在了`config.yaml`
* 如果设置writefile为true，会写入文件并输出到console
* 如果设置writefile为false，日志只输出到console
* 一个实例，定义在`loggers/logger.go`，使用时需引入，无需再创建实例
* gin的日志已被`loggers/gin_logger.go`替换，保持与其他日志格式一致


### 路由
* 路由使用`gin`框架，整体路由在`routers/router.go`
* 各模块在单独的文件中实现
* 成功返回，可使用`routers/success_resp.go`中的函数
* 错误处理已经实现了部分函数`routers/error_handlers.go`
* 错误码统一定义在`routers/error_codes.go`(待讨论)


### API文档
* 已集成swagger，生成文档时执行`swag init`即可
* 写文档可参考每个router的sample

### 部署
* 公用数据集存放路径`/dlwsdata/storage/dataset/storage`
* 私有数据集存放路径`/dlwsdata/work/user/storage`
* 更新并push至harbor

  `sudo vim /etc/hosts`
  
   增加 10.31.3.211 harbor.sigsus.cn
 
  ```shell 
  sudo vim /etc/docker/daemon.json
  ```
    ```
  {
    "registry-mirrors": [],
    "insecure-registries": [
     "https://harbor.sigsus.cn:8443"
    ],
    "debug": true,
    "experimental": false
    }
  ```
* 重启docker并登陆(需要进入 https://10.31.3.211:8443 注册)
  ```shell 
  sudo systemctl  restart docker 
  sudo docker login harbor.sigsus.cn:8443
  ```
 * 重启docker并登陆(需要进入 https://10.31.3.211:8443 注册)
   ```shell 
   sudo systemctl  restart docker 
   sudo docker login harbor.sigsus.cn:8443
   ```
 * 开始推送到harbor
    ```shell 
    cd deployment
    ./build2harbor.sh
    ```



     
