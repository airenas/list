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
      this.inputStream.subscribe((id: string) => {
        if (this.ws.readyState === WebSocket.OPEN) {
          this.ws.send(id);
        }
      });
    });

    const observable = new Observable((obs: Observer<MessageEvent>) => {
      this.ws.onmessage = obs.next.bind(obs);
      this.ws.onerror = obs.error.bind(obs);
      this.ws.onclose = (m => {
        console.log('WebSocket closed');
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
    this.inputStream.next(id);
  }
}
