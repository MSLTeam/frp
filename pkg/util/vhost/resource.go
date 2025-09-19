// Copyright 2017 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vhost

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/fatedier/frp/pkg/util/log"
	"github.com/fatedier/frp/pkg/util/version"
)

var NotFoundPagePath = ""

const (
	NotFound = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <title>404 Not Found</title>
  <style>
    .blur-bg {
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background-image: url("https://fastcdn.mihoyo.com/content-v2/hk4e/159347/d546a6006c42aebce823dc4922b05cd6_4956747493908915554.png");
      background-size: cover;
      background-repeat: no-repeat;
      background-position: center;
      filter: blur(5px); 
      -webkit-filter: blur(5px); 
      z-index: -1;
    }

    .container {
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      height: calc(100vh - 100px); /* 减去页脚的高度 */
      color: white;
    }

    h1 {
      font-size: 6rem;
      font-weight: bold;
      margin-bottom: 0;
    }

    p {
      font-size: 2rem;
      margin-top: 0;
    }

    a {
      display: inline-block;
      background-color: white;
      color: black;
      font-size: 2rem;
      padding: 1rem 2rem;
      text-decoration: none;
      border-radius: 2rem;
      margin-top: 2rem;
    }

    footer {
      position: fixed;
      bottom: 0;
      left: 0;
      width: 100%;
      height: 100px;
      background-color: rgba(0, 0, 0, 0.5); 
      display: flex;
      justify-content: center;
      align-items: center;
    }

    footer p {
      font-size: 1.5rem;
      color: white;
      margin: 0;
    }
  </style>
</head>
<body>
    <div class="blur-bg"></div> 
    <div class="container">
      <h1>MSL-Frp</h1>
      <p>这里什么也没有哦~</p>
      <a href="https://user.mslmc.net">返回MSL用户中心</a>
    </div>
    <footer>
      <p>Copyright © MSLTeam 2021-2025 All Right Received.</p> <!-- 添加页脚信息 -->
    </footer>
  </body>
  </html>
`
)

func getNotFoundPageContent() []byte {
	var (
		buf []byte
		err error
	)
	if NotFoundPagePath != "" {
		buf, err = os.ReadFile(NotFoundPagePath)
		if err != nil {
			log.Warnf("read custom 404 page error: %v", err)
			buf = []byte(NotFound)
		}
	} else {
		buf = []byte(NotFound)
	}
	return buf
}

func NotFoundResponse() *http.Response {
	header := make(http.Header)
	header.Set("server", "frp/"+version.Full())
	header.Set("Content-Type", "text/html")

	content := getNotFoundPageContent()
	res := &http.Response{
		Status:        "Not Found",
		StatusCode:    404,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        header,
		Body:          io.NopCloser(bytes.NewReader(content)),
		ContentLength: int64(len(content)),
	}
	return res
}
