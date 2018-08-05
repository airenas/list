import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'statusHuman'
})
export class StatusHumanPipe implements PipeTransform {

  transform(value: string): string {
    if (value === 'RECEIVED') {
      return 'Įkeltas. Laukia';
    }
    if (value === 'NOT_FOUND') {
      return 'Nežinomas ID';
    }
    return value;
  }
}
