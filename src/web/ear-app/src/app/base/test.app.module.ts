import { NgModule } from '@angular/core';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { APP_BASE_HREF } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http';
import { MatTabsModule, MatButtonModule, MatInputModule, MatSnackBarModule } from '@angular/material';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { FileSizeModule } from 'ngx-filesize';
import { TranscriptionService } from '../service/transcription.service';
import { Observable } from 'rxjs/Observable';
import { FileData } from '../service/file-data';
import { SendFileResult } from '../api/send-file-result';
import { TranscriptionResult } from '../api/transcription-result';
import { ActivatedRoute } from '@angular/router';
import { ParamsProviderService } from '../service/params-provider.service';
import { ResultSubscriptionService } from '../service/result-subscription.service';
import { EMPTY } from 'rxjs';

export class MockTestService implements TranscriptionService {
  sendFile(fileData: FileData): Observable<SendFileResult> {
    return EMPTY;
  }

  getResult(id: string): Observable<TranscriptionResult> {
    return EMPTY;
  }
}

export class MockSubscriptionService implements ResultSubscriptionService {
  connect(): Observable<TranscriptionResult> {
    return EMPTY;
  }
  send(id: string): void {
  }
}

export class MockActivatedRoute {
  snapshot = { paramMap: new Map() };
}

export class FileHelper {
  createFakeFile(): File {
    const arrayOfBlob = new Array<Blob>();
    arrayOfBlob.push(new Blob(['Hello Wav'], { type: 'audio/wav' }));
    return new File(arrayOfBlob, 'file.wav', { type: 'audio/wav' });
  }
}

@NgModule({
  declarations: [],
  imports: [
    BrowserModule,
    HttpClientModule,
    BrowserAnimationsModule,
    MatTabsModule, MatButtonModule, MatInputModule, FormsModule,
    MatSnackBarModule,
    MatCardModule, FileSizeModule,
    ReactiveFormsModule
  ],
  providers: [ParamsProviderService, { provide: APP_BASE_HREF, useValue: '/' },
    { provide: TranscriptionService, useClass: MockTestService },
    { provide: ResultSubscriptionService, useClass: MockSubscriptionService },
    { provide: ActivatedRoute, useClass: MockActivatedRoute }],
  bootstrap: [],
  exports: [
    MatTabsModule, MatButtonModule, MatInputModule, FormsModule,
    MatSnackBarModule, MatCardModule
  ],
})

export class TestAppModule {
}
