// Code generated by protoc-gen-browser. DO NOT EDIT.
// version:
// 	protoc-gen-browser v0.0.108<br />
// 	protoc             v4.25.3<br />
// source: api/account_money.proto<br />

export interface LogicError{
	code: number;
	msg: string;
}

export class GetMoneyLogsReq{
	src_type: string = ''
	src: string = ''
	//0:return all logs
	//>0:return the required page's data
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	start_time: number = 0//unit second
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	end_time: number = 0//unit second
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	page: number = 0
	action: string = ''
	toJSON(){
		let tmp = {}
		if(this.src_type){
			tmp["src_type"]=this.src_type
		}
		if(this.src){
			tmp["src"]=this.src
		}
		if(this.start_time){
			tmp["start_time"]=this.start_time
		}
		if(this.end_time){
			tmp["end_time"]=this.end_time
		}
		if(this.page){
			tmp["page"]=this.page
		}
		if(this.action){
			tmp["action"]=this.action
		}
		return tmp
	}
}
export class GetMoneyLogsResp{
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	page: number = 0
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	pagesize: number = 0
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	totalsize: number = 0
	logs: Array<MoneyLog|null>|null = null
	fromOBJ(obj:Object){
		if(obj["page"]){
			this.page=obj["page"]
		}
		if(obj["pagesize"]){
			this.pagesize=obj["pagesize"]
		}
		if(obj["totalsize"]){
			this.totalsize=obj["totalsize"]
		}
		if(obj["logs"] && obj["logs"].length>0){
			this.logs=new Array<MoneyLog|null>()
			for(let value of obj["logs"]){
				if(value){
					let tmp=new MoneyLog()
					tmp.fromOBJ(value)
					this.logs.push(tmp)
				}else{
					this.logs.push(null)
				}
			}
		}
	}
}
export class MoneyLog{
	user_id: string = ''
	action: string = ''//spend,recharge,refund
	unique_id: string = ''
	src_dst: string = ''
	money_type: string = ''
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	money_amount: number = 0
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	ctime: number = 0
	fromOBJ(obj:Object){
		if(obj["user_id"]){
			this.user_id=obj["user_id"]
		}
		if(obj["action"]){
			this.action=obj["action"]
		}
		if(obj["unique_id"]){
			this.unique_id=obj["unique_id"]
		}
		if(obj["src_dst"]){
			this.src_dst=obj["src_dst"]
		}
		if(obj["money_type"]){
			this.money_type=obj["money_type"]
		}
		if(obj["money_amount"]){
			this.money_amount=obj["money_amount"]
		}
		if(obj["ctime"]){
			this.ctime=obj["ctime"]
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
		success(JSON.parse(d.data))
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
const _WebPathMoneyGetMoneyLogs: string ="/account.money/get_money_logs";
//ToB means this is for internal
//ToB client must be used with https://github.com/chenjie199234/admin
//If your are not using 'admin' as your tob request's proxy gate,don't use this
export class MoneyBrowserClientToB {
	constructor(proxyhost: string,serverprojectid: Array<number>,servergroup: string){
		if(!proxyhost || proxyhost.length==0){
			throw "MoneyBrowserClientToB's proxyhost missing"
		}
		if(!serverprojectid || serverprojectid.length!=2){
			throw "MoneyBrowserClientToB's serverprojectid missing or wrong"
		}
		if(!servergroup || servergroup.length==0){
			throw "MoneyBrowserClientToB's servergroup missing"
		}
		this.host=proxyhost
		this.projectid=serverprojectid
		this.group=servergroup
	}
	//timeout's unit is millisecond,it will be used when > 0
	get_money_logs(header: Object,req: GetMoneyLogsReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: GetMoneyLogsResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		let realreq = {
			project_id:this.projectid,
			g_name:this.group,
			a_name:"account",
			path:_WebPathMoneyGetMoneyLogs,
			data:JSON.stringify(req),
		}
		call(timeout,this.host+"/admin.app/proxy",{method:"POST",headers:header,body:JSON.stringify(realreq)},error,function(arg: Object){
			let r=new GetMoneyLogsResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	private host: string
	private projectid: Array<number>
	private group: string
}
