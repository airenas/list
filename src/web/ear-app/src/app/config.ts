import {Injectable, Inject} from '@angular/core';
import {environment} from '../environments/environment';

@Injectable()
export class Config {
    public sendFileUrl: string;
    public statusUrl: string;
    public subscribeUrl: string;
    public audioUrl: string;
    public resultUrl: string;

    constructor() {
        const prefix = '';
        this.sendFileUrl = prefix + environment.sendFileUrl + 'upload';
        this.statusUrl = prefix + environment.statusUrl + 'status/';
        this.subscribeUrl = prefix + environment.statusUrl + 'subscribe';
        this.audioUrl = prefix + environment.resultUrl + 'audio/';
        this.resultUrl = prefix + environment.resultUrl + 'result/';
    }
}
