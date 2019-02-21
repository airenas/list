import { ResultSubscriptionService } from './../service/result-subscription.service';
import { TranscriptionResult } from './../api/transcription-result';
import { MatSnackBar } from '@angular/material';
import { Component, OnInit, OnDestroy, ChangeDetectorRef } from '@angular/core';
import { TranscriptionService } from '../service/transcription.service';
import { ActivatedRoute } from '@angular/router';
import { BaseComponent } from '../base/base.component';
import { ParamsProviderService } from '../service/params-provider.service';
import { Subscription } from 'rxjs/Subscription';
import { Status } from '../api/status';
import { AudioPlayer, AudioPlayerFactory } from '../utils/audio.player';
import { Config } from '../config';

@Component({
  selector: 'app-results',
  templateUrl: './results.component.html',
  styleUrls: ['./results.component.scss']
})
export class ResultsComponent extends BaseComponent implements OnInit, OnDestroy {
  transcriptionId: string;
  result: TranscriptionResult;
  error: string;
  recognizedText: string;
  private resultSubscription: Subscription;
  progress: Progress;
  status: string;
  audioPlayer: AudioPlayer;
  audioURLKeeper: AudioURLKeeper = new AudioURLKeeper(this.config);


  constructor(protected transcriptionService: TranscriptionService, protected snackBar: MatSnackBar,
    private route: ActivatedRoute, private paramsProviderService: ParamsProviderService,
    private resultSubscriptionService: ResultSubscriptionService,
    private cdr: ChangeDetectorRef, private audioPlayerFactory: AudioPlayerFactory, private config: Config) {
    super(transcriptionService, snackBar);
  }

  ngOnInit() {
    this.transcriptionId = this.route.snapshot.paramMap.get('id');
    this.audioPlayer = this.audioPlayerFactory.create('#audioWaveDiv', (ev) => this.cdr.detectChanges());
    this.audioURLKeeper.audioPlayer = this.audioPlayer;
    if (this.transcriptionId == null) {
      this.transcriptionId = this.paramsProviderService.lastId;
    } else {
      this.paramsProviderService.lastId = this.transcriptionId;
    }
    this.resultSubscription = this.resultSubscriptionService.connect().subscribe((message: TranscriptionResult) => {
      console.log('received message from server: ' + JSON.stringify(message));
      this.onResult(message);
    });
    this.refresh();
  }

  refresh() {
    this.onResult(null);
    if (this.transcriptionId) {
      this.transcriptionService.getResult(this.transcriptionId).subscribe(
        result => this.onResult(result),
        error => this.showError('Nepavyko gauti informacijos apie transkripcijos ID', <any>error)
      );
    }
  }

  transcriptionIdUpdated() {
    this.paramsProviderService.lastId = this.transcriptionId;
  }

  onResult(result: TranscriptionResult) {
    this.result = result;
    this.error = null;
    this.recognizedText = null;
    this.progress = null;
    this.status = null;
    if (this.result) {
      this.error = result.error;
      this.recognizedText = result.recognizedText;
      this.resultSubscriptionService.send(this.result.id);
      this.progress = this.prepareProgress(this.result);
      this.status = (result.status === Status.Completed) ? null : result.status;
      this.audioURLKeeper.setId(result.id);
    } else {
      this.audioURLKeeper.setId(null);
    }
  }

  prepareProgress(result: TranscriptionResult): Progress {
    if (result.status === Status.Completed) {
      return null;
    }
    const progress = new Progress();
    progress.value = result.progress;
    progress.color = result.error ? 'warn' : 'primary';
    progress.buffer = result.error || (result.status === Status.Completed) ? 100 : 90;
    return progress;
  }

  ngOnDestroy() {
    this.resultSubscription.unsubscribe();
  }

  canPlayAudio(): boolean {
    return !this.audioPlayer.isPlaying();
  }

  canStopAudio(): boolean {
    return this.audioPlayer.isPlaying();
  }

  playAudio() {
    this.audioPlayer.play();
  }

  stopAudio() {
    this.audioPlayer.pause();
  }
}

export class Progress {
  color: string;
  value: number;
  buffer: number;
}

class AudioURLKeeper {
  constructor(private config: Config) {

  }
  ID: string;
  audioPlayer: AudioPlayer;
  setId(ID: string) {
    console.log('keeper ID: ' + ID);
    if (ID === this.ID || ID == null) {
      return;
    }
    this.ID = ID;
    console.log('load ID: ' + ID);
    this.audioPlayer.load(this.config.audioUrl + ID);
  }
}
