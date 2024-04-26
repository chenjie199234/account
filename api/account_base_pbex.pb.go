// Code generated by protoc-gen-go-pbex. DO NOT EDIT.
// version:
// 	protoc-gen-pbex v0.0.115
// 	protoc         v5.26.1
// source: api/account_base.proto

package api

// return empty means pass
func (m *BaseInfoReq) Validate() (errstr string) {
	if m.GetSrcType() != "user_id" && m.GetSrcType() != "tel" && m.GetSrcType() != "email" && m.GetSrcType() != "idcard" {
		return "field: src_type in object: base_info_req check value str in failed"
	}
	return ""
}

// return empty means pass
func (m *GetOauthUrlReq) Validate() (errstr string) {
	if m.GetOauthServiceName() != "wechat" {
		return "field: oauth_service_name in object: get_oauth_url_req check value str in failed"
	}
	return ""
}

// return empty means pass
func (m *LoginReq) Validate() (errstr string) {
	if m.GetSrcType() != "tel" && m.GetSrcType() != "email" && m.GetSrcType() != "idcard" && m.GetSrcType() != "oauth" {
		return "field: src_type in object: login_req check value str in failed"
	}
	if len(m.GetSrcTypeExtra()) <= 0 {
		return "field: src_type_extra in object: login_req check value str len gt failed"
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
func (m *UpdateOauthReq) Validate() (errstr string) {
	if m.GetVerifySrcType() != "email" && m.GetVerifySrcType() != "tel" && m.GetVerifySrcType() != "oauth" {
		return "field: verify_src_type in object: update_oauth_req check value str in failed"
	}
	return ""
}

// return empty means pass
func (m *DelOauthReq) Validate() (errstr string) {
	if m.GetVerifySrcType() != "email" && m.GetVerifySrcType() != "tel" && m.GetVerifySrcType() != "oauth" {
		return "field: verify_src_type in object: del_oauth_req check value str in failed"
	}
	if len(m.GetDelOauthServiceName()) <= 0 {
		return "field: del_oauth_service_name in object: del_oauth_req check value str len gt failed"
	}
	return ""
}

// return empty means pass
func (m *IdcardDuplicateCheckReq) Validate() (errstr string) {
	if len(m.GetIdcard()) <= 0 {
		return "field: idcard in object: idcard_duplicate_check_req check value str len gt failed"
	}
	return ""
}

// return empty means pass
func (m *UpdateIdcardReq) Validate() (errstr string) {
	if m.GetVerifySrcType() != "email" && m.GetVerifySrcType() != "tel" && m.GetVerifySrcType() != "oauth" {
		return "field: verify_src_type in object: update_idcard_req check value str in failed"
	}
	if len(m.GetNewIdcard()) <= 0 {
		return "field: new_idcard in object: update_idcard_req check value str len gt failed"
	}
	return ""
}

// return empty means pass
func (m *DelIdcardReq) Validate() (errstr string) {
	if m.GetVerifySrcType() != "email" && m.GetVerifySrcType() != "tel" && m.GetVerifySrcType() != "oauth" {
		return "field: verify_src_type in object: del_idcard_req check value str in failed"
	}
	return ""
}

// return empty means pass
func (m *EmailDuplicateCheckReq) Validate() (errstr string) {
	if len(m.GetEmail()) <= 0 {
		return "field: email in object: email_duplicate_check_req check value str len gt failed"
	}
	return ""
}

// return empty means pass
func (m *UpdateEmailReq) Validate() (errstr string) {
	if m.GetVerifySrcType() != "email" && m.GetVerifySrcType() != "tel" && m.GetVerifySrcType() != "oauth" {
		return "field: verify_src_type in object: update_email_req check value str in failed"
	}
	if len(m.GetNewEmail()) <= 0 {
		return "field: new_email in object: update_email_req check value str len gt failed"
	}
	return ""
}

// return empty means pass
func (m *DelEmailReq) Validate() (errstr string) {
	if m.GetVerifySrcType() != "email" && m.GetVerifySrcType() != "tel" && m.GetVerifySrcType() != "oauth" {
		return "field: verify_src_type in object: del_email_req check value str in failed"
	}
	return ""
}

// return empty means pass
func (m *TelDuplicateCheckReq) Validate() (errstr string) {
	if len(m.GetTel()) <= 0 {
		return "field: tel in object: tel_duplicate_check_req check value str len gt failed"
	}
	return ""
}

// return empty means pass
func (m *UpdateTelReq) Validate() (errstr string) {
	if m.GetVerifySrcType() != "email" && m.GetVerifySrcType() != "tel" && m.GetVerifySrcType() != "oauth" {
		return "field: verify_src_type in object: update_tel_req check value str in failed"
	}
	if len(m.GetNewTel()) <= 0 {
		return "field: new_tel in object: update_tel_req check value str len gt failed"
	}
	return ""
}

// return empty means pass
func (m *DelTelReq) Validate() (errstr string) {
	if m.GetVerifySrcType() != "email" && m.GetVerifySrcType() != "tel" && m.GetVerifySrcType() != "oauth" {
		return "field: verify_src_type in object: del_tel_req check value str in failed"
	}
	return ""
}
