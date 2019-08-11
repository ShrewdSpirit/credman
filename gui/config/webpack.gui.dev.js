const webpack = require('webpack')
const MiniCssExtract = require('mini-css-extract-plugin')
const { CleanWebpackPlugin } = require('clean-webpack-plugin')
const Html = require('html-webpack-plugin')
const HtmlRoot = require('html-webpack-root-plugin')
const { join } = require('path')

module.exports = {
    mode: 'development',
    entry: join(__dirname, '..', 'src', 'ts', 'Main.tsx'),
    output: {
        path: join(__dirname, '..', 'dist'),
        publicPath: '/',
        filename: 'app.js'
    },
    resolve: {
        extensions: ['*', '.js', '.jsx', '.ts', '.tsx']
    },
    devtool: 'source-map',
    devServer: {
        contentBase: join(__dirname, '..', 'dist'),
        publicPath: '/',
        compress: false,
        port: 9000,
        hot: true,
    },
    plugins: [
        new MiniCssExtract({
            filename: 'app.css'
        }),
        new CleanWebpackPlugin(),
        new Html({
            filename: join(__dirname, '..', 'dist', 'index.html'),
            favicon: join(__dirname, '..', 'favicon.ico'),
            title: 'Credman'
        }),
        new HtmlRoot(),
        new webpack.HotModuleReplacementPlugin(),
        new webpack.DefinePlugin({
            'development': JSON.stringify(true)
        })
    ],
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
                        loader: MiniCssExtract.loader,
                        options: {
                            hmr: true
                        }
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
                                localIdentName: "[name]-[local]"
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
}
