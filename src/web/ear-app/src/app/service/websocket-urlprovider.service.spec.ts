import { TestBed, inject } from '@angular/core/testing';

import { WebsocketURLProviderService } from './websocket-urlprovider.service';
import { TestAppModule } from '../base/test.app.module';
import { RouterTestingModule } from '@angular/router/testing';
import { Config } from '../config';

export class MockLocation implements Location {
  ancestorOrigins: DOMStringList;
  hash: string; host: string;
  hostname: string;
  href: string;
  origin: string;
  pathname: string;
  port: string;
  protocol: string;
  search: string;
  assign(url: string): void {
  }

  reload(forcedReload?: boolean): void {
  }

  replace(url: string): void {
  }

  toString(): string {
    return '';
  }
}

describe('WebsocketURLProviderService', () => {
  let location: MockLocation;
  let config: Config;

  beforeEach(() => {
    config = new Config();
    TestBed.configureTestingModule({
      providers: [WebsocketURLProviderService,
        { provide: Config, useValue: config }
      ],
      imports: [TestAppModule, RouterTestingModule.withRoutes([])]
    });
    location = new MockLocation();
    location.protocol = 'http';
    location.port = '8000';
    location.hostname = 'host';
    location.pathname = '/';
  });

  it('should be created', inject([WebsocketURLProviderService], (service: WebsocketURLProviderService) => {
    expect(service).toBeTruthy();
  }));

  it('returns ws', inject([WebsocketURLProviderService], (service: WebsocketURLProviderService) => {
    config.subscribeUrl = 'olia';
    expect(service.getURLInternal(location, 'result')).toEqual('ws://host:8000/olia');
  }));

  it('returns path from config', inject([WebsocketURLProviderService], (service: WebsocketURLProviderService) => {
    config.subscribeUrl = 'olia';
    expect(service.getURLInternal(location, 'result')).toEqual('ws://host:8000/olia');
  }));

  it('use pathname', inject([WebsocketURLProviderService], (service: WebsocketURLProviderService) => {
    config.subscribeUrl = 'olia';
    location.pathname = '/lala';
    expect(service.getURLInternal(location, '')).toEqual('ws://host:8000/lala/olia');
  }));

  it('strip router path', inject([WebsocketURLProviderService], (service: WebsocketURLProviderService) => {
    config.subscribeUrl = 'olia';
    location.pathname = '/lala/result';
    expect(service.getURLInternal(location, 'result')).toEqual('ws://host:8000/lala/olia');
  }));

  it('handles double slash', inject([WebsocketURLProviderService], (service: WebsocketURLProviderService) => {
    config.subscribeUrl = '/olia';
    location.pathname = '/lala/result';
    expect(service.getURLInternal(location, '/result')).toEqual('ws://host:8000/lala/olia');
  }));

  it('handles no slashes', inject([WebsocketURLProviderService], (service: WebsocketURLProviderService) => {
    config.subscribeUrl = 'olia';
    location.pathname = 'lala/result';
    expect(service.getURLInternal(location, '/result')).toEqual('ws://host:8000/lala/olia');
  }));

});
