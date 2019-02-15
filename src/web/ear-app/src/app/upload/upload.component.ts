import { SendFileResult } from './../api/send-file-result';
import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { TranscriptionService } from '../service/transcription.service';
import { Router } from '@angular/router';
import { MatSnackBar, ErrorStateMatcher } from '@angular/material';
import { BaseComponent } from '../base/base.component';
import { ParamsProviderService } from '../service/params-provider.service';
import { FormControl, FormGroupDirective, NgForm } from '@angular/forms';
import { AudioPlayer, AudioPlayerFactory } from '../utils/audio.player';
import { Microphone, MicrophoneFactory } from '../utils/microphone';

export class MyErrorStateMatcher implements ErrorStateMatcher {
  isErrorState(control: FormControl | null, form: FormGroupDirective | NgForm | null): boolean {
    const isSubmitted = form && form.submitted;
    return !!(control && control.invalid && (control.dirty || control.touched || isSubmitted));
  }
}

@Component({
  selector: 'app-upload',
  templateUrl: './upload.component.html',
  styleUrls: ['./upload.component.scss']
})
export class UploadComponent extends BaseComponent implements OnInit {
  constructor(protected transcriptionService: TranscriptionService,
    private router: Router, protected snackBar: MatSnackBar, private paramsProviderService: ParamsProviderService,
    private cdr: ChangeDetectorRef, private audioPlayerFactory: AudioPlayerFactory,
    private microphoneFactory: MicrophoneFactory) {
    super(transcriptionService, snackBar);
  }

  selectedFile: File; // hold our file
  selectedFileName: string; // hold our file name
  private _email: string;
  recorder: Microphone;
  audioPlayer: AudioPlayer;

  ngOnInit() {
    this.audioPlayer = this.audioPlayerFactory.create('#audioWaveDiv', (ev) => this.cdr.detectChanges());
    this.recorder = this.microphoneFactory.create('#micWaveDiv', (ev, data) => this.recordEvent(ev, data));
    this.fileChange(this.paramsProviderService.lastSelectedFile);
    this._email = this.paramsProviderService.email;
  }

  recordEvent(ev: string, data: any): void {
    if (ev === 'data') {
      this.fileChange(new File([data], 'audio.wav'));
    } else if (ev === 'error') {
      this.showError('Nepavyko inicializuoti mikrofono.', data);
    }
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
    return this.selectedFile && this._email && !this.recorder.recording;
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
    this.recorder.start();
  }

  stopRecord() {
    this.recorder.stop();
  }
}
