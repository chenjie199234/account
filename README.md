# account
```
account是一个微服务.
运行cmd脚本可查看使用方法.windows下将./cmd.sh换为cmd.bat
./cmd.sh help 输出帮助信息
./cmd.sh pb 解析proto文件,生成桩代码
./cmd.sh sub 在该项目中创建一个新的子服务
./cmd.sh kube 新建kubernetes的配置
./cmd.sh html 新建前端html代码模版
```

## 服务端口
```
6060                                    MONITOR AND PPROF
8000                                    WEB
9000                                    CRPC
10000                                   GRPC
```

## 环境变量
```
LOG_LEVEL                               日志等级,debug,info(default),warning,error
LOG_TRACE                               是否开启链路追踪,1-开启,0-关闭(default)
LOG_TARGET                              日志输出目标,std-输出到标准输出,file-输出到文件(可执行文件相同目录),both-两者都输出
PROJECT                                 该项目所属的项目,[a-z][0-9],第一个字符必须是[a-z]
GROUP                                   该项目所属的项目组,[a-z][0-9],第一个字符必须是[a-z]
RUN_ENV                                 当前运行环境,如:test,pre,prod
DEPLOY_ENV                              部署环境,如:ali-kube-shanghai-1,ali-host-hangzhou-1
MONITOR                                 是否开启系统监控采集,0关闭,1开启

CONFIG_TYPE                             配置类型:0-使用本地配置.1-使用远程配置中心config服务
REMOTE_CONFIG_SERVICE_PROJECT           当CONFIG_TYPE为1时,远程配置中心服务的项目,[a-z][0-9],第一个字符必须是[a-z]
REMOTE_CONFIG_SERVICE_GROUP             当CONFIG_TYPE为1时,远程配置中心服务的项目组,[a-z][0-9],第一个字符必须是[a-z]
REMOTE_CONFIG_SERVICE_WEB_HOST          当CONFIG_TYPE为1时,远程配置中心服务的host,不带scheme(是否启用tls取决于NewConfigSdk时是否传入tls.Config)
REMOTE_CONFIG_SERVICE_WEB_PORT          当CONFIG_TYPE为1时,远程配置中心服务的端口,默认为80/443(取决于NewConfigSdk时是否使用tls)
REMOTE_CONFIG_SECRET                    当CONFIG_TYPE为1时,远程配置中心配置的密钥

PERMISSION_SERVICE_PROJECT              权限控制中心的项目(按需设置),[a-z][0-9],第一个字符必须是[a-z]
PERMISSION_SERVICE_GROUP                权限控制中心的项目组(按需设置),[a-z][0-9],第一个字符必须是[a-z]
PERMISSION_SERVICE_WEB_HOST             权限控制中心的host(按需设置),不带scheme(是否启用tls取决于NewConfigSdk时是否传入tls.Config)
PERMISSION_SERVICE_WEB_PORT             权限控制中心的端口(按需设置),默认为80/443(取决于NewPermissionSdk时是否使用tls)

```

## 配置文件
```
AppConfig.json该文件配置了该服务需要使用的业务配置,可热更新
SourceConfig.json该文件配置了该服务需要使用的资源配置,不热更新
```

## DB
### Mongo(Shard Mode)(Version >= 6.0)
#### Account
```
database: account

collection: user
{
    _id:ObjectId("xxx"),//user id
    password:"",
    idcard:"",//实名认证
    nick_name:"",
    tel:"",
    email:"",
    money:{
        "cny":100,
        "usd":100,
    },
}
//手动创建数据库
use account;
db.createCollection("user");
sh.shardCollection("account.user",{_id:"hashed"});

collection: user_email_index
{
    email:"",
    user_id:ObjectId("xxx"),//collection user's _id field
}
//手动创建数据库
use account;
db.createCollection("user_email_index");
db.user_email_index.createIndex({email:1},{unique:true});
sh.shardCollection("account.user_email_index",{email:"hashed"});

collection: user_tel_index
{
    tel:"",
    user_id:ObjectId("xxx"),//collection user's _id field
}
//手动创建数据库
use account;
db.createCollection("user_tel_index");
db.user_tel_index.createIndex({tel:1},{unique:true});
sh.shardCollection("account.user_tel_index",{tel:"hashed"});

collection: user_idcard_index
{
    idcard:"",//实名认证
    user_id:ObjectId("xxx"),//collection user's _id field
}
//手动创建数据库
use account;
db.createCollection("user_idcard_index");
db.user_tel_index.createIndex({idcard:1},{unique:true});
sh.shardCollection("account.user_idcard_index",{idcard:"hashed"});

collection: user_nick_name_index
{
    nick_name:"",
    user_id:ObjectId("xxx"),//collection user's _id field
}
//手动创建数据库
use account;
db.createCollection("user_nick_name_index");
db.user_nick_name_index.createIndex({nick_name:1},{unique:true});
sh.shardCollection("account.user_nick_name_index",{nick_name:"hashed"});

collection money_log
{
    _id:ObjectId("xxx"),//log id
    user_id:ObjectId("xxx"),//collection user's _id field
    action:"",//spend,recharge,refund
    unique_id:"",
    src_dst:"",
    money_type:"",
    money_amount:10,
    ext_data:"",
}
//手动创建数据库
use account;
db.createCollection("money_log");
db.money_log.createIndex({user_id:1,action:1,unique_id:1},{unique:true});
db.money_log.createIndex({action:1,src_dst:1});
sh.shardCollection("account.money_log",{user_id:"hashed"});
```
