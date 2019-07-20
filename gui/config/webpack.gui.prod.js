const Terser = require('terser-webpack-plugin')
const MiniCssExtract = require('mini-css-extract-plugin')
const OptimizeCssAssets = require('optimize-css-assets-webpack-plugin')
const { join } = require('path')

module.exports = {
    mode: 'production',
    entry: join(__dirname, '..', 'src', 'ts', 'main.tsx'),
    output: {
        path: join(__dirname, '..', 'dist'),
        publicPath: '/',
        filename: 'app.js'
    },
    resolve: {
        extensions: ['*', '.js', '.jsx', '.ts', '.tsx']
    },
    optimization: {
        minimizer: [new Terser({}), new OptimizeCssAssets({})],
        splitChunks: {
            cacheGroups: {
              styles: {
                test: /\.css$/i,
                chunks: 'all',
                enforce: true,
              },
            },
          }
    },
    plugins: [
        new MiniCssExtract({
            filename: 'app.css'
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
                    MiniCssExtract.loader,
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
    }
}
