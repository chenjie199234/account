// Code generated by protoc-gen-browser. DO NOT EDIT.
// version:
// 	protoc-gen-browser v0.0.115<br />
// 	protoc             v5.26.1<br />
// source: api/account_base.proto<br />

export interface LogicError{
	code: number;
	msg: string;
}

export class BaseInfo{
	user_id: string = ''
	idcard: string = ''
	tel: string = ''
	email: string = ''
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	ctime: number = 0
	bind_oauths: Array<string>|null = null
	//Warning!!!map's value's type is int32,be careful of sign(+,-) and overflow
	money: Map<string,number>|null = null
	ban: string = ''//if this is not empty,means this account is banned
	fromOBJ(obj:Object){
		if(obj["user_id"]){
			this.user_id=obj["user_id"]
		}
		if(obj["idcard"]){
			this.idcard=obj["idcard"]
		}
		if(obj["tel"]){
			this.tel=obj["tel"]
		}
		if(obj["email"]){
			this.email=obj["email"]
		}
		if(obj["ctime"]){
			this.ctime=obj["ctime"]
		}
		if(obj["bind_oauths"] && obj["bind_oauths"].length>0){
			this.bind_oauths=obj["bind_oauths"]
		}
		if(obj["money"] && Object.keys(obj["money"]).length>0){
			this.money=new Map<string,number>()
			for(let key of Object.keys(obj["money"])){
				this.money.set(key,obj["money"][key])
			}
		}
		if(obj["ban"]){
			this.ban=obj["ban"]
		}
	}
}
export class BaseInfoReq{
	src_type: string = ''
	src: string = ''//if this is empty,means get self's baseinfo,src_type will force to user_id and the src is from token
	toJSON(){
		let tmp = {}
		if(this.src_type){
			tmp["src_type"]=this.src_type
		}
		if(this.src){
			tmp["src"]=this.src
		}
		return tmp
	}
}
export class BaseInfoResp{
	info: BaseInfo|null = null
	fromOBJ(obj:Object){
		if(obj["info"]){
			this.info=new BaseInfo()
			this.info.fromOBJ(obj["info"])
		}
	}
}
export class DelEmailReq{
	verify_src_type: string = ''
	//when verify_src_type is oauth,this is the oauth service name
	verify_src_type_extra: string = ''
	//if this is empty,means send dynamic password
	//if this is not empty,means verify dynamic password
	verify_dynamic_password: string = ''
	toJSON(){
		let tmp = {}
		if(this.verify_src_type){
			tmp["verify_src_type"]=this.verify_src_type
		}
		if(this.verify_src_type_extra){
			tmp["verify_src_type_extra"]=this.verify_src_type_extra
		}
		if(this.verify_dynamic_password){
			tmp["verify_dynamic_password"]=this.verify_dynamic_password
		}
		return tmp
	}
}
export class DelEmailResp{
	//oldverify:server already send the dynamic password to user's email or tel(depend on the del_email_req's verify_src_type) and is waiting for verify
	//success:nothing need to do
	step: string = ''
	//if this is true,means this is the last way to login this account
	//if del this,this account will be deleted completely
	final: boolean = false
	//send dynamic password to where,this will be masked
	//when step is success,ignore this
	receiver: string = ''
	fromOBJ(obj:Object){
		if(obj["step"]){
			this.step=obj["step"]
		}
		if(obj["final"]){
			this.final=obj["final"]
		}
		if(obj["receiver"]){
			this.receiver=obj["receiver"]
		}
	}
}
export class DelIdcardReq{
	verify_src_type: string = ''
	//when verify_src_type is oauth,this is the oauth service name
	verify_src_type_extra: string = ''
	//if this is empty,means send dynamic password
	//if this is not empty,means verify dynamic password
	verify_dynamic_password: string = ''
	toJSON(){
		let tmp = {}
		if(this.verify_src_type){
			tmp["verify_src_type"]=this.verify_src_type
		}
		if(this.verify_src_type_extra){
			tmp["verify_src_type_extra"]=this.verify_src_type_extra
		}
		if(this.verify_dynamic_password){
			tmp["verify_dynamic_password"]=this.verify_dynamic_password
		}
		return tmp
	}
}
export class DelIdcardResp{
	//oldverify:server already send the dynamic password to user's email or tel(depend on the del_idcard_req's verify_src_type) and is waiting for verify
	//success:nothing need to do
	step: string = ''
	//if this is true,means this is the last way to login this account
	//if del this,this account will be deleted completely
	final: boolean = false
	//send dynamic password to where,this will be masked
	//when step is success,ignore this
	receiver: string = ''
	fromOBJ(obj:Object){
		if(obj["step"]){
			this.step=obj["step"]
		}
		if(obj["final"]){
			this.final=obj["final"]
		}
		if(obj["receiver"]){
			this.receiver=obj["receiver"]
		}
	}
}
export class DelOauthReq{
	verify_src_type: string = ''
	//when verify_src_type is oauth,this is the oauth service name
	verify_src_type_extra: string = ''
	//if this is empty,means send dynamic password
	//if this is not empty,means verify dynamic password
	verify_dynamic_password: string = ''
	del_oauth_service_name: string = ''
	toJSON(){
		let tmp = {}
		if(this.verify_src_type){
			tmp["verify_src_type"]=this.verify_src_type
		}
		if(this.verify_src_type_extra){
			tmp["verify_src_type_extra"]=this.verify_src_type_extra
		}
		if(this.verify_dynamic_password){
			tmp["verify_dynamic_password"]=this.verify_dynamic_password
		}
		if(this.del_oauth_service_name){
			tmp["del_oauth_service_name"]=this.del_oauth_service_name
		}
		return tmp
	}
}
export class DelOauthResp{
	//oldverify:server already send the dynamic password to user's email or tel(depend on the del_oauth_req's verify_src_type) and is waiting for verify
	//success:nothing need to do
	step: string = ''
	//if this is true,means this is the last way to login this account
	//if del this,this account will be deleted completely
	final: boolean = false
	//send dynamic password to where,this will be masked
	//when step is success,ignore this
	receiver: string = ''
	fromOBJ(obj:Object){
		if(obj["step"]){
			this.step=obj["step"]
		}
		if(obj["final"]){
			this.final=obj["final"]
		}
		if(obj["receiver"]){
			this.receiver=obj["receiver"]
		}
	}
}
export class DelTelReq{
	verify_src_type: string = ''
	//when verify_src_type is oauth,this is the oauth service name
	verify_src_type_extra: string = ''
	//if this is empty,means send dynamic password
	//if this is not empty,means verify dynamic password
	verify_dynamic_password: string = ''
	toJSON(){
		let tmp = {}
		if(this.verify_src_type){
			tmp["verify_src_type"]=this.verify_src_type
		}
		if(this.verify_src_type_extra){
			tmp["verify_src_type_extra"]=this.verify_src_type_extra
		}
		if(this.verify_dynamic_password){
			tmp["verify_dynamic_password"]=this.verify_dynamic_password
		}
		return tmp
	}
}
export class DelTelResp{
	//oldverify:server already send the dynamic password to user's email or tel(depend on the del_tel_req's verify_src_type) and is waiting for verify
	//success:nothing need to do
	step: string = ''
	//if this is true,means this is the last way to login this account
	//if del this,this account will be deleted completely
	final: boolean = false
	//send dynamic password to where,this will be masked
	//when step is success,ignore this
	receiver: string = ''
	fromOBJ(obj:Object){
		if(obj["step"]){
			this.step=obj["step"]
		}
		if(obj["final"]){
			this.final=obj["final"]
		}
		if(obj["receiver"]){
			this.receiver=obj["receiver"]
		}
	}
}
export class EmailDuplicateCheckReq{
	email: string = ''
	toJSON(){
		let tmp = {}
		if(this.email){
			tmp["email"]=this.email
		}
		return tmp
	}
}
export class EmailDuplicateCheckResp{
	duplicate: boolean = false
	fromOBJ(obj:Object){
		if(obj["duplicate"]){
			this.duplicate=obj["duplicate"]
		}
	}
}
export class GetOauthUrlReq{
	oauth_service_name: string = ''
	toJSON(){
		let tmp = {}
		if(this.oauth_service_name){
			tmp["oauth_service_name"]=this.oauth_service_name
		}
		return tmp
	}
}
export class GetOauthUrlResp{
	url: string = ''
	fromOBJ(obj:Object){
		if(obj["url"]){
			this.url=obj["url"]
		}
	}
}
export class IdcardDuplicateCheckReq{
	idcard: string = ''
	toJSON(){
		let tmp = {}
		if(this.idcard){
			tmp["idcard"]=this.idcard
		}
		return tmp
	}
}
export class IdcardDuplicateCheckResp{
	duplicate: boolean = false
	fromOBJ(obj:Object){
		if(obj["duplicate"]){
			this.duplicate=obj["duplicate"]
		}
	}
}
export class LoginReq{
	src_type: string = ''
	//when src_type is oauth,this is the oauth service name
	src_type_extra: string = ''
	//when src_type is idcard this can't be dynamic
	//when src_type is oauth,this can't be static
	password_type: string = ''
	//when password_type is dynamic and this is empty,means send dynamic password to email or tel.
	//when password_type is dynamic and this is not empty,means verify dynamic password.
	password: string = ''
	toJSON(){
		let tmp = {}
		if(this.src_type){
			tmp["src_type"]=this.src_type
		}
		if(this.src_type_extra){
			tmp["src_type_extra"]=this.src_type_extra
		}
		if(this.password_type){
			tmp["password_type"]=this.password_type
		}
		if(this.password){
			tmp["password"]=this.password
		}
		return tmp
	}
}
export class LoginResp{
	token: string = ''
	info: BaseInfo|null = null
	//verify:server already send the dynamic password to user's email or tel(depend on the login_req's src_type and src) and is waiting for verify
	//password:login success,but this account is new and it can be setted with a static password(optional)
	//success:nothing need to do
	step: string = ''
	fromOBJ(obj:Object){
		if(obj["token"]){
			this.token=obj["token"]
		}
		if(obj["info"]){
			this.info=new BaseInfo()
			this.info.fromOBJ(obj["info"])
		}
		if(obj["step"]){
			this.step=obj["step"]
		}
	}
}
export class TelDuplicateCheckReq{
	tel: string = ''
	toJSON(){
		let tmp = {}
		if(this.tel){
			tmp["tel"]=this.tel
		}
		return tmp
	}
}
export class TelDuplicateCheckResp{
	duplicate: boolean = false
	fromOBJ(obj:Object){
		if(obj["duplicate"]){
			this.duplicate=obj["duplicate"]
		}
	}
}
export class UpdateEmailReq{
	verify_src_type: string = ''
	//when verify_src_type is oauth,this is the oauth service name
	verify_src_type_extra: string = ''
	//if this is empty,means send dynamic password
	//if this is not empty,means verify dynamic password
	verify_dynamic_password: string = ''
	new_email: string = ''
	//if this is empty,means send dynamic password.
	//if this is not empty,means verify dynamic password.
	new_email_dynamic_password: string = ''
	toJSON(){
		let tmp = {}
		if(this.verify_src_type){
			tmp["verify_src_type"]=this.verify_src_type
		}
		if(this.verify_src_type_extra){
			tmp["verify_src_type_extra"]=this.verify_src_type_extra
		}
		if(this.verify_dynamic_password){
			tmp["verify_dynamic_password"]=this.verify_dynamic_password
		}
		if(this.new_email){
			tmp["new_email"]=this.new_email
		}
		if(this.new_email_dynamic_password){
			tmp["new_email_dynamic_password"]=this.new_email_dynamic_password
		}
		return tmp
	}
}
export class UpdateEmailResp{
	//oldverify:server already send the dynamic password to user's email or tel(depend on the update_email_req's verify_src_type) and is waiting for verify
	//newverify:server already send the dynamic password to the new email(depend on the update_email_req's new_email) and is waiting for verify
	//success:nothing need to do
	step: string = ''
	//send dynamic password to where,this will be masked
	//when step is success,ignore this
	receiver: string = ''
	fromOBJ(obj:Object){
		if(obj["step"]){
			this.step=obj["step"]
		}
		if(obj["receiver"]){
			this.receiver=obj["receiver"]
		}
	}
}
export class UpdateIdcardReq{
	verify_src_type: string = ''
	//when verify_src_type is oauth,this is the oauth service name
	verify_src_type_extra: string = ''
	//if this is empty,means send dynamic password
	//if this is not empty,means verify dynamic password
	verify_dynamic_password: string = ''
	new_idcard: string = ''
	toJSON(){
		let tmp = {}
		if(this.verify_src_type){
			tmp["verify_src_type"]=this.verify_src_type
		}
		if(this.verify_src_type_extra){
			tmp["verify_src_type_extra"]=this.verify_src_type_extra
		}
		if(this.verify_dynamic_password){
			tmp["verify_dynamic_password"]=this.verify_dynamic_password
		}
		if(this.new_idcard){
			tmp["new_idcard"]=this.new_idcard
		}
		return tmp
	}
}
export class UpdateIdcardResp{
	//oldverify:server already send the dynamic password to user's email or tel(depend on the update_idcard_req's verify_src_type) and is waiting for verify
	//success:nothing need to do
	step: string = ''
	//send dynamic password to where,this will be masked
	//when step is success,ignore this
	receiver: string = ''
	fromOBJ(obj:Object){
		if(obj["step"]){
			this.step=obj["step"]
		}
		if(obj["receiver"]){
			this.receiver=obj["receiver"]
		}
	}
}
export class UpdateOauthReq{
	verify_src_type: string = ''
	//when verify_src_type is oauth,this is the oauth service name
	verify_src_type_extra: string = ''
	//if this is empty,means send dynamic password
	//if this is not empty,means verify dynamic password
	verify_dynamic_password: string = ''
	//if verify_dynamic_password is not empty,this should not be empty too
	new_oauth_service_name: string = ''
	//if verify_dynamic_password is not empty,this should not be empty too
	new_oauth_dynamic_password: string = ''
	toJSON(){
		let tmp = {}
		if(this.verify_src_type){
			tmp["verify_src_type"]=this.verify_src_type
		}
		if(this.verify_src_type_extra){
			tmp["verify_src_type_extra"]=this.verify_src_type_extra
		}
		if(this.verify_dynamic_password){
			tmp["verify_dynamic_password"]=this.verify_dynamic_password
		}
		if(this.new_oauth_service_name){
			tmp["new_oauth_service_name"]=this.new_oauth_service_name
		}
		if(this.new_oauth_dynamic_password){
			tmp["new_oauth_dynamic_password"]=this.new_oauth_dynamic_password
		}
		return tmp
	}
}
export class UpdateOauthResp{
	//oldverify:server already send the dynamic password to user's email or tel(depend on the update_oauth_req's verify_src_type) and is waiting for verify
	//success:nothing need to do
	step: string = ''
	//send dynamic password to where,this will be masked
	//when step is success,ignore this
	receiver: string = ''
	fromOBJ(obj:Object){
		if(obj["step"]){
			this.step=obj["step"]
		}
		if(obj["receiver"]){
			this.receiver=obj["receiver"]
		}
	}
}
export class UpdateStaticPasswordReq{
	//if this is empty,means this is the first time to set the static password
	old_static_password: string = ''
	new_static_password: string = ''
	toJSON(){
		let tmp = {}
		if(this.old_static_password){
			tmp["old_static_password"]=this.old_static_password
		}
		if(this.new_static_password){
			tmp["new_static_password"]=this.new_static_password
		}
		return tmp
	}
}
export class UpdateStaticPasswordResp{
	fromOBJ(_obj:Object){
	}
}
export class UpdateTelReq{
	verify_src_type: string = ''
	//when verify_src_type is oauth,this is the oauth service name
	verify_src_type_extra: string = ''
	//if this is empty,means send dynamic password
	//if this is not empty,means verify dynamic password
	verify_dynamic_password: string = ''
	new_tel: string = ''
	//if this is empty,means send dynamic password.
	//if this is not empty,means verify dynamic password.
	new_tel_dynamic_password: string = ''
	toJSON(){
		let tmp = {}
		if(this.verify_src_type){
			tmp["verify_src_type"]=this.verify_src_type
		}
		if(this.verify_src_type_extra){
			tmp["verify_src_type_extra"]=this.verify_src_type_extra
		}
		if(this.verify_dynamic_password){
			tmp["verify_dynamic_password"]=this.verify_dynamic_password
		}
		if(this.new_tel){
			tmp["new_tel"]=this.new_tel
		}
		if(this.new_tel_dynamic_password){
			tmp["new_tel_dynamic_password"]=this.new_tel_dynamic_password
		}
		return tmp
	}
}
export class UpdateTelResp{
	//oldverify:server already send the dynamic password to user's email or tel(depend on the update_tel_req's verify_src_type) and is waiting for verify
	//newverify:server already send the dynamic password to the new tel(depend on the update_tel_req's new_tel) and is waiting for verify
	//success:nothing need to do
	step: string = ''
	//send dynamic password to where,this will be masked
	//when step is success,ignore this
	receiver: string = ''
	fromOBJ(obj:Object){
		if(obj["step"]){
			this.step=obj["step"]
		}
		if(obj["receiver"]){
			this.receiver=obj["receiver"]
		}
	}
}
//timeout's unit is millisecond,it will be used when > 0
function call(timeout: number,url: string,opts: Object,error: (arg: LogicError)=>void,success: (arg: Object)=>void){
	let tid: number|null = null
	if(timeout>0){
		const c = new AbortController()
		opts["signal"] = c.signal
		tid = setTimeout(()=>{c.abort()},timeout)
	}
	let ok=false
	fetch(url,opts)
	.then(r=>{
		ok=r.ok
		if(r.ok){
			return r.json()
		}
		return r.text()
	})
	.then(d=>{
		if(!ok){
			throw d
		}
		success(d)
	})
	.catch(e=>{
		if(e instanceof Error){
			error({code:-1,msg:e.message})
		}else if(e.length>0 && e[0]=='{' && e[e.length-1]=='}'){
			error(JSON.parse(e))
		}else{
			error({code:-1,msg:e})
		}
	})
	.finally(()=>{
		if(tid){
			clearTimeout(tid)
		}
	})
}
const _WebPathBaseGetOauthUrl: string ="/account.base/get_oauth_url";
const _WebPathBaseLogin: string ="/account.base/login";
const _WebPathBaseBaseInfo: string ="/account.base/base_info";
const _WebPathBaseUpdateStaticPassword: string ="/account.base/update_static_password";
const _WebPathBaseUpdateOauth: string ="/account.base/update_oauth";
const _WebPathBaseDelOauth: string ="/account.base/del_oauth";
const _WebPathBaseIdcardDuplicateCheck: string ="/account.base/idcard_duplicate_check";
const _WebPathBaseUpdateIdcard: string ="/account.base/update_idcard";
const _WebPathBaseDelIdcard: string ="/account.base/del_idcard";
const _WebPathBaseEmailDuplicateCheck: string ="/account.base/email_duplicate_check";
const _WebPathBaseUpdateEmail: string ="/account.base/update_email";
const _WebPathBaseDelEmail: string ="/account.base/del_email";
const _WebPathBaseTelDuplicateCheck: string ="/account.base/tel_duplicate_check";
const _WebPathBaseUpdateTel: string ="/account.base/update_tel";
const _WebPathBaseDelTel: string ="/account.base/del_tel";
//ToC means this is for users
export class BaseBrowserClientToC {
	constructor(host: string){
		if(!host || host.length==0){
			throw "BaseBrowserClientToC's host missing"
		}
		this.host=host
	}
	//timeout's unit is millisecond,it will be used when > 0
	get_oauth_url(header: Object,req: GetOauthUrlReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: GetOauthUrlResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathBaseGetOauthUrl,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new GetOauthUrlResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	login(header: Object,req: LoginReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: LoginResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathBaseLogin,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new LoginResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	base_info(header: Object,req: BaseInfoReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: BaseInfoResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathBaseBaseInfo,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new BaseInfoResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	update_static_password(header: Object,req: UpdateStaticPasswordReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: UpdateStaticPasswordResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathBaseUpdateStaticPassword,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new UpdateStaticPasswordResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	update_oauth(header: Object,req: UpdateOauthReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: UpdateOauthResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathBaseUpdateOauth,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new UpdateOauthResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	del_oauth(header: Object,req: DelOauthReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: DelOauthResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathBaseDelOauth,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new DelOauthResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	idcard_duplicate_check(header: Object,req: IdcardDuplicateCheckReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: IdcardDuplicateCheckResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathBaseIdcardDuplicateCheck,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new IdcardDuplicateCheckResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	update_idcard(header: Object,req: UpdateIdcardReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: UpdateIdcardResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathBaseUpdateIdcard,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new UpdateIdcardResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	del_idcard(header: Object,req: DelIdcardReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: DelIdcardResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathBaseDelIdcard,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new DelIdcardResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	email_duplicate_check(header: Object,req: EmailDuplicateCheckReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: EmailDuplicateCheckResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathBaseEmailDuplicateCheck,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new EmailDuplicateCheckResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	update_email(header: Object,req: UpdateEmailReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: UpdateEmailResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathBaseUpdateEmail,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new UpdateEmailResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	del_email(header: Object,req: DelEmailReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: DelEmailResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathBaseDelEmail,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new DelEmailResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	tel_duplicate_check(header: Object,req: TelDuplicateCheckReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: TelDuplicateCheckResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathBaseTelDuplicateCheck,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new TelDuplicateCheckResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	update_tel(header: Object,req: UpdateTelReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: UpdateTelResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathBaseUpdateTel,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new UpdateTelResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	del_tel(header: Object,req: DelTelReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: DelTelResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathBaseDelTel,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new DelTelResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	private host: string
}
