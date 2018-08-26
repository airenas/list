import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Config } from '../config';

@Injectable()
export class WebsocketURLProviderService {
  constructor(private config: Config, private router: Router) { }

  public getURL() {
    console.log('window.location.hostname: ' + window.location.hostname);
    console.log('router.url                ' + this.router.url);
    console.log('window.location.pathname: ' + window.location.pathname);
    const basePathURL = this.basePathName(window.location.pathname, this.router.url);
    const r = this.addURL(
      this.addURL(this.websocketProtocolByLocation() + window.location.hostname
        + this.websocketPortWithColonByLocation(), basePathURL),
      this.config.subscribeUrl);
    console.log('getURL: ' + r);
    return r;
  }

  private websocketProtocolByLocation() {
    return window.location.protocol === 'https:' ? 'wss://' : 'ws://';
  }

  private websocketPortWithColonByLocation() {
    const defaultPort = window.location.protocol === 'https:' ? '443' : '80';
    if (window.location.port !== defaultPort) {
      return ':' + window.location.port;
    } else {
      return '';
    }
  }

  private basePathName(allPath, routersURL) {
    if (allPath.endsWith(routersURL)) {
      return allPath.substring(0, allPath.length - routersURL.length);
    }
    return '';
  }

  private addURL(s1, s2) {
    if (s1 && s2 && s1.length > 0 && s2.length > 0) {
      if (s1.endsWith('/') && s2.startsWith('/')) {
        return s1 + s2.substring(1);
      }
      if (!s1.endsWith('/') && !s2.startsWith('/')) {
        return s1 + '/' + s2;
      }
    }
    return s1 + s2;
  }
}
