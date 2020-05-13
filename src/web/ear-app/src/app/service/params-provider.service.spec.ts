import { TestBed, inject } from '@angular/core/testing';
import {} from 'jasmine';

import { ParamsProviderService, LocalStorageParamsProviderService } from './params-provider.service';
import { Injectable } from '@angular/core';

@Injectable()
export class TestParamsProviderService implements ParamsProviderService {

  private _transcriptionID: string;
  lastSelectedFile: File;
  showErrorDetails = false;
  private _email: string;
  private _recognizer: string;
  private _speakerCount: string;


  constructor() {
  }
  setSpeakerCount(speakerCount: string): void {
    this._speakerCount = speakerCount;
  }
  getSpeakerCount(): string {
    return this._speakerCount;
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
    expect(new LocalStorageParamsProviderService().getRecognizer()).toBe('rec');
  }));
  it('should remember speakerCount', inject([ParamsProviderService], (service: ParamsProviderService) => {
    service.setSpeakerCount('123');
    expect(service.getSpeakerCount()).toBe('123');
  }));
  it('should remember speakerCount from local storage', inject([ParamsProviderService], (service: ParamsProviderService) => {
    service.setSpeakerCount('123');
    expect(new LocalStorageParamsProviderService().getSpeakerCount()).toBe('123');
  }));
  it('default showErrorDetails', inject([ParamsProviderService], (service: ParamsProviderService) => {
    expect(service.showErrorDetails).toBe(false);
  }));
  it('update showErrorDetails', inject([ParamsProviderService], (service: ParamsProviderService) => {
    service.showErrorDetails = true;
    expect(service.showErrorDetails).toBe(true);
  }));
});
