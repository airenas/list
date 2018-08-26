import { TestBed, inject } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';

import { TranscriptionService, HttpTranscriptionService } from './transcription.service';
import { TestAppModule } from '../base/test.app.module';

describe('HttpTranscriptionService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule, TestAppModule],
      providers: [HttpTranscriptionService]
    });
  });

  it('should be created', inject([HttpTranscriptionService], (service: HttpTranscriptionService) => {
    expect(service).toBeTruthy();
  }));
});
