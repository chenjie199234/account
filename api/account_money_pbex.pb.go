// Code generated by protoc-gen-go-pbex. DO NOT EDIT.
// version:
// 	protoc-gen-pbex v0.0.81
// 	protoc         v4.24.1
// source: api/account_money.proto

package api

// return empty means pass
func (m *GetUserMoneyLogsReq) Validate() (errstr string) {
	if m.GetSrcType() != "user_id" && m.GetSrcType() != "tel" && m.GetSrcType() != "email" && m.GetSrcType() != "idcard" && m.GetSrcType() != "nickname" {
		return "field: src_type in object: get_user_money_logs_req check value str in failed"
	}
	if len(m.GetSrc()) <= 0 {
		return "field: src in object: get_user_money_logs_req check value str len gt failed"
	}
	if m.GetStartTime() <= 0 {
		return "field: start_time in object: get_user_money_logs_req check value uint gt failed"
	}
	if m.GetEndTime() <= 0 {
		return "field: end_time in object: get_user_money_logs_req check value uint gt failed"
	}
	if m.GetAction() != "spend" && m.GetAction() != "recharge" && m.GetAction() != "refund" && m.GetAction() != "all" {
		return "field: action in object: get_user_money_logs_req check value str in failed"
	}
	return ""
}

// return empty means pass
func (m *SelfMoneyLogsReq) Validate() (errstr string) {
	if m.GetStartTime() <= 0 {
		return "field: start_time in object: self_money_logs_req check value uint gt failed"
	}
	if m.GetEndTime() <= 0 {
		return "field: end_time in object: self_money_logs_req check value uint gt failed"
	}
	if m.GetAction() != "spend" && m.GetAction() != "recharge" && m.GetAction() != "refund" && m.GetAction() != "all" {
		return "field: action in object: self_money_logs_req check value str in failed"
	}
	return ""
}