module.exports = {
    mode: 'development',
    entry: './app.js',
    output: {
        path: __dirname + '/build',
        publicPath: '/',
        filename: 'app.js'
    },
    module: {
        rules: [
            {
                test: /\.(js|jsx)$/,
                exclude: /node_modules/,
                use: ['babel-loader']
            },
            {
                test: /\.css$/i,
                use: [
                    {
                        loader: 'style-loader'
                    },
                    {
                        loader: 'css-loader',
                        options: {
                            modules: {
                                mode: 'local',
                                localIdentName: "[local]_[hash:base64:5]"
                            }
                        }
                    },
                    {
                        loader: 'less-loader'
                    }
                ],
            },
        ]
    },
    resolve: {
        extensions: ['*', '.js', '.jsx']
    },
}
