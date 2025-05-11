package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/ncruces/zenity"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	sysRuntime "runtime"
	"strings"
	"sync"
	"time"
)

type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type HttpServer struct {
	indexHTML []byte
	upGrader  websocket.Upgrader
	wsClients map[*websocket.Conn]bool
	broadcast chan []byte
	mutex     sync.RWMutex
}

func initHttpServer() *HttpServer {
	if httpServerOnce == nil {
		httpServerOnce = &HttpServer{
			upGrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
			wsClients: make(map[*websocket.Conn]bool),
			broadcast: make(chan []byte, 1000),
		}
		file, err := appOnce.assets.ReadFile("web/dist/index.html")
		if err != nil {
			globalLogger.Error().Stack().Err(err)
		} else {
			httpServerOnce.indexHTML = file
		}
	}
	return httpServerOnce
}

func (h *HttpServer) run() {
	listener, err := net.Listen("tcp", globalConfig.Host+":"+globalConfig.Port)
	if err != nil {
		log.Fatalf("无法启动监听: %v", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/ws", h.wsHandler)
	mux.HandleFunc("/api/preview", h.preview)
	mux.HandleFunc("/api/proxy-open", h.openSystemProxy)
	mux.HandleFunc("/api/proxy-unset", h.unsetSystemProxy)
	mux.HandleFunc("/api/open-directory", h.openDirectoryDialog)
	mux.HandleFunc("/api/open-file", h.openFileDialog)
	mux.HandleFunc("/api/open-folder", h.openFolder)
	mux.HandleFunc("/api/is-proxy", h.isProxy)
	mux.HandleFunc("/api/app-info", h.appInfo)
	mux.HandleFunc("/api/set-config", h.setConfig)
	mux.HandleFunc("/api/get-config", h.getConfig)
	mux.HandleFunc("/api/set-type", h.setType)
	mux.HandleFunc("/api/clear", h.clear)
	mux.HandleFunc("/api/delete", h.delete)
	mux.HandleFunc("/api/download", h.download)
	mux.HandleFunc("/api/wx-file-decode", h.wxFileDecode)
	mux.HandleFunc("/api/cert", h.handleCert)
	mux.HandleFunc("/api/get-media-infos", h.getMediaInfos)

	// Static assets endpoint
	mux.HandleFunc("/", h.staticHandler)

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Host == globalConfig.Host+":"+globalConfig.Port || r.Host == "127.0.0.1:"+globalConfig.Port && strings.Contains(r.URL.Path, "/api") {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				mux.ServeHTTP(w, r)
			} else {
				proxyOnce.Proxy.ServeHTTP(w, r)
			}
		}),
	}
	go h.handleMessages()
	LogWithLine("服务已启动，监听 http://" + globalConfig.Host + ":" + globalConfig.Port)
	if err := server.Serve(listener); err != nil {
		fmt.Printf("服务器异常: %v", err)
	}
}

func (h *HttpServer) staticHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(h.indexHTML)
		return
	}

	filePath := strings.TrimPrefix(r.URL.Path, "/")
	file, err := appOnce.assets.ReadFile("web/dist/" + filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Serve the file with correct content type
	http.ServeContent(w, r, filePath, time.Time{}, strings.NewReader(string(file)))
}

func (h *HttpServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upGrader.Upgrade(w, r, nil)
	if err != nil {
		LogWithLine("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	h.mutex.Lock()
	h.wsClients[conn] = true
	h.mutex.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			LogWithLine("WebSocket read error:", err)
			break
		}
		h.broadcast <- message
	}
	h.mutex.Lock()
	delete(h.wsClients, conn)
	h.mutex.Unlock()
}

