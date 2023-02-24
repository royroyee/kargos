const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function (app) {
    const APIServerAddress = process.env.REACT_APP_API_SERVER_ADDR;
    const APIServerPort = process.env.REACT_APP_API_SERVER_PORT;

    app.use(
        '/api',
        createProxyMiddleware({
            target: 'http://' + APIServerAddress + ':' + APIServerPort,
            changeOrigin: true,
        })
    );
};