import { SendFileResult } from './../api/send-file-result';
import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { TranscriptionService } from '../service/transcription.service';
import { Router } from '@angular/router';
import { MatSnackBar, ErrorStateMatcher } from '@angular/material';
import { BaseComponent } from '../base/base.component';
import { ParamsProviderService } from '../service/params-provider.service';
import { Validators, FormControl, FormGroupDirective, NgForm } from '@angular/forms';
import Recorder from 'recorder-js';
import { AudioPlayer, WebSurferAudioPlayer, AudioPlayerFactory } from '../utils/audio.player';

declare var WaveSurfer: any;

export class MyErrorStateMatcher implements ErrorStateMatcher {
  isErrorState(control: FormControl | null, form: FormGroupDirective | NgForm | null): boolean {
    const isSubmitted = form && form.submitted;
    return !!(control && control.invalid && (control.dirty || control.touched || isSubmitted));
  }
}

@Component({
  selector: 'app-upload',
  templateUrl: './upload.component.html',
  styleUrls: ['./upload.component.css']
})
export class UploadComponent extends BaseComponent implements OnInit {
  constructor(protected transcriptionService: TranscriptionService,
    private router: Router, protected snackBar: MatSnackBar, private paramsProviderService: ParamsProviderService,
    private cdr: ChangeDetectorRef, private audioPlayerFactory: AudioPlayerFactory) {
    super(transcriptionService, snackBar);
  }

  selectedFile: File; // hold our file
  selectedFileName: string; // hold our file name
  private _email: string;
  private wavesurfer: any = null;
  private audioPlayer: AudioPlayer;
  recording = false;
  private recorder: Recorder = null;

  ngOnInit() {
    this.audioPlayer = this.audioPlayerFactory.create('#audioWaveDiv', (ev) => this.cdr.detectChanges());
    this.fileChange(this.paramsProviderService.lastSelectedFile);
    this._email = this.paramsProviderService.email;
    this.recording = false;
  }

  openInput() {
    document.getElementById('hiddenFileInput').click();
  }

  filesChange(files: File[]) {
    if (files.length > 0) {
      this.fileChange(files[0]);
    } else {
      this.fileChange(null);
    }
  }

  fileChange(file: File) {
    this.selectedFile = null;
    this.selectedFileName = null;
    this.paramsProviderService.lastSelectedFile = file;
    if (file) {
      this.selectedFile = file;
      this.selectedFileName = this.selectedFile.name;
    }
    this.showAudioFile(file);
  }

  upload() {
    console.log('sending this to server', this.selectedFile);
    this.transcriptionService.sendFile({ file: this.selectedFile, fileName: this.selectedFileName, email: this.email })
      .subscribe(
        result => this.onResult(result),
        error => this.showError('Nepavyko nusiųsti failo', <any>error)
      );

  }

  onResult(result: SendFileResult) {
    this.fileChange(null);
    this.showInfo('Failas nusiųstas. Transkripcijos ID: ' + result.id);
    this.router.navigateByUrl('/results/' + result.id);
  }

  get email(): string {
    return this._email;
  }
  set email(email: string) {
    this._email = email;
    this.paramsProviderService.email = email;
  }

  isValid() {
    return this.selectedFile && this._email;
  }

  canPlayAudio(): boolean {
    return !this.audioPlayer.isPlaying() && this.selectedFile != null;
  }

  canStopAudio(): boolean {
    return this.audioPlayer.isPlaying();
  }

  showAudioFile(file: File) {
    if (file != null) {
      this.audioPlayer.loadFile(file);
    } else {
      this.audioPlayer.clear();
    }
  }

  playAudio() {
    this.audioPlayer.play();
  }

  stopAudio() {
    this.audioPlayer.pause();
  }

  startRecord() {
    this.recording = true;
    if (this.initMicrophone()) {
      this.wavesurfer.microphone.start();
    } else {
      this.recording = false;
    }
  }

  stopRecord() {
    if (this.wavesurfer != null) {
      this.recorder.stop().then(({ blob, buffer }) => {
        this.fileChange(new File([blob], 'audio.wav'));
      });
      this.recording = false;
      this.wavesurfer.microphone.stop();
    }
  }

  initMicrophone(): boolean {
    if (this.wavesurfer == null) {
      this.wavesurfer = WaveSurfer.create({
        container: '#micWaveDiv',
        waveColor: 'blue',
        interact: false,
        cursorWidth: 0,
        height: 40,
        plugins: [
          WaveSurfer.microphone.create()
        ]
      });
      this.wavesurfer.microphone.on('deviceReady', stream => {
        const audioContext = new AudioContext();
        this.recorder = new Recorder(audioContext, {});
        this.recorder.init(stream);
        this.recorder.start();
      });
      this.wavesurfer.microphone.on('deviceError', code => {
        this.recording = false;
        console.error('Device error: ' + code);
        this.showError('Nepavyko inicializuoti mikrofono.', <any>code);
      });
    }
    return this.wavesurfer != null;
  }
}
