import { Pipe, PipeTransform } from '@angular/core';
import { Status } from '../api/status';

@Pipe({
  name: 'resultText'
})
export class ResultTextPipe implements PipeTransform {
  transform(value: string): string {
    if (value === null) {
      return value;
    }
    const re = /\n/gi;
    return value.replace(re, '\n  ');
  }
}
