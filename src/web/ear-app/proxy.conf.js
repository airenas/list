const PROXY_CONFIG = {
  "/transcriber/": {
    "target": "http://localhost:8080",
    "secure": false,
    "pathRewrite": {"^/transcriber": ""}
  },
  "/result.provider/": {
    "target": "http://localhost:8081",
    "secure": false,
    "pathRewrite": {"^/result.provider": ""}
  },
  "/subscribe/": {
    "target": "ws://localhost:8082",
    "secure": false,
    "ws": true
  }
}

module.exports = PROXY_CONFIG;
