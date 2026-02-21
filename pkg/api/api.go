package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/fatedier/frp/pkg/msg"
)

// Based on LoCyanFrp Frp modification
type Service struct{}

var (
	apiURL        = "https://user.mslmc.net/api/frp"
	httpTransport = &http.Transport{
		DisableKeepAlives: true,
	}
)

const userAgent = "MSLFrp/1.0 (Frps)"

func MyAPIService() (s *Service, err error) {
	return &Service{}, nil
}

// ProxyStartGetCfg 简单启动获取Cfg
func (s Service) ProxyStartGetCfg(userToken string, proxyID int) (cfg string, err error) {
	api, _ := url.Parse(apiURL + "/getTunnelConfig")
	values := url.Values{}
	values.Set("id", strconv.Itoa(proxyID))
	values.Set("userToken", userToken)
	// Encode 请求参数
	api.RawQuery = values.Encode()
	defer func(u *url.URL) {
		u.RawQuery = ""
	}(api)

	client := &http.Client{Transport: httpTransport}

	req, err := http.NewRequest(http.MethodGet, api.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 非正常状态码
	if resp.StatusCode != http.StatusOK {
		errInfo := ResError{}
		if err = json.Unmarshal(body, &errInfo); err != nil {
			return "", err
		}
		return "", errInfo
	}

	response := ResGetProxyCfg{}
	if err = json.Unmarshal(body, &response); err != nil {
		return "", err
	}
	return response.Data, nil
}

// SubmitRunId 提交runID至服务器
func (s Service) SubmitRunID(apiToken string, nodeID int, pMsg *msg.NewProxy, runID string) (err error) {
	api, _ := url.Parse(apiURL + "/server/run-id")
	values := url.Values{}

	name := strings.Split(pMsg.ProxyName, ".")[1]

	values.Set("run_id", runID)
	values.Set("proxy_name", name)
	values.Set("api_token", apiToken+"|"+strconv.Itoa(nodeID))

	client := &http.Client{Transport: httpTransport}

	req, err := http.NewRequest(http.MethodPost, api.String(), strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

type PxyMsg struct {
	ServerToken  string
	UserToken    string
	ProxyName    string
	RemotePort   int
	Type         string
	CustomDomain []string
}

// TunnelCheck 校验隧道
func (s Service) VerifyTunnel(pxyMsg PxyMsg) (retStr string, err error) {
	urlStr := apiURL + "/verifyTunnel?" +
		"token=" + pxyMsg.ServerToken +
		"&userToken=" + pxyMsg.UserToken + "&name=" + pxyMsg.ProxyName +
		"&remotePort=" + strconv.Itoa(pxyMsg.RemotePort)

	switch pxyMsg.Type {
	case "http":
		urlStr = apiURL + "/verifyTunnel?" +
			"token=" + pxyMsg.ServerToken +
			"&userToken=" + pxyMsg.UserToken + "&name=" + pxyMsg.ProxyName +
			"&remotePort=80&bindDomain=" + strings.Join(pxyMsg.CustomDomain, "|")
	case "https":
		urlStr = apiURL + "/verifyTunnel?" +
			"token=" + pxyMsg.ServerToken +
			"&userToken=" + pxyMsg.UserToken + "&name=" + pxyMsg.ProxyName +
			"&remotePort=443&bindDomain=" + strings.Join(pxyMsg.CustomDomain, "|")
	}
	api, _ := url.Parse(urlStr)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, api.String(), nil)
	if err != nil {
		return "request creation failed", err
	}
	req.Header.Set("User-Agent", userAgent)
	res, err := client.Do(req)
	if err != nil {
		return "request failed", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "failed: Status Code " + strconv.Itoa(res.StatusCode), errors.New(strconv.Itoa(res.StatusCode))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "failed to read response body", err
	}

	var jsonResponse struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	if _err := json.Unmarshal(body, &jsonResponse); _err != nil {
		return "failed to parse JSON response", _err
	}

	if jsonResponse.Code != 200 {
		return jsonResponse.Msg, errors.New("ERROR: Status Code " + strconv.Itoa(jsonResponse.Code))
	}

	return "", nil
}

// ProxyCheck 校验客户端代理
func (s Service) ProxyCheck(frpToken string, pMsg *msg.NewProxy, apiToken string, nodeID int) (ok bool, err error) {
	api, _ := url.Parse(apiURL + "/server/proxy")
	domains, err := json.Marshal(pMsg.CustomDomains)
	if err != nil {
		return false, err
	}

	values := url.Values{}

	name := strings.Split(pMsg.ProxyName, ".")[1]

	// API Basic
	values.Set("frp_token", frpToken)
	values.Set("api_token", apiToken+"|"+strconv.Itoa(nodeID))

	// Proxies basic info
	values.Set("proxy_name", name)
	values.Set("proxy_type", pMsg.ProxyType)
	values.Set("use_encryption", BoolToString(pMsg.UseEncryption))
	values.Set("use_compression", BoolToString(pMsg.UseCompression))

	// Http Proxies
	values.Set("domain", string(domains))

	// TCP & UDP & STCP
	values.Set("remote_port", strconv.Itoa(pMsg.RemotePort))

	// STCP & XTCP
	values.Set("secret_key", pMsg.Sk)

	api.RawQuery = values.Encode()
	defer func(u *url.URL) {
		u.RawQuery = ""
	}(api)

	client := &http.Client{Transport: httpTransport}

	req, err := http.NewRequest(http.MethodGet, api.String(), nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	// 非正常状态码
	if resp.StatusCode != http.StatusOK {
		errInfo := ResError{}
		if err = json.Unmarshal(body, &errInfo); err != nil {
			return false, err
		}
		return false, errInfo
	}

	response := ResCheckProxy{}
	if err = json.Unmarshal(body, &response); err != nil {
		return false, err
	}
	return true, nil
}

// GetLimit 获取隧道限速信息
func (s Service) GetLimit(frpsToken string, userToken string) (inLimit, outLimit uint64, err error) {
	api, _ := url.Parse(apiURL + "/getLimit")
	values := url.Values{}
	values.Set("token", frpsToken)
	values.Set("userToken", userToken)
	api.RawQuery = values.Encode()
	defer func(u *url.URL) {
		u.RawQuery = ""
	}(api)

	client := &http.Client{Transport: httpTransport}

	req, err := http.NewRequest(http.MethodGet, api.String(), nil)
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	// 非正常状态码
	if resp.StatusCode != http.StatusOK {
		return 0, 0, errors.New(strconv.Itoa(resp.StatusCode))
	}

	response := ResGetLimit{}
	if err = json.Unmarshal(body, &response); err != nil {
		return 0, 0, err
	}
	if response.Code != 200 {
		return 0, 0, errors.New("Status Code " + strconv.Itoa(response.Code))
	}

	// 这里直接返回 uint64 应该问题不大
	return response.Data.Inbound, response.Data.Outbound, nil
}
