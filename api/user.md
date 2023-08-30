// Code generated by protoc-gen-markdown. DO NOT EDIT.<br />
// version:<br />
// 	protoc-gen-markdown v0.0.77<br />
// 	protoc              v4.24.1<br />
// source: api/user.proto<br />

## user
### login

#### Req:
```
Path:         /account.user/login
Method:       POST
Content-Type: application/json
------------------------------------------------------------------------------------------------------------
{
	//value must in ["tel","email","idcard","nickname"]
	"src_type":"str",
	//value length must > 0
	"src":"str",
	//when src_type is idcard or nickname,this can't be dynamic
	//value must in ["static","dynamic"]
	"password_type":"str",
	//when password_type is static this length must >=10
	//when password_type is dynamic and this is empty,means send dynamic password to email or tel.
	//when password_type is dynamic and this is not empty,means verify dynamic password.
	"password":"str"
}
------------------------------------------------------------------------------------------------------------
```
#### Resp:
```
Fail:    httpcode:4xx/5xx
------------------------------------------------------------------------------------------------------------
{"code":123,"msg":"error message"}
------------------------------------------------------------------------------------------------------------
Success: httpcode:200
------------------------------------------------------------------------------------------------------------
{
	"token":"str",
	//object user_info
	"info":{}
}
------------------------------------------------------------------------------------------------------------
user_info: {
	"user_id":"str",
	"idcard":"str",
	"tel":"str",
	"email":"str",
	"nick_name":"str",
	//int64 use string to avoid overflow
	"ctime":"0",
	//kv map,value-int64 use string to avoid overflow
	"money":{"str":"0","str":"0"}
}
------------------------------------------------------------------------------------------------------------
```
### update_static_password

#### Req:
```
Path:         /account.user/update_static_password
Method:       POST
Content-Type: application/json
------------------------------------------------------------------------------------------------------------
{
	//if this is empty,means this is the first time to set the static password
	"old_static_password":"str",
	//value length must >= 10
	"new_static_password":"str"
}
------------------------------------------------------------------------------------------------------------
```
#### Resp:
```
Fail:    httpcode:4xx/5xx
------------------------------------------------------------------------------------------------------------
{"code":123,"msg":"error message"}
------------------------------------------------------------------------------------------------------------
Success: httpcode:200
------------------------------------------------------------------------------------------------------------
{
}
------------------------------------------------------------------------------------------------------------
```
### update_nick_name

#### Req:
```
Path:         /account.user/update_nick_name
Method:       POST
Content-Type: application/json
------------------------------------------------------------------------------------------------------------
{
	//value length must > 0
	"new_nick_name":"str"
}
------------------------------------------------------------------------------------------------------------
```
#### Resp:
```
Fail:    httpcode:4xx/5xx
------------------------------------------------------------------------------------------------------------
{"code":123,"msg":"error message"}
------------------------------------------------------------------------------------------------------------
Success: httpcode:200
------------------------------------------------------------------------------------------------------------
{
}
------------------------------------------------------------------------------------------------------------
```
### update_email

#### Req:
```
Path:         /account.user/update_email
Method:       POST
Content-Type: application/json
------------------------------------------------------------------------------------------------------------
{
	//value must in ["email","tel"]
	"old_receiver_type":"str",
	//if this is empty,means send dynamic password
	//if this is not empty,means verify dynamic password
	"old_dynamic_password":"str",
	//value length must > 0
	"new_email":"str",
	//if this is empty,means send dynamic password.
	//if this is not empty,means verify dynamic password.
	"new_email_dynamic_password":"str"
}
------------------------------------------------------------------------------------------------------------
```
#### Resp:
```
Fail:    httpcode:4xx/5xx
------------------------------------------------------------------------------------------------------------
{"code":123,"msg":"error message"}
------------------------------------------------------------------------------------------------------------
Success: httpcode:200
------------------------------------------------------------------------------------------------------------
{
}
------------------------------------------------------------------------------------------------------------
```
### update_tel

#### Req:
```
Path:         /account.user/update_tel
Method:       POST
Content-Type: application/json
------------------------------------------------------------------------------------------------------------
{
	//value must in ["email","tel"]
	"old_receiver_type":"str",
	//if this is empty,means send dynamic password
	//if this is not empty,means verify dynamic password
	"old_dynamic_password":"str",
	//value length must > 0
	"new_tel":"str",
	//if this is empty,means send dynamic password.
	//if this is not empty,means verify dynamic password.
	"new_email_dynamic_password":"str"
}
------------------------------------------------------------------------------------------------------------
```
#### Resp:
```
Fail:    httpcode:4xx/5xx
------------------------------------------------------------------------------------------------------------
{"code":123,"msg":"error message"}
------------------------------------------------------------------------------------------------------------
Success: httpcode:200
------------------------------------------------------------------------------------------------------------
{
}
------------------------------------------------------------------------------------------------------------
```