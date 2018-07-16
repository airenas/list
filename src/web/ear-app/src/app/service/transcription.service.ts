import { TranscriptionResult } from './../api/transcription-result';
import { Config } from './../config';
import { Injectable } from '@angular/core';
import { FileData } from './file-data';
import { SendFileResult } from './../api/send-file-result';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import { throwError } from 'rxjs';

@Injectable()
export abstract class TranscriptionService {
  abstract sendFile(fileData: FileData): Observable<SendFileResult>;
  abstract getResult(id: string): Observable<TranscriptionResult>;
}

@Injectable()
export class HttpTranscriptionService implements TranscriptionService {

  sendFileUrl: string;
  resultUrl: string;
  private socket;

  constructor(public _http: HttpClient, _config: Config) {
    this.sendFileUrl = _config.sendFileUrl;
    this.resultUrl = _config.resultUrl;
  }

  sendFile(fileData: FileData): Observable<SendFileResult> {
    const formData = new FormData();
    formData.append('file', fileData.file, fileData.fileName);
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
    const headers = new Headers();
    headers.append('Accept', 'application/json');
    const httpOptions = {
      headers: new HttpHeaders({
        'Accept': 'application/json'
      })
    };
    return this._http.get(this.resultUrl + id, httpOptions)
      .map(res => {
        return <TranscriptionResult>res;
      })
      .catch(this.handleError);
  }

  protected handleError(error: Response) {
    console.error(error);
    return throwError(new Error('Serviso klaida: ' + error.text() || 'Serviso klaida'));
  }

  protected getHeader() {
    const result = new Headers();
    // result.append('Content-Type', 'application/json');
    return result;
  }
}
