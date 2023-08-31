// Code generated by protoc-gen-browser. DO NOT EDIT.
// version:
// 	protoc-gen-browser v0.0.77<br />
// 	protoc             v4.24.1<br />
// source: api/user.proto<br />

import Axios from "axios";
import Long from "long";

export interface Error{
	code: number;
	msg: string;
}

export interface LoginReq{
	src_type: string;
	src: string;
	//when src_type is idcard or nickname,this can't be dynamic
	password_type: string;
	//when password_type is static this length must >=10
	//when password_type is dynamic and this is empty,means send dynamic password to email or tel.
	//when password_type is dynamic and this is not empty,means verify dynamic password.
	password: string;
}
function LoginReqToJson(msg: LoginReq): string{
	let s: string="{"
	//src_type
	if(msg.src_type==null||msg.src_type==undefined){
		throw 'LoginReq.src_type must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.src_type)
		s+='"src_type":'+vv+','
	}
	//src
	if(msg.src==null||msg.src==undefined){
		throw 'LoginReq.src must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.src)
		s+='"src":'+vv+','
	}
	//password_type
	if(msg.password_type==null||msg.password_type==undefined){
		throw 'LoginReq.password_type must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.password_type)
		s+='"password_type":'+vv+','
	}
	//password
	if(msg.password==null||msg.password==undefined){
		throw 'LoginReq.password must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.password)
		s+='"password":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface LoginResp{
	token: string;
	info: UserInfo|null|undefined;
	//verify:server already send the dynamic password to user's email or tel(depend on the login_req's src_type and src) and is waiting for verify
	//password:login success,but this account must finish the static password set
	//success:nothing need to do
	step: string;
}
function JsonToLoginResp(jsonobj: { [k:string]:any }): LoginResp{
	let obj: LoginResp={
		token:'',
		info:null,
		step:'',
	}
	//token
	if(jsonobj['token']!=null&&jsonobj['token']!=undefined){
		if(typeof jsonobj['token']!='string'){
			throw 'LoginResp.token must be string'
		}
		obj['token']=jsonobj['token']
	}
	//info
	if(jsonobj['info']!=null&&jsonobj['info']!=undefined){
		if(typeof jsonobj['info']!='object'){
			throw 'LoginResp.info must be UserInfo'
		}
		obj['info']=JsonToUserInfo(jsonobj['info'])
	}
	//step
	if(jsonobj['step']!=null&&jsonobj['step']!=undefined){
		if(typeof jsonobj['step']!='string'){
			throw 'LoginResp.step must be string'
		}
		obj['step']=jsonobj['step']
	}
	return obj
}
export interface UpdateEmailReq{
	old_receiver_type: string;
	//if this is empty,means send dynamic password
	//if this is not empty,means verify dynamic password
	old_dynamic_password: string;
	new_email: string;
	//if this is empty,means send dynamic password.
	//if this is not empty,means verify dynamic password.
	new_email_dynamic_password: string;
}
function UpdateEmailReqToJson(msg: UpdateEmailReq): string{
	let s: string="{"
	//old_receiver_type
	if(msg.old_receiver_type==null||msg.old_receiver_type==undefined){
		throw 'UpdateEmailReq.old_receiver_type must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.old_receiver_type)
		s+='"old_receiver_type":'+vv+','
	}
	//old_dynamic_password
	if(msg.old_dynamic_password==null||msg.old_dynamic_password==undefined){
		throw 'UpdateEmailReq.old_dynamic_password must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.old_dynamic_password)
		s+='"old_dynamic_password":'+vv+','
	}
	//new_email
	if(msg.new_email==null||msg.new_email==undefined){
		throw 'UpdateEmailReq.new_email must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.new_email)
		s+='"new_email":'+vv+','
	}
	//new_email_dynamic_password
	if(msg.new_email_dynamic_password==null||msg.new_email_dynamic_password==undefined){
		throw 'UpdateEmailReq.new_email_dynamic_password must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.new_email_dynamic_password)
		s+='"new_email_dynamic_password":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface UpdateEmailResp{
	//oldverify:server already send the dynamic password to user's email or tel(depend on the update_email_req's old_receiver_type) and is waiting for verify
	//newverify:server already send the dynamic password to the new email(depend on the update_email_req's new_email) and is waiting for verify
	//success:nothing need to do
	step: string;
}
function JsonToUpdateEmailResp(jsonobj: { [k:string]:any }): UpdateEmailResp{
	let obj: UpdateEmailResp={
		step:'',
	}
	//step
	if(jsonobj['step']!=null&&jsonobj['step']!=undefined){
		if(typeof jsonobj['step']!='string'){
			throw 'UpdateEmailResp.step must be string'
		}
		obj['step']=jsonobj['step']
	}
	return obj
}
export interface UpdateNickNameReq{
	new_nick_name: string;
}
function UpdateNickNameReqToJson(msg: UpdateNickNameReq): string{
	let s: string="{"
	//new_nick_name
	if(msg.new_nick_name==null||msg.new_nick_name==undefined){
		throw 'UpdateNickNameReq.new_nick_name must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.new_nick_name)
		s+='"new_nick_name":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface UpdateNickNameResp{
}
function JsonToUpdateNickNameResp(_jsonobj: { [k:string]:any }): UpdateNickNameResp{
	let obj: UpdateNickNameResp={
	}
	return obj
}
export interface UpdateStaticPasswordReq{
	//if this is empty,means this is the first time to set the static password
	old_static_password: string;
	new_static_password: string;
}
function UpdateStaticPasswordReqToJson(msg: UpdateStaticPasswordReq): string{
	let s: string="{"
	//old_static_password
	if(msg.old_static_password==null||msg.old_static_password==undefined){
		throw 'UpdateStaticPasswordReq.old_static_password must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.old_static_password)
		s+='"old_static_password":'+vv+','
	}
	//new_static_password
	if(msg.new_static_password==null||msg.new_static_password==undefined){
		throw 'UpdateStaticPasswordReq.new_static_password must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.new_static_password)
		s+='"new_static_password":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface UpdateStaticPasswordResp{
}
function JsonToUpdateStaticPasswordResp(_jsonobj: { [k:string]:any }): UpdateStaticPasswordResp{
	let obj: UpdateStaticPasswordResp={
	}
	return obj
}
export interface UpdateTelReq{
	old_receiver_type: string;
	//if this is empty,means send dynamic password
	//if this is not empty,means verify dynamic password
	old_dynamic_password: string;
	new_tel: string;
	//if this is empty,means send dynamic password.
	//if this is not empty,means verify dynamic password.
	new_tel_dynamic_password: string;
}
function UpdateTelReqToJson(msg: UpdateTelReq): string{
	let s: string="{"
	//old_receiver_type
	if(msg.old_receiver_type==null||msg.old_receiver_type==undefined){
		throw 'UpdateTelReq.old_receiver_type must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.old_receiver_type)
		s+='"old_receiver_type":'+vv+','
	}
	//old_dynamic_password
	if(msg.old_dynamic_password==null||msg.old_dynamic_password==undefined){
		throw 'UpdateTelReq.old_dynamic_password must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.old_dynamic_password)
		s+='"old_dynamic_password":'+vv+','
	}
	//new_tel
	if(msg.new_tel==null||msg.new_tel==undefined){
		throw 'UpdateTelReq.new_tel must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.new_tel)
		s+='"new_tel":'+vv+','
	}
	//new_tel_dynamic_password
	if(msg.new_tel_dynamic_password==null||msg.new_tel_dynamic_password==undefined){
		throw 'UpdateTelReq.new_tel_dynamic_password must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.new_tel_dynamic_password)
		s+='"new_tel_dynamic_password":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface UpdateTelResp{
	//oldverify:server already send the dynamic password to user's email or tel(depend on the update_tel_req's old_receiver_type) and is waiting for verify
	//newverify:server already send the dynamic password to the new tel(depend on the update_tel_req's new_tel) and is waiting for verify
	//success:nothing need to do
	step: string;
}
function JsonToUpdateTelResp(jsonobj: { [k:string]:any }): UpdateTelResp{
	let obj: UpdateTelResp={
		step:'',
	}
	//step
	if(jsonobj['step']!=null&&jsonobj['step']!=undefined){
		if(typeof jsonobj['step']!='string'){
			throw 'UpdateTelResp.step must be string'
		}
		obj['step']=jsonobj['step']
	}
	return obj
}
export interface UserInfo{
	user_id: string;
	idcard: string;
	tel: string;
	email: string;
	nick_name: string;
	//Warning!!!Type is int64,be careful of sign(+,-)
	ctime: Long;
	//Warning!!!map's value's type is int64,be careful of sign(+,-)
	money: Map<string,Long>|null|undefined;
}
function JsonToUserInfo(jsonobj: { [k:string]:any }): UserInfo{
	let obj: UserInfo={
		user_id:'',
		idcard:'',
		tel:'',
		email:'',
		nick_name:'',
		ctime:Long.ZERO,
		money:null,
	}
	//user_id
	if(jsonobj['user_id']!=null&&jsonobj['user_id']!=undefined){
		if(typeof jsonobj['user_id']!='string'){
			throw 'UserInfo.user_id must be string'
		}
		obj['user_id']=jsonobj['user_id']
	}
	//idcard
	if(jsonobj['idcard']!=null&&jsonobj['idcard']!=undefined){
		if(typeof jsonobj['idcard']!='string'){
			throw 'UserInfo.idcard must be string'
		}
		obj['idcard']=jsonobj['idcard']
	}
	//tel
	if(jsonobj['tel']!=null&&jsonobj['tel']!=undefined){
		if(typeof jsonobj['tel']!='string'){
			throw 'UserInfo.tel must be string'
		}
		obj['tel']=jsonobj['tel']
	}
	//email
	if(jsonobj['email']!=null&&jsonobj['email']!=undefined){
		if(typeof jsonobj['email']!='string'){
			throw 'UserInfo.email must be string'
		}
		obj['email']=jsonobj['email']
	}
	//nick_name
	if(jsonobj['nick_name']!=null&&jsonobj['nick_name']!=undefined){
		if(typeof jsonobj['nick_name']!='string'){
			throw 'UserInfo.nick_name must be string'
		}
		obj['nick_name']=jsonobj['nick_name']
	}
	//ctime
	if(jsonobj['ctime']!=null&&jsonobj['ctime']!=undefined){
		if(typeof jsonobj['ctime']=='number'){
			if(!Number.isInteger(jsonobj['ctime'])){
				throw 'UserInfo.ctime must be integer'
			}
			let tmp: Long=Long.ZERO
			try{
				tmp=Long.fromNumber(jsonobj['ctime'],false)
			}catch(e){
				throw 'UserInfo.ctime must be integer'
			}
			obj['ctime']=tmp
		}else if(typeof jsonobj['ctime']=='string'){
			let tmp:Long=Long.ZERO
			try{
				tmp=Long.fromString(jsonobj['ctime'],false)
			}catch(e){
				throw 'UserInfo.ctime must be integer'
			}
			if(tmp.toString()!=jsonobj['ctime']){
				throw 'UserInfo.ctime overflow'
			}
			obj['ctime']=tmp
		}else{
			throw 'UserInfo.ctime must be integer'
		}
	}
	//money
	if(jsonobj['money']!=null&&jsonobj['money']!=undefined){
		if(typeof jsonobj['money']!='object'){
			throw 'UserInfo.money must be Map<string,Long>|null|undefined'
		}
		for(let key of Object.keys(jsonobj['money'])){
			let value=jsonobj['money'][key]
			let k: string=key
			if(typeof value=='number'){
				if(!Number.isInteger(value)){
					throw 'value in UserInfo.money must be integer'
				}
			}else if(typeof value!='string'){
				throw 'value in UserInfo46money must be integer'
			}
			let v: Long=Long.ZERO
			if(typeof value=='number'){
				try{
					v=Long.fromNumber(value,false)
				}catch(e){
					throw 'value in UserInfo46money must be integer'
				}
			}else{
				try{
					v=Long.fromString(value,false)
				}catch(e){
					throw 'value in UserInfo.money must be integer'
				}
				if(v.toString()!=value){
					throw 'value in UserInfo.money overflow'
				}
			}
			if(obj['money']==undefined){
				obj['money']=new Map<string,Long>
			}
			obj['money'].set(k,v)
		}
	}
	return obj
}
const _WebPathUserLogin: string ="/account.user/login";
const _WebPathUserUpdateStaticPassword: string ="/account.user/update_static_password";
const _WebPathUserUpdateNickName: string ="/account.user/update_nick_name";
const _WebPathUserUpdateEmail: string ="/account.user/update_email";
const _WebPathUserUpdateTel: string ="/account.user/update_tel";
//ToC means this is used for users
export class UserBrowserClientToC {
	constructor(host: string){
		if(host==null||host==undefined||host.length==0){
			throw "UserBrowserClientToC's host missing"
		}
		this.host=host
	}
	//timeout must be integer,timeout's unit is millisecond
	//don't set Content-Type in header
	login(header: { [k: string]: string },req: LoginReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: LoginResp)=>void){
		if(!Number.isInteger(timeout)){
			errorf({code:-2,msg:'timeout must be integer'})
			return
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let body: string=''
		try{
			body=LoginReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathUserLogin,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:LoginResp=JsonToLoginResp(response.data)
				successf(obj)
			}catch(e){
				let err:Error={code:-1,msg:'response error'}
				errorf(err)
			}
		})
		.catch(function(error){
			if(error.response==undefined){
				errorf({code:-2,msg:error.message})
				return
			}
			let respdata=error.response.data
			let err:Error={code:-1,msg:''}
			if(respdata.code==undefined||typeof respdata.code!='number'||!Number.isInteger(respdata.code)||respdata.msg==undefined||typeof respdata.msg!='string'){
				err.msg=respdata
			}else{
				err.code=respdata.code
				err.msg=respdata.msg
			}
			errorf(err)
		})
	}
	//timeout must be integer,timeout's unit is millisecond
	//don't set Content-Type in header
	update_static_password(header: { [k: string]: string },req: UpdateStaticPasswordReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: UpdateStaticPasswordResp)=>void){
		if(!Number.isInteger(timeout)){
			errorf({code:-2,msg:'timeout must be integer'})
			return
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let body: string=''
		try{
			body=UpdateStaticPasswordReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathUserUpdateStaticPassword,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:UpdateStaticPasswordResp=JsonToUpdateStaticPasswordResp(response.data)
				successf(obj)
			}catch(e){
				let err:Error={code:-1,msg:'response error'}
				errorf(err)
			}
		})
		.catch(function(error){
			if(error.response==undefined){
				errorf({code:-2,msg:error.message})
				return
			}
			let respdata=error.response.data
			let err:Error={code:-1,msg:''}
			if(respdata.code==undefined||typeof respdata.code!='number'||!Number.isInteger(respdata.code)||respdata.msg==undefined||typeof respdata.msg!='string'){
				err.msg=respdata
			}else{
				err.code=respdata.code
				err.msg=respdata.msg
			}
			errorf(err)
		})
	}
	//timeout must be integer,timeout's unit is millisecond
	//don't set Content-Type in header
	update_nick_name(header: { [k: string]: string },req: UpdateNickNameReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: UpdateNickNameResp)=>void){
		if(!Number.isInteger(timeout)){
			errorf({code:-2,msg:'timeout must be integer'})
			return
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let body: string=''
		try{
			body=UpdateNickNameReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathUserUpdateNickName,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:UpdateNickNameResp=JsonToUpdateNickNameResp(response.data)
				successf(obj)
			}catch(e){
				let err:Error={code:-1,msg:'response error'}
				errorf(err)
			}
		})
		.catch(function(error){
			if(error.response==undefined){
				errorf({code:-2,msg:error.message})
				return
			}
			let respdata=error.response.data
			let err:Error={code:-1,msg:''}
			if(respdata.code==undefined||typeof respdata.code!='number'||!Number.isInteger(respdata.code)||respdata.msg==undefined||typeof respdata.msg!='string'){
				err.msg=respdata
			}else{
				err.code=respdata.code
				err.msg=respdata.msg
			}
			errorf(err)
		})
	}
	//timeout must be integer,timeout's unit is millisecond
	//don't set Content-Type in header
	update_email(header: { [k: string]: string },req: UpdateEmailReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: UpdateEmailResp)=>void){
		if(!Number.isInteger(timeout)){
			errorf({code:-2,msg:'timeout must be integer'})
			return
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let body: string=''
		try{
			body=UpdateEmailReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathUserUpdateEmail,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:UpdateEmailResp=JsonToUpdateEmailResp(response.data)
				successf(obj)
			}catch(e){
				let err:Error={code:-1,msg:'response error'}
				errorf(err)
			}
		})
		.catch(function(error){
			if(error.response==undefined){
				errorf({code:-2,msg:error.message})
				return
			}
			let respdata=error.response.data
			let err:Error={code:-1,msg:''}
			if(respdata.code==undefined||typeof respdata.code!='number'||!Number.isInteger(respdata.code)||respdata.msg==undefined||typeof respdata.msg!='string'){
				err.msg=respdata
			}else{
				err.code=respdata.code
				err.msg=respdata.msg
			}
			errorf(err)
		})
	}
	//timeout must be integer,timeout's unit is millisecond
	//don't set Content-Type in header
	update_tel(header: { [k: string]: string },req: UpdateTelReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: UpdateTelResp)=>void){
		if(!Number.isInteger(timeout)){
			errorf({code:-2,msg:'timeout must be integer'})
			return
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let body: string=''
		try{
			body=UpdateTelReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathUserUpdateTel,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:UpdateTelResp=JsonToUpdateTelResp(response.data)
				successf(obj)
			}catch(e){
				let err:Error={code:-1,msg:'response error'}
				errorf(err)
			}
		})
		.catch(function(error){
			if(error.response==undefined){
				errorf({code:-2,msg:error.message})
				return
			}
			let respdata=error.response.data
			let err:Error={code:-1,msg:''}
			if(respdata.code==undefined||typeof respdata.code!='number'||!Number.isInteger(respdata.code)||respdata.msg==undefined||typeof respdata.msg!='string'){
				err.msg=respdata
			}else{
				err.code=respdata.code
				err.msg=respdata.msg
			}
			errorf(err)
		})
	}
	private host: string
}
