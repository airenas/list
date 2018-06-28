import { Injectable } from '@angular/core';

@Injectable()
export class ParamsProviderService {

  lastId: string;
  lastSelectedFile: File;

  constructor() {
  }
}
