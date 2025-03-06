; 定义安装程序的名称和输出路径
OutFile "bin\resd-mini-x64-nsis.exe"

; 定义安装程序的名称
Name "Resd Mini"

; 定义默认安装目录
InstallDir "$PROGRAMFILES\Resd Mini"

Icon "windows\icon.ico"  ; 设置安装程序的图标

; 请求管理员权限
RequestExecutionLevel admin

; 开始安装部分
Section "Install"

  ; 创建安装目录
  SetOutPath "$INSTDIR"

  ; 将可执行文件复制到安装目录
  File "bin\resd-mini-x64.exe"
  File "windows\icon.ico"

  ; 创建开始菜单快捷方式
  CreateDirectory "$SMPROGRAMS\Resd Mini"
  CreateShortcut "$SMPROGRAMS\Resd Mini\Resd Mini.lnk" "$INSTDIR\resd-mini-x64.exe"
  CreateShortcut "$DESKTOP\Resd Mini.lnk" "$INSTDIR\resd-mini-x64.exe" "" "$INSTDIR\icon.ico"

  ; 写入卸载信息
  WriteUninstaller "$INSTDIR\Uninstall.exe"

  ; 创建开始菜单中的卸载快捷方式
  CreateShortcut "$SMPROGRAMS\Resd Mini\Uninstall Resd Mini.lnk" "$INSTDIR\Uninstall.exe"

SectionEnd

; 开始卸载部分
Section "Uninstall"

  ; 删除安装目录中的所有文件
  Delete "$INSTDIR\resd-mini-x64.exe"
  Delete "$INSTDIR\Uninstall.exe"

  ; 删除开始菜单快捷方式
  Delete "$SMPROGRAMS\Resd Mini\Resd Mini.lnk"
  Delete "$SMPROGRAMS\Resd Mini\Uninstall Resd Mini.lnk"
  RMDir "$SMPROGRAMS\Resd Mini"

  ; 删除桌面快捷方式
  Delete "$DESKTOP\Resd Mini.lnk"

  ; 删除安装目录
  RMDir "$INSTDIR"

SectionEnd