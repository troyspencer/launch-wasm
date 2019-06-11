const HtmlWebpackPlugin = require('html-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const {GenerateSW} = require('workbox-webpack-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin'); 
const webpack = require('webpack'); // to access built-in plugins
const path = require('path');
const fs  = require('fs');
const lessToJs = require('less-vars-to-js');
const themeVariables = lessToJs(fs.readFileSync(path.join(__dirname, './components/theme.less'), 'utf8'));

var config = {
  resolve: {
    modules: ['components', 'node_modules']
  },
  devtool: 'source-map',
  stats: {
    cached: false,
    cachedAssets: false,
    chunks: false,
    chunkModules: false,
    chunkOrigins: false,
    modules: false
  },
  entry: {
    vendor: ['@babel/polyfill', 'react', 'react-dom'],
    client:     './components/index.js',
  },
  output: {
    path: __dirname + '/dist',
    filename: '[name].chunkhash.bundle.js',
    chunkFilename: '[name].chunkhash.bundle.js',
    publicPath: '/',
  },
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        loader: "babel-loader",
        options: {
          plugins: [
            ['import', { libraryName: "antd", style: true }]
          ]
        }
      },
      {
        test: /\.less$/,
        use: [
          {loader: 'style-loader'},
          {loader: 'css-loader'}, 
          {
            loader: 'less-loader', // compiles Less to CSS
            options: {
              modifyVars: themeVariables,
              root: path.resolve(__dirname, './'),
              javascriptEnabled: true,
           },
          }
        ],
      },

      {
        test: /\.(png|jpe?g|gif|svg|eot|ttf|woff|woff2)$/,
        loader: 'url-loader',
        options: {
          limit: 8192,
        },
      },
     ]
  },
  devServer: {
    historyApiFallback: true,
    disableHostCheck: true
  }
};

module.exports = (env, argv) => {
  if (argv.mode === 'production') {
    config.output.path = __dirname + '/../server/dist'
  }

  config.plugins =  [
    new webpack.ProgressPlugin(),
    new CleanWebpackPlugin(),
    new CopyWebpackPlugin([
      { from: 'static', to: 'static'}
    ]),
    new HtmlWebpackPlugin({
      title: 'Launch',
      template: path.join(__dirname, 'index.html'),
      filename: 'index.html',
      inject: true,
      minify: {
        collapseWhitespace: true,
        collapseInlineTagWhitespace: true,
        minifyCSS: true,
        minifyURLs: true,
        minifyJS: true,
        removeComments: true,
        removeRedundantAttributes: true
      }
    }),
    new GenerateSW(),
  ]

  return config
}