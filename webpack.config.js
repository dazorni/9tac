var webpack = require('webpack');
var path = require('path');

var BUILD_DIR = path.resolve(__dirname, 'build/public/js');
var APP_DIR = path.resolve(__dirname, 'app');

var config = {
  entry: APP_DIR + '/index.js',
  output: {
    path: BUILD_DIR,
    filename: 'app.js'
  },
  cache: true,
  devtool: 'eval-source-map',
  plugins: [
    new webpack.HotModuleReplacementPlugin(),
    new webpack.NoErrorsPlugin(),
	  new webpack.EnvironmentPlugin(["NODE_ENV"]),
  ],
  module : {
    loaders : [
      {
        test : /.+.js$/,
        include : APP_DIR,
        loader : 'babel'
      }
    ],
    preLoaders: [
      {
        test: /\.js?$/,
        loaders: ['eslint'],
        include: APP_DIR
      }
    ]
  }
};

module.exports = config;
