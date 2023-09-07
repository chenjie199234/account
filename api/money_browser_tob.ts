// Code generated by protoc-gen-browser. DO NOT EDIT.
// version:
// 	protoc-gen-browser v0.0.79<br />
// 	protoc             v4.24.1<br />
// source: api/money.proto<br />

import Axios from "axios";

export interface Error{
	code: number;
	msg: string;
}

export interface GetUserMoneyLogsReq{
	src_type: string;
	src: string;
	//0:return all logs
	//>0:return the required page's data
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	start_time: number;//unit second
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	end_time: number;//unit second
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	page: number;
	action: string;
}
function GetUserMoneyLogsReqToJson(msg: GetUserMoneyLogsReq): string{
	let s: string="{"
	//src_type
	if(msg.src_type==null||msg.src_type==undefined){
		throw 'GetUserMoneyLogsReq.src_type must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.src_type)
		s+='"src_type":'+vv+','
	}
	//src
	if(msg.src==null||msg.src==undefined){
		throw 'GetUserMoneyLogsReq.src must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.src)
		s+='"src":'+vv+','
	}
	//start_time
	if(msg.start_time==null||msg.start_time==undefined||!Number.isInteger(msg.start_time)){
		throw 'GetUserMoneyLogsReq.start_time must be integer'
	}else if(msg.start_time>4294967295||msg.start_time<0){
		throw 'GetUserMoneyLogsReq.start_time overflow'
	}else{
		s+='"start_time":'+msg.start_time+','
	}
	//end_time
	if(msg.end_time==null||msg.end_time==undefined||!Number.isInteger(msg.end_time)){
		throw 'GetUserMoneyLogsReq.end_time must be integer'
	}else if(msg.end_time>4294967295||msg.end_time<0){
		throw 'GetUserMoneyLogsReq.end_time overflow'
	}else{
		s+='"end_time":'+msg.end_time+','
	}
	//page
	if(msg.page==null||msg.page==undefined||!Number.isInteger(msg.page)){
		throw 'GetUserMoneyLogsReq.page must be integer'
	}else if(msg.page>4294967295||msg.page<0){
		throw 'GetUserMoneyLogsReq.page overflow'
	}else{
		s+='"page":'+msg.page+','
	}
	//action
	if(msg.action==null||msg.action==undefined){
		throw 'GetUserMoneyLogsReq.action must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.action)
		s+='"action":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface GetUserMoneyLogsResp{
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	page: number;
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	pagesize: number;
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	totalsize: number;
	logs: Array<MoneyLog|null|undefined>|null|undefined;
}
function JsonToGetUserMoneyLogsResp(jsonobj: { [k:string]:any }): GetUserMoneyLogsResp{
	let obj: GetUserMoneyLogsResp={
		page:0,
		pagesize:0,
		totalsize:0,
		logs:null,
	}
	//page
	if(jsonobj['page']!=null&&jsonobj['page']!=undefined){
		if(typeof jsonobj['page']!='number'||!Number.isInteger(jsonobj['page'])){
			throw 'GetUserMoneyLogsResp.page must be integer'
		}else if(jsonobj['page']>4294967295||jsonobj['page']<0){
			throw 'GetUserMoneyLogsResp.page overflow'
		}
		obj['page']=jsonobj['page']
	}
	//pagesize
	if(jsonobj['pagesize']!=null&&jsonobj['pagesize']!=undefined){
		if(typeof jsonobj['pagesize']!='number'||!Number.isInteger(jsonobj['pagesize'])){
			throw 'GetUserMoneyLogsResp.pagesize must be integer'
		}else if(jsonobj['pagesize']>4294967295||jsonobj['pagesize']<0){
			throw 'GetUserMoneyLogsResp.pagesize overflow'
		}
		obj['pagesize']=jsonobj['pagesize']
	}
	//totalsize
	if(jsonobj['totalsize']!=null&&jsonobj['totalsize']!=undefined){
		if(typeof jsonobj['totalsize']!='number'||!Number.isInteger(jsonobj['totalsize'])){
			throw 'GetUserMoneyLogsResp.totalsize must be integer'
		}else if(jsonobj['totalsize']>4294967295||jsonobj['totalsize']<0){
			throw 'GetUserMoneyLogsResp.totalsize overflow'
		}
		obj['totalsize']=jsonobj['totalsize']
	}
	//logs
	if(jsonobj['logs']!=null&&jsonobj['logs']!=undefined){
		if(!(jsonobj['logs'] instanceof Array)){
			throw 'GetUserMoneyLogsResp.logs must be Array<MoneyLog>|null|undefined'
		}
		for(let element of jsonobj['logs']){
			if(typeof element!='object'){
				throw 'element in GetUserMoneyLogsResp.logs must be MoneyLog'
			}
			if(obj['logs']==null){
				obj['logs']=new Array<MoneyLog>
			}
			obj['logs'].push(JsonToMoneyLog(element))
		}
	}
	return obj
}
export interface MoneyLog{
	user_id: string;
	action: string;//spend,recharge,refund
	unique_id: string;
	src_dst: string;
	money_type: string;
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	money_amount: number;
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	ctime: number;
}
function JsonToMoneyLog(jsonobj: { [k:string]:any }): MoneyLog{
	let obj: MoneyLog={
		user_id:'',
		action:'',
		unique_id:'',
		src_dst:'',
		money_type:'',
		money_amount:0,
		ctime:0,
	}
	//user_id
	if(jsonobj['user_id']!=null&&jsonobj['user_id']!=undefined){
		if(typeof jsonobj['user_id']!='string'){
			throw 'MoneyLog.user_id must be string'
		}
		obj['user_id']=jsonobj['user_id']
	}
	//action
	if(jsonobj['action']!=null&&jsonobj['action']!=undefined){
		if(typeof jsonobj['action']!='string'){
			throw 'MoneyLog.action must be string'
		}
		obj['action']=jsonobj['action']
	}
	//unique_id
	if(jsonobj['unique_id']!=null&&jsonobj['unique_id']!=undefined){
		if(typeof jsonobj['unique_id']!='string'){
			throw 'MoneyLog.unique_id must be string'
		}
		obj['unique_id']=jsonobj['unique_id']
	}
	//src_dst
	if(jsonobj['src_dst']!=null&&jsonobj['src_dst']!=undefined){
		if(typeof jsonobj['src_dst']!='string'){
			throw 'MoneyLog.src_dst must be string'
		}
		obj['src_dst']=jsonobj['src_dst']
	}
	//money_type
	if(jsonobj['money_type']!=null&&jsonobj['money_type']!=undefined){
		if(typeof jsonobj['money_type']!='string'){
			throw 'MoneyLog.money_type must be string'
		}
		obj['money_type']=jsonobj['money_type']
	}
	//money_amount
	if(jsonobj['money_amount']!=null&&jsonobj['money_amount']!=undefined){
		if(typeof jsonobj['money_amount']!='number'||!Number.isInteger(jsonobj['money_amount'])){
			throw 'MoneyLog.money_amount must be integer'
		}else if(jsonobj['money_amount']>4294967295||jsonobj['money_amount']<0){
			throw 'MoneyLog.money_amount overflow'
		}
		obj['money_amount']=jsonobj['money_amount']
	}
	//ctime
	if(jsonobj['ctime']!=null&&jsonobj['ctime']!=undefined){
		if(typeof jsonobj['ctime']!='number'||!Number.isInteger(jsonobj['ctime'])){
			throw 'MoneyLog.ctime must be integer'
		}else if(jsonobj['ctime']>4294967295||jsonobj['ctime']<0){
			throw 'MoneyLog.ctime overflow'
		}
		obj['ctime']=jsonobj['ctime']
	}
	return obj
}
export interface RechargeMoneyReq{
}
function RechargeMoneyReqToJson(_msg: RechargeMoneyReq): string{
	let s: string="{"
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface RechargeMoneyResp{
}
function JsonToRechargeMoneyResp(_jsonobj: { [k:string]:any }): RechargeMoneyResp{
	let obj: RechargeMoneyResp={
	}
	return obj
}
export interface RefundMoneyReq{
}
function RefundMoneyReqToJson(_msg: RefundMoneyReq): string{
	let s: string="{"
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface RefundMoneyResp{
}
function JsonToRefundMoneyResp(_jsonobj: { [k:string]:any }): RefundMoneyResp{
	let obj: RefundMoneyResp={
	}
	return obj
}
export interface SelfMoneyLogsReq{
	//0:return all logs
	//>0:return the required page's data
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	page: number;
	action: string;
}
function SelfMoneyLogsReqToJson(msg: SelfMoneyLogsReq): string{
	let s: string="{"
	//page
	if(msg.page==null||msg.page==undefined||!Number.isInteger(msg.page)){
		throw 'SelfMoneyLogsReq.page must be integer'
	}else if(msg.page>4294967295||msg.page<0){
		throw 'SelfMoneyLogsReq.page overflow'
	}else{
		s+='"page":'+msg.page+','
	}
	//action
	if(msg.action==null||msg.action==undefined){
		throw 'SelfMoneyLogsReq.action must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.action)
		s+='"action":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface SelfMoneyLogsResp{
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	page: number;
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	pagesize: number;
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	totalsize: number;
	logs: Array<MoneyLog|null|undefined>|null|undefined;
}
function JsonToSelfMoneyLogsResp(jsonobj: { [k:string]:any }): SelfMoneyLogsResp{
	let obj: SelfMoneyLogsResp={
		page:0,
		pagesize:0,
		totalsize:0,
		logs:null,
	}
	//page
	if(jsonobj['page']!=null&&jsonobj['page']!=undefined){
		if(typeof jsonobj['page']!='number'||!Number.isInteger(jsonobj['page'])){
			throw 'SelfMoneyLogsResp.page must be integer'
		}else if(jsonobj['page']>4294967295||jsonobj['page']<0){
			throw 'SelfMoneyLogsResp.page overflow'
		}
		obj['page']=jsonobj['page']
	}
	//pagesize
	if(jsonobj['pagesize']!=null&&jsonobj['pagesize']!=undefined){
		if(typeof jsonobj['pagesize']!='number'||!Number.isInteger(jsonobj['pagesize'])){
			throw 'SelfMoneyLogsResp.pagesize must be integer'
		}else if(jsonobj['pagesize']>4294967295||jsonobj['pagesize']<0){
			throw 'SelfMoneyLogsResp.pagesize overflow'
		}
		obj['pagesize']=jsonobj['pagesize']
	}
	//totalsize
	if(jsonobj['totalsize']!=null&&jsonobj['totalsize']!=undefined){
		if(typeof jsonobj['totalsize']!='number'||!Number.isInteger(jsonobj['totalsize'])){
			throw 'SelfMoneyLogsResp.totalsize must be integer'
		}else if(jsonobj['totalsize']>4294967295||jsonobj['totalsize']<0){
			throw 'SelfMoneyLogsResp.totalsize overflow'
		}
		obj['totalsize']=jsonobj['totalsize']
	}
	//logs
	if(jsonobj['logs']!=null&&jsonobj['logs']!=undefined){
		if(!(jsonobj['logs'] instanceof Array)){
			throw 'SelfMoneyLogsResp.logs must be Array<MoneyLog>|null|undefined'
		}
		for(let element of jsonobj['logs']){
			if(typeof element!='object'){
				throw 'element in SelfMoneyLogsResp.logs must be MoneyLog'
			}
			if(obj['logs']==null){
				obj['logs']=new Array<MoneyLog>
			}
			obj['logs'].push(JsonToMoneyLog(element))
		}
	}
	return obj
}
export interface SpendMoneyReq{
}
function SpendMoneyReqToJson(_msg: SpendMoneyReq): string{
	let s: string="{"
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface SpendMoneyResp{
}
function JsonToSpendMoneyResp(_jsonobj: { [k:string]:any }): SpendMoneyResp{
	let obj: SpendMoneyResp={
	}
	return obj
}
const _WebPathMoneyGetUserMoneyLogs: string ="/account.money/get_user_money_logs";
const _WebPathMoneySelfMoneyLogs: string ="/account.money/self_money_logs";
const _WebPathMoneyRechargeMoney: string ="/account.money/recharge_money";
const _WebPathMoneySpendMoney: string ="/account.money/spend_money";
const _WebPathMoneyRefundMoney: string ="/account.money/refund_money";
//ToB means this is used for internal
//ToB client must be used with https://github.com/chenjie199234/admin
//If your are not using 'admin' as your tob request's proxy gate,don't use this
export class MoneyBrowserClientToB {
	constructor(proxyhost: string,servergroup: string){
		if(proxyhost==null||proxyhost==undefined||proxyhost.length==0){
			throw "MoneyBrowserClientToB's proxyhost missing"
		}
		if(servergroup==null||servergroup==undefined||servergroup.length==0){
			throw "MoneyBrowserClientToB's servergroup missing"
		}
		this.host=proxyhost
		this.group=servergroup
	}
	//timeout must be integer,timeout's unit is millisecond
	//don't set Content-Type in header
	get_user_money_logs(header: { [k: string]: string },req: GetUserMoneyLogsReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: GetUserMoneyLogsResp)=>void){
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
				path:_WebPathMoneyGetUserMoneyLogs,
				appname:'account',
				groupname:this.group,
				data:GetUserMoneyLogsReqToJson(req),
			},
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:GetUserMoneyLogsResp=JsonToGetUserMoneyLogsResp(response.data.data)
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
	self_money_logs(header: { [k: string]: string },req: SelfMoneyLogsReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: SelfMoneyLogsResp)=>void){
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
				path:_WebPathMoneySelfMoneyLogs,
				appname:'account',
				groupname:this.group,
				data:SelfMoneyLogsReqToJson(req),
			},
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:SelfMoneyLogsResp=JsonToSelfMoneyLogsResp(response.data.data)
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
	recharge_money(header: { [k: string]: string },req: RechargeMoneyReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: RechargeMoneyResp)=>void){
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
				path:_WebPathMoneyRechargeMoney,
				appname:'account',
				groupname:this.group,
				data:RechargeMoneyReqToJson(req),
			},
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:RechargeMoneyResp=JsonToRechargeMoneyResp(response.data.data)
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
	spend_money(header: { [k: string]: string },req: SpendMoneyReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: SpendMoneyResp)=>void){
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
				path:_WebPathMoneySpendMoney,
				appname:'account',
				groupname:this.group,
				data:SpendMoneyReqToJson(req),
			},
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:SpendMoneyResp=JsonToSpendMoneyResp(response.data.data)
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
	refund_money(header: { [k: string]: string },req: RefundMoneyReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: RefundMoneyResp)=>void){
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
				path:_WebPathMoneyRefundMoney,
				appname:'account',
				groupname:this.group,
				data:RefundMoneyReqToJson(req),
			},
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			try{
				let obj:RefundMoneyResp=JsonToRefundMoneyResp(response.data.data)
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
