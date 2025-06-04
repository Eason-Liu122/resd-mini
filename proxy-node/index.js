const http = require('http');
const https = require('https');
const httpProxy = require('http-proxy');
const fs = require('fs');
const path = require('path');
const crypto = require('crypto');
const url = require('url');

class Proxy {
  constructor() {
    this.proxy = httpProxy.createProxyServer({
      secure: false,
      xfwd: true
    });
    this.init();
  }

  init() {
    this.setCa();
    this.setupProxyEvents();
    this.startServer();
  }

  setCa() {
    // 加载证书和私钥
    try {
      this.caCert = `
-----BEGIN CERTIFICATE-----
MIIDwzCCAqugAwIBAgIUFAnC6268dp/z1DR9E1UepiWgWzkwDQYJKoZIhvcNAQEL
BQAwcDELMAkGA1UEBhMCQ04xEjAQBgNVBAgMCUNob25ncWluZzESMBAGA1UEBwwJ
Q2hvbmdxaW5nMQ4wDAYDVQQKDAVnb3dhczEWMBQGA1UECwwNSVQgRGVwYXJ0bWVu
dDERMA8GA1UEAwwIZ293YXMuY24wIBcNMjQwMjE4MDIwOTI2WhgPMjEyNDAxMjUw
MjA5MjZaMHAxCzAJBgNVBAYTAkNOMRIwEAYDVQQIDAlDaG9uZ3FpbmcxEjAQBgNV
BAcMCUNob25ncWluZzEOMAwGA1UECgwFZ293YXMxFjAUBgNVBAsMDUlUIERlcGFy
dG1lbnQxETAPBgNVBAMMCGdvd2FzLmNuMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEA3A7dt7eoqAaBxv2Npjo8Z7VkGvXT93jZfpgAuuNuQ5RLcnOnMzQC
CrrjPcLfsAMA0AIK3eUWsXXKSR9SZTJBLQRZCJHZ9AIPfA+58JVQPTjd8UIuQZJf
rDf6FjhPJTsLzcjTU+mT7t6lEimPEl2VWN9eXWqs9nkVrJtqLao6m1hoYfXOxRh6
96/WgBtPHcmjujryteBiSITVflDjx+YQzDGsbqw7fM52klMPd2+w/vmhJ4pxq6P7
Ni2OBvdXYDPIuLfPFFqG16arORjBkyNCJy19iOuh5LXh+EUX11wvbLwNgsTd8j9v
eBSD+4HUUNQhiXiXJbs7I7cdFYthvb609QIDAQABo1MwUTAdBgNVHQ4EFgQUdI8p
aY1A47rWCRvQKSTRCCk6FoMwHwYDVR0jBBgwFoAUdI8paY1A47rWCRvQKSTRCCk6
FoMwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEArMCAfqidgXL7
cW5TAZTCqnUeKzbbqMJgk6iFsma8scMRsUXz9ZhF0UVf98376KvoJpy4vd81afbi
TehQ8wVBuKTtkHeh/MkXMWC/FU4HqSjtvxpic2+Or5dMjIrfa5VYPgzfqNaBIUh4
InD5lo8b/n5V+jdwX7RX9VYAKug6QZlCg5YSKIvgNRChb36JmrGcvsp5R0Vejnii
e3oowvgwikqm6XR6BEcRpPkztqcKST7jPFGHiXWsAqiibc+/plMW9qebhfMXEGhQ
5yVNeSxX2zqasZvP/fRy+3I5iVilxtKvJuVpPZ0UZzGS0CJ/lF67ntibktiPa3sR
D8HixYbEDg==
-----END CERTIFICATE-----
`
      this.caKey = `
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDcDt23t6ioBoHG
/Y2mOjxntWQa9dP3eNl+mAC6425DlEtyc6czNAIKuuM9wt+wAwDQAgrd5RaxdcpJ
H1JlMkEtBFkIkdn0Ag98D7nwlVA9ON3xQi5Bkl+sN/oWOE8lOwvNyNNT6ZPu3qUS
KY8SXZVY315daqz2eRWsm2otqjqbWGhh9c7FGHr3r9aAG08dyaO6OvK14GJIhNV+
UOPH5hDMMaxurDt8znaSUw93b7D++aEninGro/s2LY4G91dgM8i4t88UWobXpqs5
GMGTI0InLX2I66HkteH4RRfXXC9svA2CxN3yP294FIP7gdRQ1CGJeJcluzsjtx0V
i2G9vrT1AgMBAAECggEAF0obfQ4a82183qqHC0iui+tOpOvPeyl3G0bLDPx09wIC
2iITV//xF2GgGzE8q0wmEd2leMZ+GFn3BrYh6kPfUfxbz+RfxMtTCDZB34xt6YzT
MG1op9ft+DQUa7WZ6r7NCQJwGzllRqqZncp4MeFlpPo+6nQXyh4WhSYNnredbENE
uPZ63Kme4RZfMvtVso+XgAQM3oDih0onv1YitmNQpL9rRzlthTfybAT4737DBINq
zsmBNE6QIsXnSKpzo11OtDgof2QM9ac6eAXf73oTpDxfodwCotILytKn+8WYvlR+
T15uuknb4M3XI1FPVolkF4qtK5SLAAbVzV4DsCmuIQKBgQD6bTKKbL2huvU6dEKx
bgS079LfQUxxOTClgwkhVsMxRtvcPBnHYMAsPK4mnMhEh9x+TF6wxMx0pmhQluPI
ZULNBj/qdoiBL0RwVLA+9jgE0NeWB3XXFDsEavQBr9Q8CC0uzrsgsxFcvHpqqs2Q
RtngxRWtJP06D6mKC23s4YjDHwKBgQDg9KUCFqOmWcRXyeg9gYMC4jFFQw4lUQBd
sYpqSMHDw1b+T1W/dCPbwbxZL/+d8y930BYy9QYDtQwHdLyXCH0pHM7S6rfgr5xk
2Szd8xBUIqmeV/zcR00mTeQHJ1M50VHfclAVgZgkpWSoLwbX+bXyx/mfqLAtynZ5
yU9RfrT5awKBgQC0uJ8TlFvZXjFgyMvkfY/5/2R/ZwFCaFI573FkVNeyNP+vVNQJ
tUGZ6wSGqvg/tIgjwPtIuA0QVZLMLcgeMy1dBhiUHIxwJetO4V77YPaWSxx5kdKx
r1DT5FdI7FnOJNxufhQ/CdsKwJ3bYn3Mk8TiV3hIJnx0LR9dltfybeQjYwKBgDOY
6aApATBOtrJMJXC2HA61QwfX8Y6tnZ/f8RefyJHWZEXAfLKFORRWw5TRZZgdB247
1Furx81h4Xh0Vi1uTQb5DJdkLvjiTsTy60+dSMmDidQ/6ke8Mv3uL7dUVcqVMGpI
FgZYy0TcitHot3EiXZFqPN9aGc7m+XXFruPKZEgxAoGBAMA96jsow7CzulU+GRW8
Njg4zWuAEVErgPoNBcOXAVWLCTU/qGIEMNpZL6Ok34kf13pJDMjQ8eDuQHu5CSqf
0ul5Zy85fwfVq2IvNAyYT8eflQprTejFw22CHhfPBfADVW9ro8dK/Jw+J/31Vh7V
ILKEQKmPPzKs7kp/7Nz+2cT3
-----END PRIVATE KEY-----
`
    } catch (err) {
      console.error('启动代理服务失败：', err.message);
      return;
    }
  }

