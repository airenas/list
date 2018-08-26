import {Injectable, Inject} from '@angular/core';
import {environment} from '../environments/environment';

@Injectable()
export class Config {
    public sendFileUrl: string;
    public resultUrl: string;
    public subscribeUrl: string;

    constructor() {
        const prefix = '';
        this.sendFileUrl = prefix + environment.sendFileUrl + 'upload';
        this.resultUrl = prefix + environment.resultUrl + 'result/';
        this.subscribeUrl = prefix + environment.subscribeUrl;
    }
}
