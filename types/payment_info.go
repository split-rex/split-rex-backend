package types

import (
	"database/sql/driver"
	"encoding/json"
)

type PaymentInfo map[string]map[int]string
//{"payment_info": 
// 	{
// 		account_number:"account_name"
// 	}
//}

func (paymentInfo *PaymentInfo) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &paymentInfo)
}

func (paymentInfo PaymentInfo) Value() (driver.Value, error) {
	val, err := json.Marshal(paymentInfo)
	return string(val), err
}

func (paymentInfo PaymentInfo) GormDataType() string {
	return "string"
}

func (paymentInfo PaymentInfo) Count() int {
	return len(paymentInfo)
}