  setupProxyEvents() {
    // 请求拦截
    this.proxy.on('proxyReq', (proxyReq, req, res) => {
      this.httpRequestEvent(proxyReq, req, res);
    });

    // 响应拦截
    this.proxy.on('proxyRes', (proxyRes, req, res) => {
      console.log('proxyRes....', proxyRes)
      console.log('proxyRes....', req)
      console.log('proxyRes....', res)
      this.httpResponseEvent(proxyRes, req, res);
    });

    // 错误处理
    this.proxy.on('error', (err, req, res) => {
      console.error('代理错误:', err.message);
      res.writeHead(500, { 'Content-Type': 'text/plain' });
      res.end('代理服务器错误');
    });
  }

  startServer() {
    // HTTP代理服务器
    this.httpServer = http.createServer((req, res) => {
      this.proxy.web(req, res, { target: `http://${req.headers.host}` });
    });

    // HTTPS代理服务器
    this.httpsServer = https.createServer({
      key: this.caKey,
      cert: this.caCert
    }, (req, res) => {
      this.proxy.web(req, res, { target: `https://${req.headers.host}` });
    });

    // 处理CONNECT请求（HTTPS协议）
    this.httpServer.on('connect', (req, socket, head) => {
      this.handleConnect(req, socket, head);
    });

    // 启动服务器
    this.httpServer.listen(8080, () => {
      console.log('HTTP代理服务器已启动，监听端口8080');
    });

    this.httpsServer.listen(8443, () => {
      console.log('HTTPS代理服务器已启动，监听端口8443');
    });
  }

