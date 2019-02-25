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

export class Progress {
  color: string;
  value: number;
  buffer: number;
}

class AudioURLKeeper {
  private ID: string = null;
  URL: string = null;
  private lastLoadedURL: string = null;
  isAudio = false;

  constructor(private config: Config, private audioPlayer: AudioPlayer) {
  }

  setAudio(result: TranscriptionResult) {
    console.log('keeper ID: ' + (result == null ? 'null' : result.id));
    if (result == null || result.status === Status.NOT_FOUND) {
      this.URL = null;
      this.ID = null;
      this.isAudio = false;
      return;
    }
    if (result.id === this.ID) {
      return;
    }
    this.ID = result.id;
    this.URL = this.config.audioUrl + result.id;
    this.isAudio = this.ID != null;
    if (this.isAudio && this.lastLoadedURL !== this.URL) {
      console.log('load URL: ' + this.URL);
      this.audioPlayer.load(this.URL);
      this.lastLoadedURL = this.URL;
    }
  }
}

class FileURLKeeper {
  contains = false;
  result: string;
  nbest: string;
  latticeGz: string;
  lattice: string;

  constructor(private config: Config) {
  }

  setID(result: TranscriptionResult) {
    this.contains = false;
    if (result == null || result.status !== Status.Completed) {
      return;
    }
    this.contains = true;
    this.result = this.config.resultUrl + result.id + '/result.txt';
    this.lattice = this.config.resultUrl + result.id + '/lat.txt';
    this.latticeGz = this.config.resultUrl + result.id + '/lat.gz';
    this.nbest = this.config.resultUrl + result.id + '/lat.nb10.txt';
  }
}

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
  audioURLKeeper: AudioURLKeeper = null;
  fileKeeper: FileURLKeeper = null;


  constructor(protected transcriptionService: TranscriptionService, protected snackBar: MatSnackBar,
    private route: ActivatedRoute, private paramsProviderService: ParamsProviderService,
    private resultSubscriptionService: ResultSubscriptionService,
    private cdr: ChangeDetectorRef, private audioPlayerFactory: AudioPlayerFactory, private config: Config) {
    super(transcriptionService, snackBar);
  }

  ngOnInit() {
    this.transcriptionId = this.route.snapshot.paramMap.get('id');
    this.audioPlayer = this.audioPlayerFactory.create('#audioWaveDiv', (ev) => this.cdr.detectChanges());
    this.audioURLKeeper = new AudioURLKeeper(this.config, this.audioPlayer);
    this.fileKeeper = new FileURLKeeper(this.config);
    if (this.transcriptionId == null) {
      this.transcriptionId = this.paramsProviderService.getTranscriptionID();
    } else {
      this.paramsProviderService.setTranscriptionID(this.transcriptionId);
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
        error => this.showError('Nepavyko gauti informacijos apie transkripcijos ID.', error)
      );
    }
  }

  transcriptionIdUpdated() {
    this.paramsProviderService.setTranscriptionID(this.transcriptionId);
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
    }
    this.audioURLKeeper.setAudio(this.result);
    this.fileKeeper.setID(this.result);
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
