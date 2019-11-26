const PROXY_CONFIG = {
  "/transcriber/": {
    "target": "http://localhost:80/ausis",
    "secure": false
    //"pathRewrite": {"^/transcriber": ""}
  },
  "/status.service/": {
    "target": "http://localhost:80/ausis",
    "secure": false
    //"pathRewrite": {"^/result.provider": ""}
  },
  "/result.service/": {
    "target": "http://localhost:80/ausis",
    "secure": false
  },
  "/status.service/subscribe": {
    "target": "ws://localhost:80/ausis",
    "secure": false,
    "ws": true,
    "logLevel": "debug"
  }
}

module.exports = PROXY_CONFIG;
