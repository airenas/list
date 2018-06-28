import { TestBed, inject } from '@angular/core/testing';

import { TranscriptionService } from './transcription.service';

describe('TranscriptionService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [TranscriptionService]
    });
  });

  it('should be created', inject([TranscriptionService], (service: TranscriptionService) => {
    expect(service).toBeTruthy();
  }));
});
