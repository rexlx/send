module.exports = {
    configureWebpack: {
      devServer: {
        headers: { "Access-Control-Allow-Origin": "*" },
        disableHostCheck: true,
        proxy: "http://storage.nullferatu.com:8888/",
      }
    }
  };
