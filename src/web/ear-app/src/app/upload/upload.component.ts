import { SendFileResult } from './../api/send-file-result';
import { Component, OnInit } from '@angular/core';
import { TranscriptionService } from '../service/transcription.service';
import { Router } from '@angular/router';
import { MatSnackBar, ErrorStateMatcher } from '@angular/material';
import { BaseComponent } from '../base/base.component';
import { ParamsProviderService } from '../service/params-provider.service';
import { Validators, FormControl, FormGroupDirective, NgForm } from '@angular/forms';

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
    private router: Router, protected snackBar: MatSnackBar, private paramsProviderService: ParamsProviderService) {
    super(transcriptionService, snackBar);
  }

  selectedFile: File; // hold our file
  selectedFileName: string; // hold our file name
  private _email: string;

  ngOnInit() {
    this.fileChange(this.paramsProviderService.lastSelectedFile);
    this._email = this.paramsProviderService.email;
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
}
