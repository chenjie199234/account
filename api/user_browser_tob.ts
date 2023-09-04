// Code generated by protoc-gen-browser. DO NOT EDIT.
// version:
// 	protoc-gen-browser v0.0.78<br />
// 	protoc             v4.24.1<br />
// source: api/user.proto<br />

import Axios from "axios";

export interface Error{
	code: number;
	msg: string;
}

export interface EmailDuplicateCheckReq{
	email: string;
}
function EmailDuplicateCheckReqToJson(msg: EmailDuplicateCheckReq): string{
	let s: string="{"
	//email
	if(msg.email==null||msg.email==undefined){
		throw 'EmailDuplicateCheckReq.email must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.email)
		s+='"email":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface EmailDuplicateCheckResp{
	duplicate: boolean;
}
function JsonToEmailDuplicateCheckResp(jsonobj: { [k:string]:any }): EmailDuplicateCheckResp{
	let obj: EmailDuplicateCheckResp={
		duplicate:false,
	}
	//duplicate
	if(jsonobj['duplicate']!=null&&jsonobj['duplicate']!=undefined){
		if(typeof jsonobj['duplicate']!='boolean'){
			throw 'EmailDuplicateCheckResp.duplicate must be boolean'
		}
		obj['duplicate']=jsonobj['duplicate']
	}
	return obj
}
export interface GetUserInfoReq{
	src_type: string;
	src: string;
}
function GetUserInfoReqToJson(msg: GetUserInfoReq): string{
	let s: string="{"
	//src_type
	if(msg.src_type==null||msg.src_type==undefined){
		throw 'GetUserInfoReq.src_type must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.src_type)
		s+='"src_type":'+vv+','
	}
	//src
	if(msg.src==null||msg.src==undefined){
		throw 'GetUserInfoReq.src must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.src)
		s+='"src":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface GetUserInfoResp{
	info: UserInfo|null|undefined;
}
function JsonToGetUserInfoResp(jsonobj: { [k:string]:any }): GetUserInfoResp{
	let obj: GetUserInfoResp={
		info:null,
	}
	//info
	if(jsonobj['info']!=null&&jsonobj['info']!=undefined){
		if(typeof jsonobj['info']!='object'){
			throw 'GetUserInfoResp.info must be UserInfo'
		}
		obj['info']=JsonToUserInfo(jsonobj['info'])
	}
	return obj
}
export interface IdcardDuplicateCheckReq{
	idcard: string;
}
function IdcardDuplicateCheckReqToJson(msg: IdcardDuplicateCheckReq): string{
	let s: string="{"
	//idcard
	if(msg.idcard==null||msg.idcard==undefined){
		throw 'IdcardDuplicateCheckReq.idcard must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.idcard)
		s+='"idcard":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface IdcardDuplicateCheckResp{
	duplicate: boolean;
}
function JsonToIdcardDuplicateCheckResp(jsonobj: { [k:string]:any }): IdcardDuplicateCheckResp{
	let obj: IdcardDuplicateCheckResp={
		duplicate:false,
	}
	//duplicate
	if(jsonobj['duplicate']!=null&&jsonobj['duplicate']!=undefined){
		if(typeof jsonobj['duplicate']!='boolean'){
			throw 'IdcardDuplicateCheckResp.duplicate must be boolean'
		}
		obj['duplicate']=jsonobj['duplicate']
	}
	return obj
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
export interface NickNameDuplicateCheckReq{
	nick_name: string;
}
function NickNameDuplicateCheckReqToJson(msg: NickNameDuplicateCheckReq): string{
	let s: string="{"
	//nick_name
	if(msg.nick_name==null||msg.nick_name==undefined){
		throw 'NickNameDuplicateCheckReq.nick_name must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.nick_name)
		s+='"nick_name":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface NickNameDuplicateCheckResp{
	duplicate: boolean;
}
function JsonToNickNameDuplicateCheckResp(jsonobj: { [k:string]:any }): NickNameDuplicateCheckResp{
	let obj: NickNameDuplicateCheckResp={
		duplicate:false,
	}
	//duplicate
	if(jsonobj['duplicate']!=null&&jsonobj['duplicate']!=undefined){
		if(typeof jsonobj['duplicate']!='boolean'){
			throw 'NickNameDuplicateCheckResp.duplicate must be boolean'
		}
		obj['duplicate']=jsonobj['duplicate']
	}
	return obj
}
export interface SelfUserInfoReq{
}
function SelfUserInfoReqToJson(_msg: SelfUserInfoReq): string{
	let s: string="{"
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface SelfUserInfoResp{
	info: UserInfo|null|undefined;
}
function JsonToSelfUserInfoResp(jsonobj: { [k:string]:any }): SelfUserInfoResp{
	let obj: SelfUserInfoResp={
		info:null,
	}
	//info
	if(jsonobj['info']!=null&&jsonobj['info']!=undefined){
		if(typeof jsonobj['info']!='object'){
			throw 'SelfUserInfoResp.info must be UserInfo'
		}
		obj['info']=JsonToUserInfo(jsonobj['info'])
	}
	return obj
}
export interface TelDuplicateCheckReq{
	tel: string;
}
function TelDuplicateCheckReqToJson(msg: TelDuplicateCheckReq): string{
	let s: string="{"
	//tel
	if(msg.tel==null||msg.tel==undefined){
		throw 'TelDuplicateCheckReq.tel must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.tel)
		s+='"tel":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface TelDuplicateCheckResp{
	duplicate: boolean;
}
function JsonToTelDuplicateCheckResp(jsonobj: { [k:string]:any }): TelDuplicateCheckResp{
	let obj: TelDuplicateCheckResp={
		duplicate:false,
	}
	//duplicate
	if(jsonobj['duplicate']!=null&&jsonobj['duplicate']!=undefined){
		if(typeof jsonobj['duplicate']!='boolean'){
			throw 'TelDuplicateCheckResp.duplicate must be boolean'
		}
		obj['duplicate']=jsonobj['duplicate']
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
export interface UpdateIdcardReq{
	new_idcard: string;
}
function UpdateIdcardReqToJson(msg: UpdateIdcardReq): string{
	let s: string="{"
	//new_idcard
	if(msg.new_idcard==null||msg.new_idcard==undefined){
		throw 'UpdateIdcardReq.new_idcard must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.new_idcard)
		s+='"new_idcard":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface UpdateIdcardResp{
}
function JsonToUpdateIdcardResp(_jsonobj: { [k:string]:any }): UpdateIdcardResp{
	let obj: UpdateIdcardResp={
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
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	ctime: number;
	//Warning!!!map's value's type is int32,be careful of sign(+,-) and overflow
	money: Map<string,number>|null|undefined;
}
function JsonToUserInfo(jsonobj: { [k:string]:any }): UserInfo{
	let obj: UserInfo={
		user_id:'',
		idcard:'',
		tel:'',
		email:'',
		nick_name:'',
		ctime:0,
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
		if(typeof jsonobj['ctime']!='number'||!Number.isInteger(jsonobj['ctime'])){
			throw 'UserInfo.ctime must be integer'
		}else if(jsonobj['ctime']>4294967295||jsonobj['ctime']<0){
			throw 'UserInfo.ctime overflow'
		}
		obj['ctime']=jsonobj['ctime']
	}
	//money
	if(jsonobj['money']!=null&&jsonobj['money']!=undefined){
		if(typeof jsonobj['money']!='object'){
			throw 'UserInfo.money must be Map<string,number>|null|undefined'
		}
		for(let key of Object.keys(jsonobj['money'])){
			let value=jsonobj['money'][key]
			let k: string=key
			if(typeof value!='number'||!Number.isInteger(value)){
				throw 'value in UserInfo.money must be integer'
			}else if(value>2147483647&&value<-2147483648){
				throw 'value in UserInfo.money overflow'
			}
			let v: number=value
			if(obj['money']==undefined){
				obj['money']=new Map<string,number>
			}
			obj['money'].set(k,v)
		}
	}
	return obj
}
const _WebPathUserGetUserInfo: string ="/account.user/get_user_info";
const _WebPathUserLogin: string ="/account.user/login";
const _WebPathUserSelfUserInfo: string ="/account.user/self_user_info";
const _WebPathUserUpdateStaticPassword: string ="/account.user/update_static_password";
const _WebPathUserIdcardDuplicateCheck: string ="/account.user/idcard_duplicate_check";
const _WebPathUserUpdateIdcard: string ="/account.user/update_idcard";
const _WebPathUserNickNameDuplicateCheck: string ="/account.user/nick_name_duplicate_check";
const _WebPathUserUpdateNickName: string ="/account.user/update_nick_name";
const _WebPathUserEmailDuplicateCheck: string ="/account.user/email_duplicate_check";
const _WebPathUserUpdateEmail: string ="/account.user/update_email";
const _WebPathUserTelDuplicateCheck: string ="/account.user/tel_duplicate_check";
const _WebPathUserUpdateTel: string ="/account.user/update_tel";
//ToB means this is used for internal
//ToB client must be used with https://github.com/chenjie199234/admin
//If your are not using 'admin' as your tob request's proxy gate,don't use this
export class UserBrowserClientToB {
	constructor(proxyhost: string,servergroup: string){
		if(proxyhost==null||proxyhost==undefined||proxyhost.length==0){
			throw "UserBrowserClientToB's proxyhost missing"
		}
		if(servergroup==null||servergroup==undefined||servergroup.length==0){
			throw "UserBrowserClientToB's servergroup missing"
		}
		this.host=proxyhost
		this.group=servergroup
	}
	//timeout must be integer,timeout's unit is millisecond
	//don't set Content-Type in header
	get_user_info(header: { [k: string]: string },req: GetUserInfoReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: GetUserInfoResp)=>void){
		if(!Number.isInteger(timeout)){
			throw 'timeout must be integer'
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let config={
			url:'/admin.app/proxy',
			method: 'post',
			baseURL: this.host,
			headers: header,
			data:{
				path:_WebPathUserGetUserInfo,
				appname:'account',
				groupname:this.group,
				data:GetUserInfoReqToJson(req),
			},
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:GetUserInfoResp=JsonToGetUserInfoResp(response.data.data)
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
	login(header: { [k: string]: string },req: LoginReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: LoginResp)=>void){
		if(!Number.isInteger(timeout)){
			throw 'timeout must be integer'
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let config={
			url:'/admin.app/proxy',
			method: 'post',
			baseURL: this.host,
			headers: header,
			data:{
				path:_WebPathUserLogin,
				appname:'account',
				groupname:this.group,
				data:LoginReqToJson(req),
			},
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:LoginResp=JsonToLoginResp(response.data.data)
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
	self_user_info(header: { [k: string]: string },req: SelfUserInfoReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: SelfUserInfoResp)=>void){
		if(!Number.isInteger(timeout)){
			throw 'timeout must be integer'
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let config={
			url:'/admin.app/proxy',
			method: 'post',
			baseURL: this.host,
			headers: header,
			data:{
				path:_WebPathUserSelfUserInfo,
				appname:'account',
				groupname:this.group,
				data:SelfUserInfoReqToJson(req),
			},
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:SelfUserInfoResp=JsonToSelfUserInfoResp(response.data.data)
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
			throw 'timeout must be integer'
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let config={
			url:'/admin.app/proxy',
			method: 'post',
			baseURL: this.host,
			headers: header,
			data:{
				path:_WebPathUserUpdateStaticPassword,
				appname:'account',
				groupname:this.group,
				data:UpdateStaticPasswordReqToJson(req),
			},
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:UpdateStaticPasswordResp=JsonToUpdateStaticPasswordResp(response.data.data)
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
	idcard_duplicate_check(header: { [k: string]: string },req: IdcardDuplicateCheckReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: IdcardDuplicateCheckResp)=>void){
		if(!Number.isInteger(timeout)){
			throw 'timeout must be integer'
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let config={
			url:'/admin.app/proxy',
			method: 'post',
			baseURL: this.host,
			headers: header,
			data:{
				path:_WebPathUserIdcardDuplicateCheck,
				appname:'account',
				groupname:this.group,
				data:IdcardDuplicateCheckReqToJson(req),
			},
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:IdcardDuplicateCheckResp=JsonToIdcardDuplicateCheckResp(response.data.data)
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
	update_idcard(header: { [k: string]: string },req: UpdateIdcardReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: UpdateIdcardResp)=>void){
		if(!Number.isInteger(timeout)){
			throw 'timeout must be integer'
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let config={
			url:'/admin.app/proxy',
			method: 'post',
			baseURL: this.host,
			headers: header,
			data:{
				path:_WebPathUserUpdateIdcard,
				appname:'account',
				groupname:this.group,
				data:UpdateIdcardReqToJson(req),
			},
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:UpdateIdcardResp=JsonToUpdateIdcardResp(response.data.data)
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
	nick_name_duplicate_check(header: { [k: string]: string },req: NickNameDuplicateCheckReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: NickNameDuplicateCheckResp)=>void){
		if(!Number.isInteger(timeout)){
			throw 'timeout must be integer'
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let config={
			url:'/admin.app/proxy',
			method: 'post',
			baseURL: this.host,
			headers: header,
			data:{
				path:_WebPathUserNickNameDuplicateCheck,
				appname:'account',
				groupname:this.group,
				data:NickNameDuplicateCheckReqToJson(req),
			},
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:NickNameDuplicateCheckResp=JsonToNickNameDuplicateCheckResp(response.data.data)
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
			throw 'timeout must be integer'
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let config={
			url:'/admin.app/proxy',
			method: 'post',
			baseURL: this.host,
			headers: header,
			data:{
				path:_WebPathUserUpdateNickName,
				appname:'account',
				groupname:this.group,
				data:UpdateNickNameReqToJson(req),
			},
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:UpdateNickNameResp=JsonToUpdateNickNameResp(response.data.data)
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
	email_duplicate_check(header: { [k: string]: string },req: EmailDuplicateCheckReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: EmailDuplicateCheckResp)=>void){
		if(!Number.isInteger(timeout)){
			throw 'timeout must be integer'
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let config={
			url:'/admin.app/proxy',
			method: 'post',
			baseURL: this.host,
			headers: header,
			data:{
				path:_WebPathUserEmailDuplicateCheck,
				appname:'account',
				groupname:this.group,
				data:EmailDuplicateCheckReqToJson(req),
			},
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:EmailDuplicateCheckResp=JsonToEmailDuplicateCheckResp(response.data.data)
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
			throw 'timeout must be integer'
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let config={
			url:'/admin.app/proxy',
			method: 'post',
			baseURL: this.host,
			headers: header,
			data:{
				path:_WebPathUserUpdateEmail,
				appname:'account',
				groupname:this.group,
				data:UpdateEmailReqToJson(req),
			},
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:UpdateEmailResp=JsonToUpdateEmailResp(response.data.data)
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
	tel_duplicate_check(header: { [k: string]: string },req: TelDuplicateCheckReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: TelDuplicateCheckResp)=>void){
		if(!Number.isInteger(timeout)){
			throw 'timeout must be integer'
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let config={
			url:'/admin.app/proxy',
			method: 'post',
			baseURL: this.host,
			headers: header,
			data:{
				path:_WebPathUserTelDuplicateCheck,
				appname:'account',
				groupname:this.group,
				data:TelDuplicateCheckReqToJson(req),
			},
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:TelDuplicateCheckResp=JsonToTelDuplicateCheckResp(response.data.data)
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
			throw 'timeout must be integer'
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let config={
			url:'/admin.app/proxy',
			method: 'post',
			baseURL: this.host,
			headers: header,
			data:{
				path:_WebPathUserUpdateTel,
				appname:'account',
				groupname:this.group,
				data:UpdateTelReqToJson(req),
			},
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:UpdateTelResp=JsonToUpdateTelResp(response.data.data)
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
	private group: string
}
