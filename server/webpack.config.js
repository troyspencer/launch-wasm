const HtmlWebpackPlugin = require('html-webpack-plugin');
const AddAssetHtmlPlugin = require('add-asset-html-webpack-plugin');

module.exports = {
  resolve: {
    modules: ['static', 'node_modules']
  },
  devtool: 'source-map',
  entry: {
    vendor: ['@babel/polyfill', 'react', 'react-dom'],
    client:     './static/index.js',
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
        use: {
          loader: "babel-loader"
        }
      },

     ]
  },
  devServer: {
    historyApiFallback: true,
    disableHostCheck: true
  },
  plugins: [
    new HtmlWebpackPlugin({
      title: 'GoWasm!',
      template: './static/index.html',
      filename: './index.html',
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
    new AddAssetHtmlPlugin({ filepath: require.resolve('./static/init_go.js') }),
    new AddAssetHtmlPlugin({ filepath: require.resolve('./static/wasm_exec.js') }),
    new AddAssetHtmlPlugin({ filepath: require.resolve('./static/manifest.json') }),
    new AddAssetHtmlPlugin({ filepath: require.resolve('./static/launch-144.png') }),
    new AddAssetHtmlPlugin({ filepath: require.resolve('./static/launch-192.png') })
  ]
};