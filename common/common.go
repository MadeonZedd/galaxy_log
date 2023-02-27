package common

import "encoding/json"

const TimeFormat = "2006-01-02 15:04:05"
const (
	CompanyIDKey  = "CompanyID"
	UserIDKey     = "UserID"
	PermissionKey = "Permission"
	EmptyString   = ""
)

func ToJsonString(v interface{}) string {
	if ret, err := Marshal(v); err != nil {
		return err.Error()
	} else {
		return string(ret)
	}
}

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
