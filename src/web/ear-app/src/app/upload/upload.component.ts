import { SendFileResult } from './../api/send-file-result';
import { Component, OnInit } from '@angular/core';
import { TranscriptionService } from '../service/transcription.service';
import { Router } from '@angular/router';
import { MatSnackBar } from '@angular/material';
import { BaseComponent } from '../base/base.component';
import { ParamsProviderService } from '../service/params-provider.service';

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
  selectedFileName: string; // hold our file

  ngOnInit() {
    this.fileChange(this.paramsProviderService.lastSelectedFile);
  }


  /**
   * this is used to trigger the input
   */
  openInput() {
    // your can use ElementRef for this later
    document.getElementById('fileInput').click();
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

  /**
  * this is used to perform the actual upload
  */
  upload() {
    console.log('sending this to server', this.selectedFile);
    this.transcriptionService.sendFile({ file: this.selectedFile, fileName: this.selectedFileName })
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
}
