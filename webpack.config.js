const path = require("path");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const TerserPlugin = require("terser-webpack-plugin");
const CssMinimizerPlugin = require("css-minimizer-webpack-plugin");
const CopyPlugin = require("copy-webpack-plugin");
const { CleanWebpackPlugin } = require("clean-webpack-plugin");

const staticPath = path.join(__dirname, "static");

module.exports = {
  entry: "./resource/js/app.js",
  // output setting
  output: {
    filename: "js/app.js", // js output
    path: staticPath, // output path
  },
  module: {
    rules: [
      {
        test: /\.html$/,
        loader: "html-loader",
      },
      {
        // extension .js case
        test: /\.js$/,
        use: [
          {
            // use babel
            loader: "babel-loader",
            // babel options
            options: {
              presets: [
                // es2020 to es5
                "@babel/preset-env",
              ],
            },
          },
        ],
      },
      {
        test: /\.svg$/,
        use: [
          {
            loader: "html-loader",
            options: {
              minimize: true,
            },
          },
        ],
      },
      {
        test: /\.(sa|sc|c)ss$/,
        exclude: /node_modules/,
        use: [
          MiniCssExtractPlugin.loader,
          {
            loader: "css-loader",
            options: { url: false },
          },
          "sass-loader",
        ],
      },
    ],
  },
  plugins: [
    new CleanWebpackPlugin(),
    new MiniCssExtractPlugin({
      filename: "css/app.css", // css output
    }),
    new CopyPlugin({
      patterns: [
        {
          from: "./node_modules/admin-lte/plugins",
          to: path.join(staticPath, "plugins"),
        },
        {
          from: "./resource/image",
          to: path.join(staticPath, "image"),
        },
      ],
    }),
  ],
  optimization: {
    minimize: true,
    minimizer: [
      new TerserPlugin(), // js minify
      new CssMinimizerPlugin(), // css minify
    ],
  },
  performance: { hints: false },
};
