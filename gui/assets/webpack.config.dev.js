const webpack = require('webpack')

module.exports = {
    mode: 'development',
    entry: './src/ts/Main.tsx',
    output: {
        path: __dirname + '/build',
        publicPath: '/',
        filename: 'app.js'
    },
    devtool: 'source-map',
    devServer: {
        contentBase: __dirname + '/build',
        compress: false,
        port: 9000,
        hot: true,
    },
    module: {
        rules: [
            {
                test: /\.(js|jsx|ts|tsx)$/i,
                exclude: /node_modules/,
                loader: 'babel-loader'
            },
            {
                test: /\.svg$/i,
                use: [
                    'svg-inline-loader',
                ]
            },
            {
                test: /\.(png|jpg|gif|woff)$/i,
                loader: 'url-loader'
            },
            {
                test: /\.less$/i,
                use: [
                    {
                        loader: 'style-loader'
                    },
                    {
                        loader: 'css-modules-typescript-loader',
                        options: {
                            mode: process.env.CI ? 'verify' : 'emit'
                        }
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
                        loader: 'postcss-loader'
                    },
                    {
                        loader: 'less-loader'
                    }
                ],
            },
        ]
    },
    plugins: [
        // new CleanWebpackPlugin(),
        new webpack.HotModuleReplacementPlugin(),
    ],
    resolve: {
        extensions: ['*', '.js', '.jsx', '.ts', '.tsx']
    },
}
