import { EditorURLProviderService } from './../service/editor-urlprovider.service';
import { MicrophoneFactory } from './../utils/microphone';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatTooltipModule } from '@angular/material/tooltip';
import { Config } from './../config';
import { NgModule, Injectable } from '@angular/core';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { APP_BASE_HREF } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTabsModule } from '@angular/material/tabs';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NgxFilesizeModule } from 'ngx-filesize';
import { TranscriptionService } from '../service/transcription.service';
import { Observable } from 'rxjs/Observable';
import { FileData } from '../service/file-data';
import { SendFileResult } from '../api/send-file-result';
import { TranscriptionResult } from '../api/transcription-result';
import { ActivatedRoute } from '@angular/router';
import { ParamsProviderService } from '../service/params-provider.service';
import { ResultSubscriptionService } from '../service/result-subscription.service';
import { EMPTY, of } from 'rxjs';
import { AudioPlayerFactory } from '../utils/audio.player';
import { TestMicrophoneFactory } from '../utils/microphone.specs';
import { TestAudioPlayerFactory } from '../utils/audio.player.specs';
import { TestParamsProviderService } from '../service/params-provider.service.spec';
import { Recognizer } from '../api/recognizer';
import { RouterTestingModule } from '@angular/router/testing';

@Injectable()
export class MockTestService implements TranscriptionService {
  getRecognizers(): Observable<Recognizer[]> {
    return of<Recognizer[]>([{ id: 'rID', name: 'rName', description: 'rDescr' }]);
  }
  sendFile(fileData: FileData): Observable<SendFileResult> {
    return EMPTY;
  }

  getResult(id: string): Observable<TranscriptionResult> {
    return EMPTY;
  }
}

@Injectable()
export class MockSubscriptionService implements ResultSubscriptionService {
  connect(): Observable<TranscriptionResult> {
    return EMPTY;
  }
  send(id: string): void {
  }
}

@Injectable()
export class MockActivatedRoute {
  snapshot = { paramMap: new Map(), queryParamMap: new Map() };
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
    if (element !== null && element.name === '#document' && element.parent === null) {
      return true;
    }
    if (element === null || element.nativeElement === null) {
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
    NoopAnimationsModule,
    MatTabsModule, MatButtonModule, MatInputModule, FormsModule,
    MatSnackBarModule, MatProgressBarModule, MatProgressSpinnerModule,
    MatCardModule, NgxFilesizeModule, MatSelectModule, MatMenuModule,
    ReactiveFormsModule,
    RouterTestingModule
  ],
  providers: [{ provide: APP_BASE_HREF, useValue: '/' },
  { provide: ParamsProviderService, useClass: TestParamsProviderService },
  { provide: TranscriptionService, useClass: MockTestService },
  { provide: ResultSubscriptionService, useClass: MockSubscriptionService },
  { provide: ActivatedRoute, useClass: MockActivatedRoute },
  { provide: Config, useClass: Config },
  { provide: AudioPlayerFactory, useClass: TestAudioPlayerFactory },
  { provide: MicrophoneFactory, useClass: TestMicrophoneFactory },
  { provide: EditorURLProviderService, useClass: EditorURLProviderService }
  ],
  bootstrap: [],
  exports: [
    MatTabsModule, MatButtonModule, MatInputModule, FormsModule,
    MatSnackBarModule, MatCardModule, MatProgressBarModule, MatProgressSpinnerModule, MatSelectModule, MatTooltipModule,
    MatMenuModule, NoopAnimationsModule,
    NgxFilesizeModule
  ],
})

export class TestAppModule {
}
