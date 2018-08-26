const PROXY_CONFIG = {
  "/transcriber/": {
    "target": "http://localhost:7050/ausis",
    "secure": false
    //"pathRewrite": {"^/transcriber": ""}
  },
  "/result.provider/": {
    "target": "http://localhost:7050/ausis",
    "secure": false
    //"pathRewrite": {"^/result.provider": ""}
  },
  "/result.provider/subscribe": {
    "target": "ws://localhost:7050/ausis",
    "secure": false,
    "ws": true,
    "logLevel": "debug"
  }
}

module.exports = PROXY_CONFIG;
