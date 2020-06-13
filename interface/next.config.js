module.exports = {
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
