module.exports = {
    mode: 'production',
    entry: './src/ts/Main.tsx',
    output: {
        path: __dirname + '/build',
        publicPath: '/',
        filename: 'app.js'
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
    resolve: {
        extensions: ['*', '.js', '.jsx', '.ts', '.tsx']
    },
}
