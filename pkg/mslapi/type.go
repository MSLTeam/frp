package mslapi

import "fmt"

type ResGetProxyCfg struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    string `json:"data"`
}
type ResCheckFrpToken struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct{}
}

type ResCheckProxy struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct{}
}

type ResGetLimit struct {
	Code int `json:"code"`
	Data struct {
		Inbound  uint64 `json:"inbound"`
		Outbound uint64 `json:"outbound"`
	}
}

type ResError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (r ResError) Error() string {
	return fmt.Sprintf("MSLF API Error (Status: %d, Message: %s)", r.Status, r.Message)
}
