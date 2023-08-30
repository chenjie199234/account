// Code generated by protoc-gen-go-pbex. DO NOT EDIT.
// version:
// 	protoc-gen-pbex v0.0.77
// 	protoc         v4.24.1
// source: api/user.proto

package api

// return empty means pass
func (m *GetUserInfoReq) Validate() (errstr string) {
	if m.GetSrcType() != "user_id" && m.GetSrcType() != "tel" && m.GetSrcType() != "email" && m.GetSrcType() != "idcard" {
		return "field: src_type in object: get_user_info_req check value str in failed"
	}
	if len(m.GetSrc()) <= 0 {
		return "field: src in object: get_user_info_req check value str len gt failed"
	}
	return ""
}

// return empty means pass
func (m *LoginReq) Validate() (errstr string) {
	if m.GetSrcType() != "tel" && m.GetSrcType() != "email" && m.GetSrcType() != "idcard" && m.GetSrcType() != "nickname" {
		return "field: src_type in object: login_req check value str in failed"
	}
	if len(m.GetSrc()) <= 0 {
		return "field: src in object: login_req check value str len gt failed"
	}
	if m.GetPasswordType() != "static" && m.GetPasswordType() != "dynamic" {
		return "field: password_type in object: login_req check value str in failed"
	}
	return ""
}

// return empty means pass
func (m *UpdateStaticPasswordReq) Validate() (errstr string) {
	if len(m.GetNewStaticPassword()) < 10 {
		return "field: new_static_password in object: update_static_password_req check value str len gte failed"
	}
	return ""
}

// return empty means pass
func (m *UpdateNickNameReq) Validate() (errstr string) {
	if len(m.GetNewNickName()) <= 0 {
		return "field: new_nick_name in object: update_nick_name_req check value str len gt failed"
	}
	return ""
}

// return empty means pass
func (m *UpdateEmailReq) Validate() (errstr string) {
	if m.GetOldReceiverType() != "email" && m.GetOldReceiverType() != "tel" {
		return "field: old_receiver_type in object: update_email_req check value str in failed"
	}
	if len(m.GetNewEmail()) <= 0 {
		return "field: new_email in object: update_email_req check value str len gt failed"
	}
	return ""
}

// return empty means pass
func (m *UpdateTelReq) Validate() (errstr string) {
	if m.GetOldReceiverType() != "email" && m.GetOldReceiverType() != "tel" {
		return "field: old_receiver_type in object: update_tel_req check value str in failed"
	}
	if len(m.GetNewTel()) <= 0 {
		return "field: new_tel in object: update_tel_req check value str len gt failed"
	}
	return ""
}
