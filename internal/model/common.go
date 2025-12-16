// internal/model/common.go
package model

type UsersMobile struct {
	Uid        int64  `json:"uid"`
	Mid        string `json:"mid"`
	Region     string `json:"region"`
	Encrypt    string `json:"encrypt"`
	CreateTime int    `json:"create_time"`
}