func (h *HttpServer) preview(w http.ResponseWriter, r *http.Request) {
	realURL := r.URL.Query().Get("url")
	if realURL == "" {
		http.Error(w, "Missing 'url' parameter", http.StatusBadRequest)
		return
	}
	realURL, _ = url.QueryUnescape(realURL)
	parsedURL, err := url.Parse(realURL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	request, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		http.Error(w, "Failed to fetch the resource", http.StatusInternalServerError)
		return
	}

	if rangeHeader := r.Header.Get("Range"); rangeHeader != "" {
		request.Header.Set("Range", rangeHeader)
	}

	//request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
	//request.Header.Set("Referer", parsedURL.Scheme+"://"+parsedURL.Host+"/")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		http.Error(w, "Failed to fetch the resource", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)

	if contentRange := resp.Header.Get("Content-Range"); contentRange != "" {
		w.Header().Set("Content-Range", contentRange)
	}

	buffer := make([]byte, 32*1024) // 32KB buffer
	_, err = io.CopyBuffer(w, resp.Body, buffer)
	if err != nil {
		http.Error(w, "Failed to serve the resource", http.StatusInternalServerError)
	}
	return
}

func (h *HttpServer) handleCert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-x509-ca-data")
	w.Header().Set("Content-Disposition", "attachment;filename=res-downloader-public.crt")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(appOnce.PublicCrt)))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, io.NopCloser(bytes.NewReader(appOnce.PublicCrt)))
}

func (h *HttpServer) getMediaInfos(w http.ResponseWriter, r *http.Request) {
    // 假设resourceOnce有一个获取所有mediaInfo的方法GetAllMediaInfos()
    mediaInfos := resourceOnce.getAllMediaInfos()
    
    h.writeJson(w, ResponseData{
        Code: 1,
        Data: mediaInfos,
    })
}

func (h *HttpServer) handleMessages() {
	for {
		msg := <-h.broadcast
		h.mutex.RLock()
		for client := range h.wsClients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Printf("写入消息错误: %v", err)
				client.Close()
				h.mutex.Lock()
				delete(h.wsClients, client)
				h.mutex.Unlock()
			}
		}
		h.mutex.RUnlock()
	}
}

func (h *HttpServer) send(t string, data interface{}) {
	jsonData, err := json.Marshal(map[string]interface{}{
		"type": t,
		"data": data,
	})
	if err != nil {
		LogWithLine("Error converting map to JSON:", err)
		return
	}
	h.broadcast <- jsonData
}

func (h *HttpServer) writeJson(w http.ResponseWriter, data ResponseData) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		globalLogger.err(err)
	}
}

func (h *HttpServer) openDirectoryDialog(w http.ResponseWriter, r *http.Request) {
	folder, err := zenity.SelectFile(zenity.Filename(""), zenity.Directory())
	if err != nil {
		h.writeJson(w, ResponseData{Code: 0, Message: err.Error()})
		DialogErr("选择文件夹时出错: %v" + err.Error())
		return
	}

	h.writeJson(w, ResponseData{
		Code: 1,
		Data: map[string]interface{}{
			"folder": folder,
		},
	})
}

func (h *HttpServer) openFileDialog(w http.ResponseWriter, r *http.Request) {
	file, err := zenity.SelectFile(
		zenity.Filename(""),
		zenity.FileFilters{
			{"Video files", []string{"*.mp4"}, false},
		})
	if err != nil {
		h.writeJson(w, ResponseData{Code: 0, Message: err.Error()})
		DialogErr("选择文件夹时出错: %v" + err.Error())
		return
	}

	h.writeJson(w, ResponseData{
		Code: 1,
		Data: map[string]interface{}{
			"file": file,
		},
	})
	return
}

