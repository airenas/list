import { Injectable } from '@angular/core';

@Injectable()
export abstract class ParamsProviderService {
  lastSelectedFile: File;
  abstract setEmail(email: string): void;
  abstract getEmail(): string;
  abstract setTranscriptionID(id: string): void;
  abstract getTranscriptionID(): string;
}

export class LocalStorageParamsProviderService implements ParamsProviderService {

  private _transcriptionID: string;
  lastSelectedFile: File;
  private _email: string;

  constructor() {
  }

  setEmail(email: string): void {
    this._email = email;
    localStorage.setItem('email', email);
  }

  getEmail(): string {
    if (this._email == null) {
      this._email = localStorage.getItem('email');
    }
    return this._email;
  }

  setTranscriptionID(id: string): void {
    this._transcriptionID = id;
    localStorage.setItem('transcriptionID', id);
  }

  getTranscriptionID(): string {
    if (this._transcriptionID == null) {
      this._transcriptionID = localStorage.getItem('transcriptionID');
    }
    return this._transcriptionID;
  }
}
