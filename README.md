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
6060                                    PPROF and prometheus(if METRIC is prometheus)
7000                                    RAW TCP OR WEBSOCKET
8000                                    WEB
9000                                    CRPC
10000                                   GRPC
```

## 环境变量
```
PROJECT                                 该项目所属的项目,[a-z][0-9],第一个字符必须是[a-z]
GROUP                                   该项目所属的项目下的小组,[a-z][0-9],第一个字符必须是[a-z]
RUN_ENV                                 当前运行环境,如:test,pre,prod
DEPLOY_ENV                              部署环境,如:ali-kube-shanghai-1,ali-host-hangzhou-1
TRACE                                   是否开启链路追踪,空-不启用,不空-trace输出方式,[log,otlphttp,otlpgrpc,zipkin]
ZIPKIN_URL                              当TRACE为zipkin时,该变量为zipkin服务器的推送url
OTEL_EXPORTER_OTLP_TRACES_ENDPOINT      当TRACE为otlphttp或otlpgrpc时,该变量为otlp服务器的推送url
METRIC                                  是否开启系统监控采集,空-不启用,不空-metric输出方式,[log,otlphttp,otlpgrpc,prometheus]
OTEL_EXPORTER_OTLP_METRICS_ENDPOINT     当METRIC为otlphttp或otlpgrpc时,该变量为otlp服务器的推送url
OTEL_EXPORTER_OTLP_ENDPOINT             二合一,可取代OTEL_EXPORTER_OTLP_TRACES_ENDPOINT和OTEL_EXPORTER_OTLP_METRICS_ENDPOINT,但优先级比前两者低

CONFIG_TYPE                             配置类型:0-使用本地配置.1-使用admin服务的远程配置中心功能
REMOTE_CONFIG_SECRET                    当CONFIG_TYPE为1时,admin服务中,该服务使用的配置加密密钥,最长31个字符
ADMIN_SERVICE_PROJECT                   当使用admin服务的远程配置中心,服务发现,权限管理功能时,需要设置该环境变量,该变量为admin服务所属的项目,[a-z][0-9],第一个字符必须是[a-z]
ADMIN_SERVICE_GROUP                     当使用admin服务的远程配置中心,服务发现,权限管理功能时,需要设置该环境变量,该变量为admin服务所属的项目下的小组,[a-z][0-9],第一个字符必须是[a-z]
ADMIN_SERVICE_WEB_HOST                  当使用admin服务的远程配置中心,服务发现,权限管理功能时,需要设置该环境变量,该变量为admin服务的host,不带scheme(tls取决于NewSdk时是否传入tls.Config)
ADMIN_SERVICE_WEB_PORT                  当使用admin服务的远程配置中心,服务发现,权限管理功能时,需要设置该环境变量,该变量为admin服务的web端口,默认为80/443(取决于NewSdk时是否使用tls)
ADMIN_SERVICE_CONFIG_ACCESS_KEY         当使用admin服务的远程配置中心功能时,admin服务的授权码
ADMIN_SERVICE_DISCOVER_ACCESS_KEY       当使用admin服务的服务发现功能时,admin服务的授权码
ADMIN_SERVICE_PERMISSION_ACCESS_KEY     当使用admin服务的权限控制功能时,admin服务的授权码
```

## 配置文件
```
AppConfig.json该文件配置了该服务需要使用的业务配置,可热更新
SourceConfig.json该文件配置了该服务需要使用的资源配置,不热更新
```

## Cache
### Redis(Version >= 7.0)
#### account_redis(Cluster mode is better)

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
    tel:"",
    email:"",
    oauths:{
        "service_name_1":"id in this service",
        "service_name_2":"id in this service",
    },
    money:{
        "cny":100,
        "usd":100,
    },
    btime:123,//timestamp,unit nanoseconds,>0 means this account is banned,==0 means not banned
    breason:"",//ban reason
}
//手动创建数据库
use account;
db.createCollection("user");
sh.shardCollection("account.user",{_id:"hashed"});

collection: user_oauth_index
{
    service:"",//service_name+'|'+id
    user_id:ObjectId("xxx"),//collection user's _id field
}
//手动创建数据库
use account;
db.createCollection("user_oauth_index");
db.user_oauth_index.createIndex({service:1},{unique:true});
sh.shardCollection("account.user_oauth_index",{service:"hashed"});

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
db.user_idcard_index.createIndex({idcard:1},{unique:true});
sh.shardCollection("account.user_idcard_index",{idcard:"hashed"});

collection money_log
{
    _id:ObjectId("xxx"),//log id
    user_id:ObjectId("xxx"),//collection user's _id field
    action:"",//spend,recharge,refund
    unique_id:"",
    src_dst:"",
    money_type:"",
    money_amount:10,
}
//手动创建数据库
use account;
db.createCollection("money_log");
db.money_log.createIndex({user_id:1,action:1,unique_id:1},{unique:true});
db.money_log.createIndex({action:1,src_dst:1});
sh.shardCollection("account.money_log",{user_id:"hashed"});
```
