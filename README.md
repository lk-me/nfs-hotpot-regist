# nfs-hotpot-regist
    `nfs-hotpot-regist`目前为dbus框架，用于后期存储nfs系统上的hotpot信息，使用go开发。

### 功能
- 单例
- 调用者校验，包括权限与调用者名称
- 提供存储与读取接口
- 将信息存储值属性中，实现秒读取
- 提供日志模块，将日志存储至系统日志中

### /debian/ubuntu
- apt install libdbus-1-dev
- 在主目录使用makefile编译
- sudo ./nfs-hotpot-regist

### 模块介绍
- `/cmd`主程序目录
- `/config`存放配置文件
- `/test`测试程序，目前用c开发测试dbus接口
- `/pkg/logger`日志模块
- `/pkg/hotpot`主程序逻辑模块
- `/pkg/module`主程序调用模块

### todo
- 接口功能待实现
- 日志模块后期修改为输出至文件
- 单元测试