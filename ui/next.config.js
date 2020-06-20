module.exports = {
    publicRuntimeConfig: {
        apiUrl: process.env.API_URL,
        channelUrl: process.env.CHANNEL_URL,
        environment: process.env.NODE_ENV,
    },
    webpackDevMiddleware: config => {
        config.watchOptions = {
            poll: 1000,
            aggregateTimeout: 300,
            ignored: [
                /\.git\//,
                /\.next\//,
                /node_modules/,
            ]
        }
        return config
    },
}
