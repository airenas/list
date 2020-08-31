import { NgxFilesizeModule } from 'ngx-filesize';
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
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTabsModule } from '@angular/material/tabs';
import { MatMenuModule } from '@angular/material/menu';
import { MatSelectModule } from '@angular/material/select';
import { MatTooltipModule } from '@angular/material/tooltip';
import { ShowOnDirtyErrorStateMatcher, ErrorStateMatcher } from '@angular/material/core';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { Config } from './config';
import { HttpClientModule } from '@angular/common/http';
import { StatusHumanPipe } from './pipes/status-human.pipe';
import { ParamsProviderService, LocalStorageParamsProviderService } from './service/params-provider.service';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { WebsocketURLProviderService } from './service/websocket-urlprovider.service';
import { MicrophoneFactory } from './utils/microphone';
import { EditorURLProviderService } from './service/editor-urlprovider.service';

@NgModule({
  declarations: [
    AppComponent,
    UploadComponent,
    ResultsComponent,
    StatusHumanPipe,
    ResultTextPipe
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    BrowserAnimationsModule,
    MatTabsModule, MatButtonModule, MatInputModule, FormsModule,
    MatSnackBarModule, MatProgressBarModule, MatProgressSpinnerModule,
    MatCardModule, NgxFilesizeModule, MatSelectModule,
    ReactiveFormsModule, MatTooltipModule, MatMenuModule
  ],
  providers: [Config,
    WebsocketURLProviderService, EditorURLProviderService,
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