  handleConnect(req, socket, head) {
    const targetUrl = url.parse(`https://${req.url}`);
    const options = {
      port: 443,
      host: targetUrl.hostname,
      method: 'CONNECT',
      path: targetUrl.path
    };

    const proxySocket = net.connect(options, () => {
      socket.write('HTTP/1.1 200 Connection established\r\n\r\n');
      proxySocket.write(head);
      proxySocket.pipe(socket);
      socket.pipe(proxySocket);
    });

    proxySocket.on('error', (err) => {
      console.error('CONNECT请求错误:', err.message);
      socket.end(`HTTP/1.1 500 Connection error\r\n\r\n`);
    });
  }

  httpRequestEvent(proxyReq, req, res) {
    const host = req.headers.host || '';
    const path = req.url || '';

    if (host.includes('res-downloader.666666.com') && path.includes('/wechat')) {
      const query = new url.URLSearchParams(url.parse(req.url).query);
      if (globalConfig.WxAction && query.get('type') === '1') {
        return this.handleWechatRequest(req, res);
      } else if (!globalConfig.WxAction && query.get('type') === '2') {
        return this.handleWechatRequest(req, res);
      } else {
        return this.buildEmptyResponse(res);
      }
    }
  }

  handleWechatRequest(req, res) {
    // 处理微信请求的逻辑
    console.log('处理微信请求:', req.url);
    
    // 这里需要根据具体需求实现
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ status: 'success', message: '微信请求已处理' }));
  }

  buildEmptyResponse(res) {
    res.writeHead(200, { 'Content-Type': 'text/plain' });
    res.end('内容不存在');
  }

  httpResponseEvent(proxyRes, req, res) {
    if (!proxyRes || !req || (proxyRes.statusCode !== 200 && proxyRes.statusCode !== 206)) {
      return;
    }

    const host = req.headers.host || '';
    const path = req.url || '';

    let body = [];
    proxyRes.on('data', (chunk) => body.push(chunk));
    proxyRes.on('end', () => {
      body = Buffer.concat(body).toString();
      console.log(body)
      
      if (host.endsWith('channels.weixin.qq.com') && 
          (path.includes('/web/pages/feed') || path.includes('/web/pages/home'))) {
        body = this.replaceWxJsContent(body, '.js"', `.js?v=${this.v()}"`);
      }

      if (host.endsWith('res.wx.qq.com')) {
        if (path.endsWith(`.js?v=${this.v()}`)) {
          body = this.replaceWxJsContent(body, '.js"', `.js?v=${this.v()}"`);
        }

        if (path.includes('web/web-finder/res/js/virtual_svg-icons-register.publish')) {
          body = body.replace(/get\s*media\(\)\{/, `
            get media(){
              if(this.objectDesc){
                fetch("https://res-downloader.666666.com/wechat?type=1", {
                  method: "POST",
                  mode: "no-cors",
                  body: JSON.stringify(this.objectDesc),
                });
              };
          `);
        }
      }

      // 修改响应头和内容
      proxyRes.headers['content-length'] = Buffer.byteLength(body);
      res.writeHead(proxyRes.statusCode, proxyRes.headers);
      res.end(body);
    });
  }

  replaceWxJsContent(body, old, new_) {
    // 确保替换是全局的且安全处理正则表达式特殊字符
    const escapeRegExp = (string) => {
      return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    };
    
    const regex = new RegExp(escapeRegExp(old), 'g');
    return body.replace(regex, new_);
  }

  v() {
    return appOnce.Version;
  }
}

// 单例模式
let proxyOnce = null;

function initProxy() {
  if (!proxyOnce) {
    proxyOnce = new Proxy();
  }
  return proxyOnce;
}

// 初始化代理
initProxy();

// 全局配置示例
const globalConfig = {
  WxAction: true,
  UpstreamProxy: '',
  OpenProxy: false,
  Port: 8080
};

// 应用实例示例
const appOnce = {
  Version: '1.0.0',
  PublicCrt: path.join(__dirname, 'public.crt'),
  PrivateKey: path.join(__dirname, 'private.key')
};    