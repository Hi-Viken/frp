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
	NotFound = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>404 · Page Not Found</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
    * { margin: 0; padding: 0; box-sizing: border-box; }

    :root {
        --bg: #0d1117;
        --fg: #e6edf3;
        --muted: #8b949e;
        --accent: #58a6ff;
        --accent2: #bc8cff;
    }

    html, body { height: 100%; }

    body {
        background: var(--bg);
        color: var(--fg);
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
                     "Helvetica Neue", Arial, "PingFang SC", "Microsoft YaHei", sans-serif;
        display: flex;
        align-items: center;
        justify-content: center;
        overflow: hidden;
        position: relative;
    }

    .glow {
        position: fixed;
        border-radius: 50%;
        filter: blur(120px);
        opacity: .22;
        pointer-events: none;
    }
    .glow-a {
        width: 480px; height: 480px;
        background: var(--accent);
        top: -160px; left: -120px;
        animation: drift 14s ease-in-out infinite alternate;
    }
    .glow-b {
        width: 420px; height: 420px;
        background: var(--accent2);
        bottom: -140px; right: -100px;
        animation: drift 18s ease-in-out infinite alternate-reverse;
    }
    @keyframes drift {
        from { transform: translate(0, 0) scale(1); }
        to   { transform: translate(60px, 40px) scale(1.15); }
    }

    .card {
        position: relative;
        text-align: center;
        padding: 48px 40px;
        max-width: 560px;
        width: calc(100% - 48px);
        background: rgba(22, 27, 34, .55);
        border: 1px solid rgba(240, 246, 252, .1);
        border-radius: 20px;
        backdrop-filter: blur(14px);
        -webkit-backdrop-filter: blur(14px);
        box-shadow: 0 24px 64px rgba(0, 0, 0, .45);
        animation: rise .7s cubic-bezier(.22, 1, .36, 1) both;
    }
    @keyframes rise {
        from { opacity: 0; transform: translateY(24px); }
        to   { opacity: 1; transform: translateY(0); }
    }

    .code {
        font-size: clamp(88px, 18vw, 132px);
        font-weight: 800;
        line-height: 1;
        letter-spacing: -4px;
        background: linear-gradient(120deg, var(--accent), var(--accent2));
        -webkit-background-clip: text;
        background-clip: text;
        color: transparent;
        text-shadow: 0 0 60px rgba(88, 166, 255, .25);
    }

    .card h1 {
        margin-top: 18px;
        font-size: 22px;
        font-weight: 600;
        letter-spacing: .5px;
    }

    .card p {
        margin-top: 12px;
        color: var(--muted);
        font-size: 15px;
        line-height: 1.7;
    }

    .card p em { color: var(--fg); font-style: normal; }

    .divider {
        margin: 28px auto;
        width: 72px;
        height: 2px;
        border-radius: 2px;
        background: linear-gradient(90deg, transparent, var(--accent), transparent);
    }

    .footer {
        font-size: 13px;
        color: var(--muted);
    }
    .footer a {
        color: var(--accent);
        text-decoration: none;
        border-bottom: 1px solid rgba(88, 166, 255, .35);
        transition: border-color .2s;
    }
    .footer a:hover { border-bottom-color: var(--accent); }

    .sig {
        margin-top: 10px;
        font-size: 12px;
        color: rgba(139, 148, 158, .7);
    }

    @media (prefers-reduced-motion: reduce) {
        .glow, .card { animation: none; }
    }
</style>
</head>
<body>
    <div class="glow glow-a"></div>
    <div class="glow glow-b"></div>

    <main class="card">
        <div class="code">404</div>
        <h1>Page Not Found</h1>
        <p>Sorry, the page you are looking for is currently unavailable.<br>
        Please try again later.</p>
        <div class="divider"></div>
        <p class="footer">The server is powered by <a href="https://github.com/fatedier/frp">frp</a>.</p>
        <p class="sig"><em>Faithfully yours, frp.</em></p>
    </main>
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
