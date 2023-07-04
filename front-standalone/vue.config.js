const path = require("path");
function resolve(dir) {
  return path.join(__dirname, dir);
}

module.exports = {
    publicPath: process.env.VUE_APP_ROUTER_BASE,
    outputDir: 'dist',
    devServer: {
        host:'localhost',
        port: 8080
    },
    lintOnSave: false,
    chainWebpack: (config) => {
        config.resolve.alias
          .set("/@", resolve("src"))
          .set("/@u", resolve("src/utils"));
        config.resolve.extensions
            .add('ts')
            .add('tsx');
    }
}
