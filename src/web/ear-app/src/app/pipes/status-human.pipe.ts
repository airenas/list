import { Pipe, PipeTransform } from '@angular/core';
import { Status } from '../api/status';

@Pipe({
  name: 'statusHuman'
})
export class StatusHumanPipe implements PipeTransform {

  transform(value: string): string {
    if (value === Status.Uploaded) {
      return 'Įkeltas. Laukia';
    }
    if (value === Status.Completed) {
      return 'Pabaigta';
    }
    if (value === Status.AudioConvert) {
      return 'Konvertuojamas audio failas';
    }
    if (value === Status.Diarization) {
      return 'Skaidomas audio failas';
    }
    if (value === Status.Transcription) {
      return 'Transkribuojamas';
    }
    if (value === Status.Rescore) {
      return 'Perskaičiuojamas su kalbos modeliu';
    }
    if (value === Status.ResultMake) {
      return 'Ruošiamas rezultatas';
    }
    if (value === Status.NOT_FOUND) {
      return 'Nežinomas ID';
    }
    return value;
  }
}
