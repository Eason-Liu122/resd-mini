package core

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ncruces/zenity"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
	"fmt"
)

func Empty(data interface{}) {
}

func LogWithLine(format string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1) // 获取调用者的文件路径和行号
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s:%d: %s\n", file, line, msg)
}

func DialogErr(message string) {
	err := zenity.Warning(message,
		zenity.Title("Warning"),
		zenity.WarningIcon)
	if err != nil {
		LogWithLine("DialogErr:", err)
	}
}

func IsDevelopment() bool {
	return os.Getenv("APP_ENV") == "development"
}

func FileExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0777)
	}
	return nil
}

func TypeSuffix(mime string) (string, string) {
	switch strings.ToLower(mime) {
	case "image/png",
		"image/webp",
		"image/jpeg",
		"image/jpg",
		"image/gif",
		"image/avif",
		"image/bmp",
		"image/tiff",
		"image/heic",
		"image/x-icon",
		"image/svg+xml",
		"image/vnd.adobe.photoshop":
		return "image", ".png"
	case "audio/mpeg",
		"audio/wav",
		"audio/aiff",
		"audio/x-aiff",
		"audio/aac",
		"audio/ogg",
		"audio/flac",
		"audio/midi",
		"audio/x-midi",
		"audio/x-ms-wma",
		"audio/opus",
		"audio/webm",
		"audio/mp4",
		"audio/mp3":
		return "audio", ".mp3"
	case "video/mp4",
		"video/webm",
		"video/ogg",
		"video/x-msvideo",
		"video/mpeg",
		"video/quicktime",
		"video/x-ms-wmv",
		"video/3gpp",
		"video/x-matroska":
		return "video", ".mp4"
	case "audio/video",
		"video/x-flv":
		return "live", ".mp4"
	case "application/vnd.apple.mpegurl",
		"application/x-mpegurl":
		return "m3u8", ".m3u8"
	case "application/pdf":
		return "pdf", ".pdf"
	case "application/vnd.ms-powerpoint",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation":
		return "ppt", ".ppt"
	case "application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return "xls", ".xls"
	case "application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return "doc", ".doc"

	}
	return "", ""
}

func BuildReferer(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return u.Scheme + "://" + u.Host + "/"
}

func Md5(data string) string {
	hashNew := md5.New()
	hashNew.Write([]byte(data))
	hash := hashNew.Sum(nil)
	return hex.EncodeToString(hash)
}

func FormatSize(size float64) string {
	if size > 1048576 {
		return fmt.Sprintf("%.2fMB", float64(size)/1048576)
	}
	if size > 1024 {
		return fmt.Sprintf("%.2fKB", float64(size)/1024)
	}
	return fmt.Sprintf("%.0fb", size)
}

func GetTopLevelDomain(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	host := parsedURL.Hostname()
	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return ""
	}
	return strings.Join(parts[len(parts)-2:], ".")
}

func GetCurrentDateTimeFormatted() string {
	now := time.Now()
	return fmt.Sprintf("%04d%02d%02d%02d%02d%02d",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second())
}

func OpenBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	default:
		return fmt.Errorf("unsupported platform")
	}

	return exec.Command(cmd, args...).Start()
}
