module.exports = {
    configureWebpack: {
      devServer: {
        headers: { "Access-Control-Allow-Origin": "*" },
        disableHostCheck: true,
        proxy: "http://localhost:8888/",
      }
    }
  };
