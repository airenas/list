import { TestBed, inject } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';

import { TranscriptionService, HttpTranscriptionService } from './transcription.service';
import { TestAppModule } from '../base/test.app.module';
import { HttpErrorResponse } from '@angular/common/http';

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

  function makeError(error: string): string {
    return HttpTranscriptionService.asString(new HttpErrorResponse({ error: error }));
  }

  it('should parse messages', (() => {
    expect(makeError('error')).toBe('Serviso klaida');
    expect(makeError('Wrong email')).toBe('Neteisingas El. paštas');
    expect(makeError('No email')).toBe('Nenurodytas El. paštas');
    expect(makeError('No file')).toBe('Nenurodytas failas');
    expect(makeError('No recognizer')).toBe('Nepavyko parinkti atpažintuvą');
    expect(makeError('Unknown recognizer: aaaa')).toBe('Nepavyko parinkti atpažintuvą');
  }));
});
