package core

import (
	"embed"
	"github.com/getlantern/systray"
	"github.com/vrischmann/userdir"
	"os"
	"path/filepath"
	// sysRuntime "runtime"
	// "strings"
	"time"
)

type App struct {
	assets      embed.FS
	AppName     string            `json:"AppName"`
	Version     string            `json:"Version"`
	Description string            `json:"Description"`
	Copyright   string            `json:"Copyright"`
	UserDir     string            `json:"-"`
	LockFile    string            `json:"-"`
	PublicCrt   []byte            `json:"-"`
	PrivateKey  []byte            `json:"-"`
	IsProxy     bool              `json:"-"`
	MenuQuit    *systray.MenuItem `json:"-"`
	MenuProxy   *systray.MenuItem `json:"-"`
	MenuOpen    *systray.MenuItem `json:"-"`
}

var (
	appOnce        *App
	globalConfig   *Config
	globalLogger   *Logger
	resourceOnce   *Resource
	systemOnce     *SystemSetup
	proxyOnce      *Proxy
	httpServerOnce *HttpServer
)

func GetApp(assets embed.FS) *App {
	if appOnce == nil {
		appOnce = &App{
			assets:      assets,
			AppName:     "resd-mini",
			Version:     "1.0.3",
			Description: "resd-mini是一款集网络资源嗅探 + 高速下载功能于一体的软件，高颜值、高性能和多样化，提供个人用户下载自己上传到各大平台的网络资源功能！",
			Copyright:   "Copyright © 2025",
			PublicCrt: []byte(`
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
`),
			PrivateKey: []byte(`
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
`),
		}
		appOnce.UserDir = filepath.Join(userdir.GetConfigHome(), appOnce.AppName)
		LogWithLine("UserDir：", appOnce.UserDir)
		appOnce.LockFile = filepath.Join(appOnce.UserDir, "install.lock")
		initLogger()
		initConfig()
		initProxy()
		initResource()
		initHttpServer()
		initSystem()
	}
	return appOnce
}

func (a *App) Startup() {
	systray.Run(a.OnReady, a.OnExit)
}

func (a *App) OnExit() {
	a.UnsetSystemProxy()
	globalLogger.Close()
}

func (a *App) OnReady() {
	systray.SetIcon(a.getIcon())
	systray.SetTitle("")
	systray.SetTooltip(a.Description)

	a.MenuProxy = systray.AddMenuItem("启动代理", "设置系统代理")
	a.MenuOpen = systray.AddMenuItem("打开面板", "打开管理面板")
	a.MenuQuit = systray.AddMenuItem("退出", "退出应用")

	go httpServerOnce.run()

	time.AfterFunc(200*time.Millisecond, func() {
		// if globalConfig.AutoProxy {
		// 	appOnce.OpenSystemProxy()
		// }
		// _ = OpenBrowser("http://127.0.0.1:" + globalConfig.Port)
	})

	go func() {
		for {
			select {
			case <-a.MenuOpen.ClickedCh:
				_ = OpenBrowser("http://127.0.0.1:" + globalConfig.Port)
			case <-a.MenuProxy.ClickedCh:
				if appOnce.IsProxy {
					a.UnsetSystemProxy()
				} else {
					a.OpenSystemProxy()
				}
				httpServerOnce.send("updateProxyStatus", map[string]interface{}{
					"isProxy": appOnce.IsProxy,
				})
			case <-a.MenuQuit.ClickedCh:
				systray.Quit()
				os.Exit(0)
			}
		}
	}()

	go func() {
		if a.isInstall() {
			return
		}
		err := os.MkdirAll(a.UserDir, os.ModePerm)
		if err != nil {
			return
		}
		a.installCert()
	}()
}

func (a *App) getIcon() []byte {
	file, err := a.assets.ReadFile("web/dist/favicon.ico")
	if err != nil {
		return []byte{}
	}
	return file
}

func (a *App) installCert() {
	// if res, err := systemOnce.installCert(); err != nil {
	// 	DialogErr("证书安装失败，err:" + err.Error() + ", " + res)
	// 	if sysRuntime.GOOS == "darwin" {
	// 		LogWithLine(`证书安装失败，请执行如下命令，安装证书文件:`)
	// 		LogWithLine(`echo "输入本地登录密码" && sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain "`)
	// 	} else if sysRuntime.GOOS == "windows" && strings.Contains(err.Error(), "Access is denied.") {
	// 		LogWithLine(`证书安装失败，首次启用本软件，请使用鼠标右键选择以管理员身份运行`)
	// 	} else if sysRuntime.GOOS == "linux" && strings.Contains(err.Error(), "Access is denied.") {
	// 		LogWithLine("证书路径: " + systemOnce.CertFile + ", 请手动安装，安装完成后请执行: touch" + a.LockFile)
	// 	} else {
	// 		LogWithLine(`证书安装失败，` + sysRuntime.GOOS)
	// 	}
	// } else {
	// 	if err := a.lock(); err != nil {
	// 		globalLogger.err(err)
	// 	}
	// }
}

func (a *App) OpenSystemProxy() bool {
	defer func() {
		if a.IsProxy {
			a.MenuProxy.SetTitle("关闭代理")
		} else {
			a.MenuProxy.SetTitle("开启代理")
		}
	}()
	if a.IsProxy {
		return true
	}
	err := systemOnce.setProxy()
	if err == nil {
		a.IsProxy = true
		return true
	}
	DialogErr("设置失败" + err.Error())
	return false
}

func (a *App) UnsetSystemProxy() bool {
	defer func() {
		if a.IsProxy {
			a.MenuProxy.SetTitle("关闭代理")
		} else {
			a.MenuProxy.SetTitle("开启代理")
		}
	}()

	if !a.IsProxy {
		return true
	}
	err := systemOnce.unsetProxy()
	if err == nil {
		a.IsProxy = false
		return true
	}
	DialogErr("设置失败")
	return false
}

func (a *App) isInstall() bool {
	return FileExist(a.LockFile)
}

func (a *App) lock() error {
	err := os.WriteFile(a.LockFile, []byte("success"), 0777)
	if err != nil {
		return err
	}
	return nil
}
