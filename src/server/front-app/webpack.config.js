const { join } = require('path');
const webpack = require('webpack');
const { VueLoaderPlugin } = require('vue-loader');
// const { HotModuleReplacementPlugin } = require('webpack');
const HTMLWebpackPlugin = require('html-webpack-plugin');
const WorkboxPlugin = require('workbox-webpack-plugin');
const CopyPlugin = require("copy-webpack-plugin");
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const autoprefixer = require('autoprefixer');

module.exports = {
    mode: 'development',
    entry: {
        "app": join(__dirname,'js/app.js'),
        "styles": join(__dirname, 'css/styles.sass')
    },
    output: {
        path: join(__dirname, 'dist'),
        filename: '[name].min.js',
        clean: true,
    },
    module: {
        rules: [
            {
                test: /\.js$/,
                loader: 'babel-loader',
                options: {
                    presets: ['@babel/preset-env']
                }
            }, {
                test: /.vue$/,
                loader: 'vue-loader'
            },
            {
                test: /\.css$/,
                use: [
                    'postcss-loader',
                    process.env.NODE_ENV !== 'production'
                        ? 'vue-style-loader'
                        : MiniCssExtractPlugin.loader,
                    'css-loader',
                ]
            },
            {
                test: /\.s[ac]ss$/i,
                use: [
                    process.env.NODE_ENV !== 'production'
                        ? 'vue-style-loader'
                        : MiniCssExtractPlugin.loader,
                    'css-loader',
                    {
                        loader: 'sass-loader',
                        options: {
                            implementation: require("sass"),
                            sassOptions: {
                                indentedSyntax: true,
                            }
                        }
                    },
                    // 'postcss-loader',
                ],
            },
        ]
    },
    plugins: [
        new VueLoaderPlugin(),
        new webpack.DefinePlugin({
            // Drop Options API from bundle
            __VUE_OPTIONS_API__: true,
            __VUE_PROD_DEVTOOLS__: process.env.NODE_ENV !== 'production',
        }),
        new HTMLWebpackPlugin({
            showErrors: true,
            cache: true,
            template: join(__dirname, 'index.template.html'),
            filename: join(__dirname, 'index.html'),
        }),
        // new WorkboxPlugin.GenerateSW({
        //     clientsClaim: true,
        //     skipWaiting: true,
        // }),
        new MiniCssExtractPlugin({
            filename: '[name].css',
        }),
        // new CopyPlugin({
        //     patterns: [
        //         {
        //             from: join(__dirname, 'front/dist'),
        //             to: join(__dirname, 'static'),
        //         },
        //     ],
        // }),
        // new CopyPlugin({
        //     patterns: [
        //         {
        //             from: join(__dirname, 'front/images'),
        //             to: join(__dirname, 'static/images'),
        //         },
        //     ],
        // }),
        // new CopyPlugin({
        //     patterns: [
        //         {
        //             from: join(__dirname, 'front/favicon.ico'),
        //             to: join(__dirname, 'static/favicon.ico'),
        //         },
        //     ],
        // }),
        // new CopyPlugin({
        //     patterns: [
        //         {
        //             from: join(__dirname, 'front/manifest.json'),
        //             to: join(__dirname, 'static/manifest.json'),
        //         },
        //     ],
        // }),
    ]
}