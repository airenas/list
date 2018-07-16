import { SendFileResult } from './../api/send-file-result';
import { Component, OnInit } from '@angular/core';
import { TranscriptionService } from '../service/transcription.service';
import { Router } from '@angular/router';
import { MatSnackBar } from '@angular/material';

export abstract class BaseComponent {
  constructor(protected transcriptionService: TranscriptionService, protected snackBar: MatSnackBar) { }

  showError(msg: string, error: any) {
    console.error('Error sending file', error);
    this.snackBar.open(msg, null, { duration: 3000 });
  }

  showInfo(info: any) {
    console.log('Info ', info);
    this.snackBar.open(info, null, { duration: 3000 });
  }
}
