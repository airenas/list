import { Pipe, PipeTransform } from '@angular/core';
import { TranscriptionResult } from '../api/transcription-result';
import { ErrorCode } from '../api/error-codes';

@Pipe({
  name: 'error'
})
export class ErrorPipe implements PipeTransform {

  transform(value: TranscriptionResult): string {
    if (value.errorCode === ErrorCode.TooShortAudio) {
      return 'Per trumpas įrašas';
    }
    if (value.errorCode === ErrorCode.WrongFormat) {
      return 'Blogas formatas';
    }
    return value.error;
  }
}
