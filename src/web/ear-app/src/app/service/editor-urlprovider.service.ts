import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Config } from '../config';

@Injectable()
export class EditorURLProviderService {
  constructor(private config: Config, private router: Router) { }

  public getURL(audioURL: string, latticeURL: string): string {
    return this.getURLInternal(window.location, this.router.url, audioURL, latticeURL);
  }

  getURLInternal(location: Location, routerURL: string, audioURL: string, latticeURL: string): string {
    const basePathURL = this.basePathName(location.pathname, routerURL);
    let mainUrl = location.protocol + '//' + location.hostname + this.getPort(location);
    mainUrl = this.addURL(mainUrl, basePathURL);
    let result = mainUrl;
    result = this.addURL(result, this.config.editorUrl);
    result = this.addURL(result, encodeURIComponent(this.addURL(mainUrl, latticeURL)));
    result = this.addURL(result, encodeURIComponent(this.addURL(mainUrl, audioURL)));
    return result;
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
