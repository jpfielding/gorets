const env = process.env.CONFIG_ENV || 'local';

const webpack = require('webpack');
const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const cssnext = require('postcss-cssnext');
const postcssFlexbugFixes = require('postcss-flexbugs-fixes');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const postcssImport = require("postcss-import")


const config = require(`./config/${env}`);

console.log(`Building gorets explorer with env: ${env} config: ${config}`);

const definePlugin = new webpack.DefinePlugin({
  config: JSON.stringify(config),
});

const providePlugin = new webpack.ProvidePlugin({
  fetch: 'imports?this=>global!exports?global.fetch!whatwg-fetch',
});

const htmlPlugin = new HtmlWebpackPlugin({
  filename: 'index.html',
  template: 'mustache!./src/index.html',
  inject: false,
  conf: config,
});

module.exports = {
  entry: {
    app: ['babel-polyfill', './src/index.jsx'],
  },
  resolve: {
    root: path.resolve('./src'),
    extensions: ['', '.js', '.jsx'],
  },
  output: {
    path: path.join(__dirname, '/build'),
    filename: 'bundle.js',
    publicPath: `${config.staticAssetPath}/`,
  },
  devtool: 'source-map',
  debug: true,
  plugins: [
    definePlugin,
    htmlPlugin,
    providePlugin,
    new ExtractTextPlugin('[name].css'),
    // new CopyWebpackPlugin([
    //   { from: 'src/index.mustache' },
    // ]),
  ],
  externals: {
    react: 'window.React',
    'react-dom': 'window.ReactDOM',
  },
  module: {
    preLoaders: [
      {
        test: /.(js|jsx)$/i,
        loader: 'eslint-loader',
        include: [/src/],
        exclude: /node_modules/,
      },
    ],
    loaders: [
      {
        test: /.(js|jsx)$/i,
        exclude: /node_modules/,
        include: [/src/],
        loader: 'babel',
        query: {
          presets: [
            'es2015',
            'react',
            'stage-2',
          ],
        },
      },
      // { test: /(page|manifest)\.json/i, loader: 'file?name=[name].[ext]' },
      // { test: /\.json/i, loader: 'json-loader' },
      // { test: /\.mustache$/i, loader: 'raw' },
      { test: /\.css$/,
        loader: ExtractTextPlugin.extract('style-loader', 'css-loader?minimize&-autoprefixer!postcss-loader') },
      // { test: /\.(woff|woff2|svg|ico)$/i, loader: 'url?limit=10000' },
      // { test: /\.(ttf|eot|jpe?g|png|gif)$/i, loader: 'file' },
    ],
  },
  postcss: () => (
    [cssnext({
      browsers: ['> 1%', 'last 2 versions', 'iOS 8', 'iOS 7'],
      flexbox: true,
    }),
    postcssImport,
    postcssFlexbugFixes]
  ),
  devServer: {
    headers: { 'Access-Control-Allow-Origin': '*' },
    historyApiFallback: true,
    compress: true,
  },
};