func (h *HttpServer) openFolder(w http.ResponseWriter, r *http.Request) {
	var data struct {
		FilePath string `json:"filePath"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err == nil && data.FilePath == "" {
		return
	}

	filePath := data.FilePath
	var cmd *exec.Cmd

	switch sysRuntime.GOOS {
	case "darwin":
		// macOS
		cmd = exec.Command("open", "-R", filePath)
	case "windows":
		// Windows
		cmd = exec.Command("explorer", "/select,", filePath)
	case "linux":
		// linux
		// 尝试使用不同的文件管理器
		cmd = exec.Command("nautilus", filePath)
		if err := cmd.Start(); err != nil {
			cmd = exec.Command("thunar", filePath)
			if err := cmd.Start(); err != nil {
				cmd = exec.Command("dolphin", filePath)
				if err := cmd.Start(); err != nil {
					cmd = exec.Command("pcmanfm", filePath)
					if err := cmd.Start(); err != nil {
						globalLogger.err(err)
						h.writeJson(w, ResponseData{Code: 0, Message: err.Error()})
						return
					}
				}
			}
		}
	default:
		h.writeJson(w, ResponseData{Code: 0, Message: "unsupported platform"})
		return
	}

	err = cmd.Start()
	if err != nil {
		globalLogger.err(err)
		h.writeJson(w, ResponseData{Code: 0, Message: err.Error()})
		return
	}
	h.writeJson(w, ResponseData{Code: 1})
}

func (h *HttpServer) openSystemProxy(w http.ResponseWriter, r *http.Request) {
	appOnce.OpenSystemProxy()
	h.writeJson(w, ResponseData{
		Code: 1,
		Data: map[string]bool{
			"isProxy": appOnce.IsProxy,
		},
	})
}

func (h *HttpServer) unsetSystemProxy(w http.ResponseWriter, r *http.Request) {
	appOnce.UnsetSystemProxy()
	h.writeJson(w, ResponseData{
		Code: 1,
		Data: map[string]bool{
			"isProxy": appOnce.IsProxy,
		},
	})
}

func (h *HttpServer) isProxy(w http.ResponseWriter, r *http.Request) {
	h.writeJson(w, ResponseData{
		Code: 1,
		Data: map[string]interface{}{
			"isProxy": appOnce.IsProxy,
		},
	})
}

func (h *HttpServer) appInfo(w http.ResponseWriter, r *http.Request) {
	h.writeJson(w, ResponseData{
		Code: 1,
		Data: appOnce,
	})
}

func (h *HttpServer) getConfig(w http.ResponseWriter, r *http.Request) {
	h.writeJson(w, ResponseData{
		Code: 1,
		Data: globalConfig,
	})
}

func (h *HttpServer) setConfig(w http.ResponseWriter, r *http.Request) {
	var data Config
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.writeJson(w, ResponseData{Code: 0, Message: err.Error()})
		return
	}
	globalConfig.setConfig(data)
	h.writeJson(w, ResponseData{Code: 1})
}

func (h *HttpServer) setType(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Type string `json:"type"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err == nil {
		if data.Type != "" {
			resourceOnce.setResType(strings.Split(data.Type, ","))
		} else {
			resourceOnce.setResType([]string{})
		}
	}

	h.writeJson(w, ResponseData{Code: 1})
}

func (h *HttpServer) clear(w http.ResponseWriter, r *http.Request) {
	resourceOnce.clear()
	h.writeJson(w, ResponseData{Code: 1})
}

func (h *HttpServer) delete(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Sign string `json:"sign"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err == nil && data.Sign != "" {
		resourceOnce.delete(data.Sign)
	}
	h.writeJson(w, ResponseData{Code: 1})
}

func (h *HttpServer) download(w http.ResponseWriter, r *http.Request) {
	var data struct {
		MediaInfo
		DecodeStr string `json:"decodeStr"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.writeJson(w, ResponseData{Code: 0, Message: err.Error()})
		return
	}
	resourceOnce.download(data.MediaInfo, data.DecodeStr)
	h.writeJson(w, ResponseData{Code: 1})
}

func (h *HttpServer) wxFileDecode(w http.ResponseWriter, r *http.Request) {
	var data struct {
		MediaInfo
		Filename  string `json:"filename"`
		DecodeStr string `json:"decodeStr"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.writeJson(w, ResponseData{Code: 0, Message: err.Error()})
		return
	}
	savePath, err := resourceOnce.wxFileDecode(data.MediaInfo, data.Filename, data.DecodeStr)
	if err != nil {
		h.writeJson(w, ResponseData{Code: 0, Message: err.Error()})
		return
	}
	h.writeJson(w, ResponseData{
		Code: 1,
		Data: map[string]string{
			"save_path": savePath,
		},
	})
}
