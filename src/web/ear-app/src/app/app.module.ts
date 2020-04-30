import { ResultTextPipe } from './pipes/result-text.pipe';
import { AudioPlayerFactory } from './utils/audio.player';
import { ResultSubscriptionService, WSResultSubscriptionService } from './service/result-subscription.service';
import { HttpTranscriptionService, TranscriptionService } from './service/transcription.service';
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { UploadComponent } from './upload/upload.component';
import { ResultsComponent } from './results/results.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatTabsModule, MatButtonModule, MatInputModule, MatSnackBarModule, MatTooltipModule, MatSelectModule } from '@angular/material';
import { ShowOnDirtyErrorStateMatcher, ErrorStateMatcher, MatProgressSpinnerModule } from '@angular/material';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { FileSizeModule } from 'ngx-filesize';
import { Config } from './config';
import { HttpClientModule } from '@angular/common/http';
import { StatusHumanPipe } from './pipes/status-human.pipe';
import { ParamsProviderService, LocalStorageParamsProviderService } from './service/params-provider.service';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { WebsocketURLProviderService } from './service/websocket-urlprovider.service';
import { MicrophoneFactory } from './utils/microphone';
import { ErrorPipe } from './pipes/error.pipe';

@NgModule({
  declarations: [
    AppComponent,
    UploadComponent,
    ResultsComponent,
    StatusHumanPipe,
    ResultTextPipe,
    ErrorPipe
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    BrowserAnimationsModule,
    MatTabsModule, MatButtonModule, MatInputModule, FormsModule,
    MatSnackBarModule, MatProgressBarModule, MatProgressSpinnerModule,
    MatCardModule, FileSizeModule, MatSelectModule,
    ReactiveFormsModule, MatTooltipModule
  ],
  providers: [Config,
    WebsocketURLProviderService,
    { provide: ParamsProviderService, useClass: LocalStorageParamsProviderService },
    { provide: ResultSubscriptionService, useClass: WSResultSubscriptionService },
    { provide: TranscriptionService, useClass: HttpTranscriptionService },
    { provide: AudioPlayerFactory, useClass: AudioPlayerFactory },
    { provide: MicrophoneFactory, useClass: MicrophoneFactory },
    { provide: ErrorStateMatcher, useClass: ShowOnDirtyErrorStateMatcher },
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
