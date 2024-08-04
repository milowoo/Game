# game frame

## 简介
这个是一款针对小游戏/语音房的框架代码，采用 GO 语言实现, 采用多种通用的插件技术,包括 redis mongo, nats nacos等。同时为方便开发调试，提供了功能测试工具&压力测试工具。


#### 特性

- 可适应开发、测试、生成环境的配置；
- 统一的输出格式
- 整合mgo三方库的连接池及简易调用方法；
- 整合redis go三方库的连接池及简易调用方法、管道调用方法；
- 整合nats 消息中间件
- 整合 nacos配置以及服务管理


#### 目录结构

    game
    |_ucenter 用户管理服务
       |_conf 存放配置文件
       |_src  代码目录
          |_config 配置文件处理
          |_constants 服务通用宏定义
          |_domain 服务通用struct定义
          |_handler 业务处理
          |_internal 服务全局变量定义
          |_log 日志处理方法
          |_main 服务入口
          |_mongo mongo db 封装接口
          |_mq  nats消息中间件
          |_pb  proto 文件定义
          |_redis redis db 封装接口
          |_utils 服务通用接口封装  
    |_game_mgr 游戏管理服务 （包括服务管理，业务处理）
       |_conf 存放配置文件
       |_src  代码目录
          |_config 配置文件处理
          |_constants 服务通用宏定义
          |_domain 服务通用struct定义
          |_handler 业务处理
          |_internal 服务全局变量定义
          |_log 日志处理方法
          |_main 服务入口
          |_mongo mongo db 封装接口
          |_mq  nats消息中间件
          |_pb  proto 文件定义
          |_redis redis db 封装接口
          |_utils 服务通用接口封装
    |_match 匹配服务 （包括服务管理，业务处理）
       |_conf 存放配置文件
       |_src  代码目录
          |_config 配置文件处理
          |_constants 服务通用宏定义
          |_domain 服务通用struct定义
          |_handler 业务处理
          |_internal 服务全局变量定义
          |_log 日志处理方法
          |_main 服务入口
          |_mq  nats消息中间件
          |_pb  proto 文件定义
          |_redis redis db 封装接口
          |_utils 服务通用接口封装戏
    |_gateway 网关服务
       |_conf 存放配置文件
       |_src  代码目录（包括服务管理，业务处理）
          |_config 配置文件处理
          |_constants 服务通用宏定义
          |_domain 服务通用struct定义
          |_handler 业务处理
          |_internal 服务全局变量定义
          |_log 日志处理方法
          |_main 服务入口
          |_mq  nats消息中间件
          |_pb  proto 文件定义
          |_redis redis db 封装接口
          |_utils 服务通用接口封装
        
    |_game_frame 统一的小游戏框架
       |_conf 存放配置文件
       |_src  代码目录（包括服务管理，业务处理）
          |_config 配置文件处理
          |_constants 服务通用宏定义
          |_domain 服务通用struct定义
          |_handler 业务处理
          |_internal 服务全局变量定义
          |_log 日志处理方法
          |_main 服务入口
          |_mq  nats消息中间件
          |_pb  proto 文件定义
          |_redis redis db 封装接口
          |_utils 服务通用接口封装


#### 用例
##### 方式一：download
- 下载项目:

    `git clone https://github.com/milowoo/Game {your project name}`

- 由于使用go module,请自定义go.mod文件的replace本地代码目录

- 在/config目录下重置或新增配置项，并解析到全局变量

- 运行程序：

    `go run main.go`


    



