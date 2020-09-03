import { TestBed, inject } from '@angular/core/testing';

import { TestAppModule } from '../base/test.app.module';
import { RouterTestingModule } from '@angular/router/testing';
import { Config } from '../config';
import { MockLocation } from './websocket-urlprovider.service.spec';
import { EditorURLProviderService } from './editor-urlprovider.service';

describe('EditorURLProviderService', () => {
  let location: MockLocation;
  let config: Config;

  beforeEach(() => {
    config = new Config();
    TestBed.configureTestingModule({
      providers: [EditorURLProviderService,
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

  it('should be created', inject([EditorURLProviderService], (service: EditorURLProviderService) => {
    expect(service).toBeTruthy();
  }));

  it('should use url from config', inject([EditorURLProviderService], (service: EditorURLProviderService) => {
    config.editorUrl = 'olia/';
    expect(service.getURLInternal(location, 'result', 'audio/id', 'result/id/lattice/lat.txt'))
      .toContain('http//host:8000/olia/');
  }));

  it('should strip router path', inject([EditorURLProviderService], (service: EditorURLProviderService) => {
    config.editorUrl = 'olia/';
    location.pathname = '/lala/result';
    expect(service.getURLInternal(location, 'result', 'audio/id', 'result/id/lattice/lat.txt'))
      .toContain('http//host:8000/lala/olia/');
  }));

  it('should add audio param', inject([EditorURLProviderService], (service: EditorURLProviderService) => {
    config.editorUrl = 'olia/';
    location.pathname = '/lala/result';
    expect(service.getURLInternal(location, 'result', 'audio/id', 'result/id/lattice/lat.txt'))
      .toContain('/' + encodeURIComponent('http//host:8000/lala/audio/id'));
  }));

  it('should add lattice param', inject([EditorURLProviderService], (service: EditorURLProviderService) => {
    config.editorUrl = 'olia/';
    location.pathname = '/lala/result';
    expect(service.getURLInternal(location, 'result', 'audio/id', 'result/id/lattice/lat.txt'))
      .toContain('/' + encodeURIComponent('http//host:8000/lala/result/id/lattice/lat.txt') + '/');
  }));
});
