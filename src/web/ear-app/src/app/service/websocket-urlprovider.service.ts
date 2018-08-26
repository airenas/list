import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Config } from '../config';

@Injectable()
export class WebsocketURLProviderService {
  constructor(private config: Config, private router: Router) { }

  public getURL() {
    return this.getURLInternal(window.location, this.router.url);
  }

  getURLInternal(location: Location, routerURL: string): string {
    const basePathURL = this.basePathName(location.pathname, routerURL);
    let result = this.getProtocol(location) + location.hostname + this.getPort(location);
    result = this.addURL(result, basePathURL);
    result = this.addURL(result, this.config.subscribeUrl);
    return result;
  }

  private getProtocol(location: Location) {
    return window.location.protocol === 'https:' ? 'wss://' : 'ws://';
  }

  private getPort(location: Location) {
    const defaultPort = location.protocol === 'https:' ? '443' : '80';
    if (location.port !== defaultPort) {
      return ':' + location.port;
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
