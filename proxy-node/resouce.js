const express = require('express');
const bodyParser = require('body-parser');

const app = express();
const port = 8080;

// 使用body-parser中间件解析JSON请求体
app.use(bodyParser.json());

// 允许跨域请求
app.use((req, res, next) => {
    res.setHeader('Access-Control-Allow-Origin', '*');
    res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
    res.setHeader('Access-Control-Allow-Headers', 'Content-Type');
    next();
});

// 处理OPTIONS请求（预检请求）
app.options('/api/wechat', (req, res) => {
    res.sendStatus(200);
});

// 处理POST请求
app.post('/api/wechat', (req, res) => {
    const type = req.query.type;
    const requestBody = req.body;

    console.log(`收到类型为 ${type} 的请求`);
    console.log('请求体:', requestBody);

    // 根据type参数处理不同类型的请求
    let responseData;
    switch (type) {
        case '1':
            responseData = {
                code: 200,
                message: '成功处理类型1请求',
                data: {
                    result: '这是类型1的处理结果',
                    timestamp: new Date().toISOString()
                }
            };
            break;
        default:
            responseData = {
                code: 400,
                message: `不支持的请求类型: ${type}`,
                data: null
            };
            res.status(400);
    }

    // 返回JSON响应
    res.json(responseData);
});

// 启动服务器
app.listen(port, () => {
    console.log(`服务器运行在 http://localhost:${port}`);
});    