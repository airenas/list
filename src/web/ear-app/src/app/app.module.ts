import { ResultSubscriptionService, WSResultSubscriptionService } from './service/result-subscription.service';
import { HttpTranscriptionService, TranscriptionService } from './service/transcription.service';
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { UploadComponent } from './upload/upload.component';
import { ResultsComponent } from './results/results.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatTabsModule, MatButtonModule, MatInputModule, MatSnackBarModule, MatTooltipModule } from '@angular/material';
import { ShowOnDirtyErrorStateMatcher, ErrorStateMatcher } from '@angular/material';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { FileSizeModule } from 'ngx-filesize';
import { Config } from './config';
import { HttpClientModule } from '@angular/common/http';
import { StatusHumanPipe } from './pipes/status-human.pipe';
import { ParamsProviderService } from './service/params-provider.service';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { WebsocketURLProviderService } from './service/websocket-urlprovider.service';

@NgModule({
  declarations: [
    AppComponent,
    UploadComponent,
    ResultsComponent,
    StatusHumanPipe
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    BrowserAnimationsModule,
    MatTabsModule, MatButtonModule, MatInputModule, FormsModule,
    MatSnackBarModule, MatProgressBarModule,
    MatCardModule, FileSizeModule,
    ReactiveFormsModule, MatTooltipModule
  ],
  providers: [Config,
    ParamsProviderService,
    WebsocketURLProviderService,
    { provide: ResultSubscriptionService, useClass: WSResultSubscriptionService },
    { provide: TranscriptionService, useClass: HttpTranscriptionService },
    { provide: ErrorStateMatcher, useClass: ShowOnDirtyErrorStateMatcher },
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
