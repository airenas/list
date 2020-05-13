import { TranscriptionService } from '../service/transcription.service';
import { MatSnackBar } from '@angular/material/snack-bar';

export abstract class BaseComponent {

  constructor(protected transcriptionService: TranscriptionService, protected snackBar: MatSnackBar) { }

  showError(msg: string, error: any) {
    console.error('Error', error);
    this.snackBar.open(this.asString(msg, error), null, { duration: 3000 });
  }

  asString(msg: string, error: any): string {
    return msg +  ' ' + String(error);
  }

  showInfo(info: any) {
    console.log('Info ', info);
    this.snackBar.open(info, null, { duration: 3000 });
  }
}
