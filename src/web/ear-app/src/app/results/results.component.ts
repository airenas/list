import { EditorURLProviderService } from './../service/editor-urlprovider.service';
import { ResultSubscriptionService } from './../service/result-subscription.service';
import { TranscriptionResult } from './../api/transcription-result';
import { MatMenuTrigger } from '@angular/material/menu';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Component, OnInit, OnDestroy, ChangeDetectorRef, ViewChild } from '@angular/core';
import { TranscriptionService } from '../service/transcription.service';
import { ActivatedRoute } from '@angular/router';
import { BaseComponent } from '../base/base.component';
import { ParamsProviderService } from '../service/params-provider.service';
import { Subscription } from 'rxjs/Subscription';
import { Status } from '../api/status';
import { AudioPlayer, AudioPlayerFactory } from '../utils/audio.player';
import { Config } from '../config';
import { ErrorPipe } from '../pipes/error.pipe';

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
    //    console.log('keeper ID: ' + (result == null ? 'null' : result.id));
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
      //      console.log('load URL: ' + this.URL);
      this.audioPlayer.load(this.URL);
      this.lastLoadedURL = this.URL;
    }
  }

  getURL() {
    return this.URL;
  }
}

class FileURLKeeper {

  constructor(private config: Config) {
  }

  static latRestoredURL = 'lat.restored.txt';

  contains = false;
  result: string;
  resultFinal: string;
  nbest: string;
  lattice: string;
  latticeGz: string;
  latticeRestored: string;
  latticeRestoredGz: string;
  webVTT: string;
  URLPrefix: string;

  dwnFiles = [
    {
      id: 'dfResult',
      url: 'result.txt',
      title: 'Rezultatas - paprastas (.txt)'
    },
    {
      id: 'dfResultFinal',
      url: 'resultFinal.txt',
      title: 'Rezultatas (.txt)'
    },
    {
      id: 'dfLat',
      url: 'lat.txt',
      title: 'Kaldi grafas (.txt)'
    },
    {
      id: 'dfLatGz',
      url: 'lat.gz',
      title: 'Kaldi grafas (.gz)'
    },
    {
      id: 'dfN10',
      url: 'lat.nb10.txt',
      title: 'Kaldi 10-geriausių variantų (.txt)'
    },
    {
      id: 'dfLatRescore',
      url: FileURLKeeper.latRestoredURL,
      title: 'Grafas redagavimui (.txt)'
    },
    {
      id: 'dfLatRescoreGz',
      url: 'lat.restored.gz',
      title: 'Grafas redagavimui (.gz)'
    },
    {
      id: 'dfWebVTT',
      url: 'webvtt.txt',
      title: 'WebVTT (.txt)'
    },
  ];

  setID(result: TranscriptionResult) {
    this.contains = false;
    if (result == null || result.status !== Status.Completed) {
      this.URLPrefix = '';
      return;
    }
    this.contains = true;
    this.URLPrefix = this.config.resultUrl + result.id;
  }

  getURL(URLSufix: string) {
    return this.URLPrefix + '/' + URLSufix;
  }

  download(URLSufix: string) {
    const URL = this.getURL(URLSufix);
    window.open(URL, '_blank');
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
  resultSubscription: Subscription;
  progress: Progress;
  status: string;
  audioPlayer: AudioPlayer;
  audioURLKeeper: AudioURLKeeper = null;
  fileKeeper: FileURLKeeper = null;
  errDetailsClick = 0;
  @ViewChild('menuTrigger') menuTrigger: MatMenuTrigger;


  constructor(protected transcriptionService: TranscriptionService, protected snackBar: MatSnackBar,
    private route: ActivatedRoute, private paramsProviderService: ParamsProviderService,
    private resultSubscriptionService: ResultSubscriptionService,
    private cdr: ChangeDetectorRef, private audioPlayerFactory: AudioPlayerFactory, private config: Config,
    private editorUrlSrvice: EditorURLProviderService) {
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
      console.log('Got msg from Websocket');
      this.onResult(message);
    });
    this.refresh();
  }

  ngOnDestroy() {
    console.log('Destroy result');
    this.audioPlayer.destroy();
    this.resultSubscription.unsubscribe();
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
    this.recognizedText = null;
    this.progress = null;
    this.status = null;
    if (this.result != null) {
      this.recognizedText = result.recognizedText;
      this.resultSubscriptionService.send(this.result.id);
      this.progress = this.prepareProgress(this.result);
      this.status = (result.status === Status.Completed) ? null : result.status;
    }
    this.audioURLKeeper.setAudio(this.result);
    this.fileKeeper.setID(this.result);
    this.setError();
  }

  setError() {
    if (this.result) {
      this.error = new ErrorPipe(this.paramsProviderService.showErrorDetails).transform(this.result);
    } else {
      this.error = null;
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

  showErrDetails() {
    this.errDetailsClick++;
    if (this.errDetailsClick > 4 && this.paramsProviderService.showErrorDetails !== true) {
      this.paramsProviderService.showErrorDetails = true;
      this.showInfo('Detalus klaidų rodymas įjungtas');
      this.setError();
    }
  }

  openEditor() {
    const url = this.editorUrlSrvice.getURL(this.audioURLKeeper.getURL(),
      this.fileKeeper.getURL(FileURLKeeper.latRestoredURL));
    console.log('Editor: ' + decodeURIComponent(url));
    window.open(url);
  }
}
