import request from '@/api/request'

export default {
    openSystemProxy() {
        return request({
            url: 'api/proxy-open',
            method: 'post',
        })
    },
    unsetSystemProxy() {
        return request({
            url: 'api/proxy-unset',
            method: 'post',
        })
    },
    openDirectoryDialog() {
        return request({
            url: 'api/open-directory',
            method: 'get'
        })
    },
    openFileDialog() {
        return request({
            url: 'api/open-file',
            method: 'get'
        })
    },
    openFolder(data: object) {
        return request({
            url: 'api/open-folder',
            method: 'post',
            data: data
        })
    },
    isProxy() {
        return request({
            url: 'api/is-proxy',
            method: 'get'
        })
    },
    appInfo() {
        return request({
            url: 'api/app-info',
            method: 'get',
        })
    },
    getConfig() {
        return request({
            url: 'api/get-config',
            method: 'get',
        })
    },
    setConfig(data: object) {
        return request({
            url: 'api/set-config',
            method: 'post',
            data: data
        })
    },
    setType(data: string[]) {
        return request({
            url: 'api/set-type',
            method: 'post',
            data: {
                type: data.toString()
            }
        })
    },
    clear() {
        return request({
            url: 'api/clear',
            method: 'post'
        })
    },
    delete(data: object) {
        return request({
            url: 'api/delete',
            method: 'post',
            data: data
        })
    },
    download(data: object) {
        return request({
            url: 'api/download',
            method: 'post',
            data: data
        })
    },
    wxFileDecode(data: object) {
        return request({
            url: 'api/wx-file-decode',
            method: 'post',
            data: data
        })
    },
}