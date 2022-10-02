module.exports = {
    configureWebpack: {
      devServer: {
        headers: { "Access-Control-Allow-Origin": "*" },
        disableHostCheck: true,
        proxy: "http://api.domain.com:8888/",
      }
    }
  };
