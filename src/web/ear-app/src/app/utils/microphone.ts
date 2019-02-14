import { Injectable } from '@angular/core';
import Recorder from 'recorder-js';

@Injectable()
export abstract class Microphone {
    recording: boolean;
    abstract stop();
    abstract start();
}

@Injectable()
export class MicrophoneFactory {
    create(divName: string, eventHandler: RecorderEvent): Microphone {
        return new WebSurferMicrophone(divName, eventHandler);
    }
}

declare var WaveSurfer: any;

export type RecorderEvent = (event: string, data: any) => void;

@Injectable()
export class WebSurferMicrophone implements Microphone {
    recording = false;
    private wavesurfer: any = null;
    private recorder: Recorder = null;

    stop() {
        if (this.wavesurfer != null) {
            this.recorder.stop().then(({ blob, buffer }) => {
                this.eventHandler('data', blob);
//                this.fileChange(new File([blob], 'audio.wav'));
            });
            this.recording = false;
            this.wavesurfer.microphone.stop();
        }
    }
    start() {
        this.recording = true;
        if (this.initMicrophone()) {
            this.wavesurfer.microphone.start();
        } else {
            this.recording = false;
        }
    }

    constructor(private divName: string, private eventHandler: RecorderEvent) {
    }

    initMicrophone(): boolean {
        if (this.wavesurfer == null) {
            this.wavesurfer = WaveSurfer.create({
                container: this.divName,
                waveColor: 'blue',
                interact: false,
                cursorWidth: 0,
                height: 40,
                plugins: [
                    WaveSurfer.microphone.create()
                ]
            });
            this.wavesurfer.microphone.on('deviceReady', stream => {
                const audioContext = new AudioContext();
                this.recorder = new Recorder(audioContext, {});
                this.recorder.init(stream);
                this.recorder.start();
            });
            this.wavesurfer.microphone.on('deviceError', code => {
                this.recording = false;
                this.eventHandler('error', code);
//                console.error('Device error: ' + code);
//                this.showError('Nepavyko inicializuoti mikrofono.', <any>code);
            });
        }
        return this.wavesurfer != null;
    }
}
