import { QueueingSubject } from 'queueing-subject';
import { WebsocketURLProviderService } from './websocket-urlprovider.service';
import { TranscriptionResult } from './../api/transcription-result';
import { Injectable } from '@angular/core';
import { Observable, Observer, Subject } from 'rxjs';
import 'rxjs/add/operator/share';

@Injectable()
export abstract class ResultSubscriptionService {
  abstract connect(): Observable<TranscriptionResult>;
  abstract send(id: string): void;
}

@Injectable()
export class WSResultSubscriptionService implements ResultSubscriptionService {
  private ws: WebSocket;
  private inputStream: QueueingSubject<string>;

  constructor(private urlProvider: WebsocketURLProviderService) {
    this.inputStream = new QueueingSubject<string>();
  }

  public connect(): Observable<TranscriptionResult> {
    if (this.ws != null) {
      this.ws.close();
    }

    console.log('Connecting to ' + this.urlProvider.getURL());
    this.ws = new WebSocket(this.urlProvider.getURL());
    this.inputStream = new QueueingSubject<string>();
    this.ws.onopen = (m => {
      console.log('Websocket opened');
      this.inputStream.subscribe((id: string) => {
        if (this.ws.readyState === WebSocket.OPEN) {
          console.log('Send to websocket ' + id);
          this.ws.send(id);
        } else {
          console.warn('Websocket in not ready state');
        }
      });
    });

    const observable = Observable.create((obs: Observer<MessageEvent>) => {
      this.ws.onmessage = obs.next.bind(obs);
      this.ws.onerror = obs.error.bind(obs);
      this.ws.onclose = (m => {
        console.log('got complete ' + m);
        this.inputStream.unsubscribe();
        obs.complete();
      });
      return this.ws.close.bind(this.ws);
    });
    return observable.map(
      (response: MessageEvent): TranscriptionResult => {
        return JSON.parse(response.data);
      }
    );
  }

  public send(id: string): void {
    console.log('Try send to websocket ' + id);
    this.inputStream.next(id);
  }
}
