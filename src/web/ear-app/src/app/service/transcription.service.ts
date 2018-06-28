import { TranscriptionResult } from './../api/transcription-result';
import { Config } from './../config';
import { Injectable } from '@angular/core';
import { FileData } from './file-data';
import { SendFileResult } from './../api/send-file-result';
import { Http, RequestOptions, Headers } from '@angular/http';
import { Observable } from 'rxjs/Rx';
import { map } from 'rxjs/operators';


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

  constructor(public _http: Http, _config: Config) {
    this.sendFileUrl = _config.sendFileUrl;
    this.resultUrl = _config.resultUrl;
  }

  sendFile(fileData: FileData): Observable<SendFileResult> {
    const formData = new FormData();
    formData.append('file', fileData.file, fileData.fileName);
    const headers = new Headers();
    headers.append('Accept', 'application/json');
    const options = new RequestOptions({ headers: headers });
    return this._http.post(this.sendFileUrl, formData, options)
      .map(res => {
        return <SendFileResult>res.json();
      })
      .catch(this.handleError);
  }

  getResult(id: string): Observable<TranscriptionResult> {
    const headers = new Headers();
    headers.append('Accept', 'application/json');
    const options = new RequestOptions({ headers: headers });
    return this._http.get(this.resultUrl + id, options)
      .map(res => {
        return <TranscriptionResult>res.json();
      })
      .catch(this.handleError);
  }

  protected handleError(error: Response) {
    console.error(error);
    return Observable.throw('Serviso klaida: ' + error.text() || 'Serviso klaida');
  }

  protected getHeader() {
    const result = new Headers();
    // result.append('Content-Type', 'application/json');
    return result;
  }
}