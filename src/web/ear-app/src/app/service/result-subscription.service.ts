import { TranscriptionResult } from './../api/transcription-result';
import { Injectable } from '@angular/core';
import websocketConnect from 'rxjs-websockets';
import { QueueingSubject } from 'queueing-subject';
import { Observable } from 'rxjs/Observable';
import { Config } from '../config';

@Injectable()
export class ResultSubscriptionService {
  private inputStream: QueueingSubject<string>;

  constructor(private config: Config) {
  }

  public connect(): Observable<TranscriptionResult> {
    // Using share() causes a single websocket to be created when the first
    // observer subscribes. This socket is shared with subsequent observers
    // and closed when the observer count falls to zero.
    return websocketConnect(
        this.config.subscribeUrl,
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
