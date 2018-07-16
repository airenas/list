import { ResultSubscriptionService } from './../service/result-subscription.service';
import { TranscriptionResult } from './../api/transcription-result';
import { MatSnackBar } from '@angular/material';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { TranscriptionService } from '../service/transcription.service';
import { ActivatedRoute } from '@angular/router';
import { BaseComponent } from '../base/base.component';
import { ParamsProviderService } from '../service/params-provider.service';
import { Observable } from 'rxjs';
import { Subscription } from 'rxjs/Subscription';

@Component({
  selector: 'app-results',
  templateUrl: './results.component.html',
  styleUrls: ['./results.component.css']
})
export class ResultsComponent extends BaseComponent implements OnInit, OnDestroy {
  transcriptionId: string;
  result: TranscriptionResult;
  error: string;
  recognizedText: string;
  private resultSubscription: Subscription;


  constructor(protected transcriptionService: TranscriptionService, protected snackBar: MatSnackBar,
    private route: ActivatedRoute, private paramsProviderService: ParamsProviderService,
    private resultSubscriptionService: ResultSubscriptionService) {
    super(transcriptionService, snackBar);
  }

  ngOnInit() {
    this.transcriptionId = this.route.snapshot.paramMap.get('id');
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
        error => this.showError('Nepavyko gauti informacijos apie transkripcijos ID',  <any>error)
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
    if (this.result) {
      this.error = result.error;
      this.recognizedText = result.recognizedText;
      this.resultSubscriptionService.send(this.result.id);
    }
  }

  ngOnDestroy() {
    this.resultSubscription.unsubscribe();
  }
}
