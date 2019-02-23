import { TestBed, inject } from '@angular/core/testing';

import { ParamsProviderService } from './params-provider.service';

export class TestParamsProviderService implements ParamsProviderService {

  private _transcriptionID: string;
  lastSelectedFile: File;
  private _email: string;

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
}


describe('ParamsProviderService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [ParamsProviderService]
    });
  });

  it('should be created', inject([ParamsProviderService], (service: ParamsProviderService) => {
    expect(service).toBeTruthy();
  }));
});
