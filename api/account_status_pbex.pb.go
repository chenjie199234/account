// Code generated by protoc-gen-go-pbex. DO NOT EDIT.
// version:
// 	protoc-gen-pbex v0.0.116
// 	protoc         v5.27.0
// source: api/account_status.proto

package api

// return empty means pass
func (m *Pingreq) Validate() (errstr string) {
	if m.GetTimestamp() <= 0 {
		return "field: timestamp in object: pingreq check value int gt failed"
	}
	return ""
}
