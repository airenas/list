import { TranscriptionResult } from './../api/transcription-result';
import { Config } from './../config';
import { Injectable } from '@angular/core';
import { FileData } from './file-data';
import { SendFileResult } from './../api/send-file-result';
import { HttpClient, HttpHeaders, HttpErrorResponse } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';
import { of } from 'rxjs';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import { throwError } from 'rxjs';
import { Recognizer } from '../api/recognizer';

@Injectable()
export abstract class TranscriptionService {
  abstract sendFile(fileData: FileData): Observable<SendFileResult>;
  abstract getResult(id: string): Observable<TranscriptionResult>;
  abstract getRecognizers(): Observable<Recognizer[]>;
}

@Injectable()
export class HttpTranscriptionService implements TranscriptionService {

  sendFileUrl: string;
  statusUrl: string;
  recognizersUrl: string;
  private socket;

  static asString(error: HttpErrorResponse): string {
    if (error !== null) {
      const value = String(error.error);
      if (value.includes('Wrong email')) {
        return 'Neteisingas El. paštas';
      }
      if (value.includes('No email')) {
        return 'Nenurodytas El. paštas';
      }
      if (value.includes('No file')) {
        return 'Nenurodytas failas';
      }
      if (value.includes('No recognizer') || value.includes('Unknown recognizer:')) {
        return 'Nepavyko parinkti atpažintuvą';
      }
    }
    return 'Serviso klaida';
  }

  constructor(public _http: HttpClient, _config: Config) {
    this.sendFileUrl = _config.sendFileUrl;
    this.statusUrl = _config.statusUrl;
    this.recognizersUrl = _config.recognizersUrl;
  }

  sendFile(fileData: FileData): Observable<SendFileResult> {
    const formData = new FormData();
    formData.append('file', fileData.file, fileData.fileName);
    formData.append('email', fileData.email);
    formData.append('recognizer', fileData.recognizer === '' ? '' : fileData.recognizer);
    formData.append('numberOfSpeakers', fileData.speakerCount === '' ? '' : fileData.speakerCount);
    if (fileData.skipNumJoin === true) {
      formData.append('skipNumJoin', '1');
    }
    const httpOptions = {
      headers: new HttpHeaders({
        'Accept': 'application/json'
      })
    };

    return this._http.post(this.sendFileUrl, formData, httpOptions)
      .map(res => {
        return <SendFileResult>res;
      })
      .catch(this.handleError);
  }

  getResult(id: string): Observable<TranscriptionResult> {
    const httpOptions = {
      headers: new HttpHeaders({
        'Accept': 'application/json'
      })
    };
    return this._http.get(this.statusUrl + id, httpOptions)
      .map(res => {
        return <TranscriptionResult>res;
      })
      .catch(e => this.handleError(e));
  }

  getRecognizers(): Observable<Recognizer[]> {
    const httpOptions = {
      headers: new HttpHeaders({
        'Accept': 'application/json'
      })
    };
    return this._http.get(this.recognizersUrl, httpOptions)
      .map(res => {
        return <Recognizer[]>res;
      })
      .catch(e => this.handleError(e));
  }

  protected handleError(error: HttpErrorResponse): Observable<never> {
    console.error(error);
    const errStr = HttpTranscriptionService.asString(error);
    return throwError(errStr);
  }

  protected getHeader() {
    const result = new Headers();
    // result.append('Content-Type', 'application/json');
    return result;
  }
}
