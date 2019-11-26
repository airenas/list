import { TestBed, inject } from '@angular/core/testing';
import {} from 'jasmine';

import { ParamsProviderService, LocalStorageParamsProviderService } from './params-provider.service';

export class TestParamsProviderService implements ParamsProviderService {

  private _transcriptionID: string;
  lastSelectedFile: File;
  private _email: string;
  private _recognizer: string;


  constructor() {
  }

  setEmail(email: string): void {
    this._email = email;
  }

  getEmail(): string {
    return this._email;
  }

  setTranscriptionID(id: string): void {
    this._transcriptionID = id;
  }

  getTranscriptionID(): string {
    return this._transcriptionID;
  }

  setRecognizer(recognizer: string): void {
    this._recognizer = recognizer;
  }

  getRecognizer(): string {
    return this._recognizer;
  }
}

describe('ParamsProviderService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [{ provide: ParamsProviderService, useClass: LocalStorageParamsProviderService }]
    });
  });

  it('should be created', inject([ParamsProviderService], (service: ParamsProviderService) => {
    expect(service).toBeTruthy();
  }));

  it('should remember email', inject([ParamsProviderService], (service: ParamsProviderService) => {
    service.setEmail('olia');
    expect(service.getEmail()).toBe('olia');
  }));
  it('should remember ID', inject([ParamsProviderService], (service: ParamsProviderService) => {
    service.setTranscriptionID('id');
    expect(service.getTranscriptionID()).toBe('id');
  }));
  it('should remember recognizer', inject([ParamsProviderService], (service: ParamsProviderService) => {
    service.setRecognizer('rolia');
    expect(service.getRecognizer()).toBe('rolia');
  }));
  it('should remember email from local storage', inject([ParamsProviderService], (service: ParamsProviderService) => {
    service.setEmail('olia2');
    expect(new LocalStorageParamsProviderService().getEmail()).toBe('olia2');
  }));
  it('should remember ID from local storage', inject([ParamsProviderService], (service: ParamsProviderService) => {
    service.setTranscriptionID('id2');
    expect(new LocalStorageParamsProviderService().getTranscriptionID()).toBe('id2');
  }));
  it('should remember recognizer from local storage', inject([ParamsProviderService], (service: ParamsProviderService) => {
    service.setRecognizer('rec');
    expect(service.getRecognizer()).toBe('rec');
  }));
});
