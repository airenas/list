import { WebsocketURLProviderService } from './websocket-urlprovider.service';
import { TranscriptionResult } from './../api/transcription-result';
import { Injectable } from '@angular/core';
import websocketConnect from 'rxjs-websockets';
import { QueueingSubject } from 'queueing-subject';
import { Observable } from 'rxjs/Observable';
import { Config } from '../config';
import 'rxjs/add/operator/share';
import { Router } from '@angular/router';

@Injectable()
export abstract class ResultSubscriptionService {
  abstract connect(): Observable<TranscriptionResult>;
  abstract send(id: string): void;
}

@Injectable()
export class WSResultSubscriptionService implements ResultSubscriptionService {
  private inputStream: QueueingSubject<string>;

  constructor(private urlProvider: WebsocketURLProviderService) {
  }

  public connect(): Observable<TranscriptionResult> {
    // Using share() causes a single websocket to be created when the first
    // observer subscribes. This socket is shared with subsequent observers
    // and closed when the observer count falls to zero.
    return websocketConnect(
      this.urlProvider.getURL(),
      this.inputStream = new QueueingSubject<string>()
    ).messages.share().map(data => {
      return <TranscriptionResult>JSON.parse(data);
    }
    );
  }

  public send(id: string): void {
    // If the websocket is not connected then the QueueingSubject will ensure
    // that messages are queued and delivered when the websocket reconnects.
    // A regular Subject can be used to discard messages sent when the websocket
    // is disconnected.
    this.inputStream.next(id);
  }
}
