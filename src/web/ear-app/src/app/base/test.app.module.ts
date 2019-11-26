import { MicrophoneFactory } from './../utils/microphone';
import { MatProgressSpinnerModule, MatSelectModule, MatTooltipModule } from '@angular/material';
import { Config } from './../config';
import { NgModule } from '@angular/core';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { APP_BASE_HREF } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http';
import { MatTabsModule, MatButtonModule, MatInputModule, MatSnackBarModule, MatProgressBarModule } from '@angular/material';
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
import { AudioPlayerFactory } from '../utils/audio.player';
import { TestMicrophoneFactory } from '../utils/microphone.specs';
import { TestAudioPlayerFactory } from '../utils/audio.player.specs';
import { TestParamsProviderService } from '../service/params-provider.service.spec';
import { Recognizer } from '../api/recognizer';

export class MockTestService implements TranscriptionService {
  getRecognizers(): Observable<Recognizer[]> {
    return EMPTY;
  }
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

export class TestHelper {
  static Visible(element: any): boolean {
    if (element == null || element.nativeElement == null) {
      return false;
    }
    return !element.nativeElement.hasAttribute('hidden') && (element.parent == null || TestHelper.Visible(element.parent));
  }
}

@NgModule({
  declarations: [],
  imports: [
    BrowserModule,
    HttpClientModule,
    BrowserAnimationsModule,
    MatTabsModule, MatButtonModule, MatInputModule, FormsModule,
    MatSnackBarModule, MatProgressBarModule, MatProgressSpinnerModule,
    MatCardModule, FileSizeModule, MatSelectModule,
    ReactiveFormsModule
  ],
  providers: [{ provide: APP_BASE_HREF, useValue: '/' },
  { provide: ParamsProviderService, useClass: TestParamsProviderService },
  { provide: TranscriptionService, useClass: MockTestService },
  { provide: ResultSubscriptionService, useClass: MockSubscriptionService },
  { provide: ActivatedRoute, useClass: MockActivatedRoute },
  { provide: Config, useClass: Config },
  { provide: AudioPlayerFactory, useClass: TestAudioPlayerFactory },
  { provide: MicrophoneFactory, useClass: TestMicrophoneFactory }
  ],
  bootstrap: [],
  exports: [
    MatTabsModule, MatButtonModule, MatInputModule, FormsModule,
    MatSnackBarModule, MatCardModule, MatProgressBarModule, MatProgressSpinnerModule, MatSelectModule, MatTooltipModule
  ],
})

export class TestAppModule {
}
