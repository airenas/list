import { Injectable } from '@angular/core';

export abstract class Microphone {
    recording: boolean;
    abstract stop();
    abstract start();
}

declare var WaveSurfer: any;
declare var WebAudioRecorder: any;

export type RecorderEvent = (event: string, data: any) => void;

export class WebSurferMicrophone implements Microphone {
    recording = false;
    private wavesurfer: any = null;
    private recorder: any = null;

    stop() {
        if (this.wavesurfer != null) {
            // console.log('stopping recording');
            this.recorder.finishRecording();
            // console.log('stopping wavesurfer');
            this.recording = false;
            this.wavesurfer.microphone.stop();
        }
    }
    start() {
        this.recording = true;
        if (this.initMicrophone()) {
            this.wavesurfer.microphone.start();
            this.eventHandler('started', null);
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
                const input = audioContext.createMediaStreamSource(stream);
                this.recorder = new WebAudioRecorder(input, {
                    workerDir: 'assets/',
                    numChannels: 1,
                    encoding: 'wav',
                    onEncoderLoading: (recorder, encoding) => {
//                        console.log('Loading ' + encoding + ' encoder...');
                    },
                    onEncoderLoaded: (recorder, encoding) => {
//                        console.log(encoding + ' encoder loaded');
                    }
                });

                this.recorder.onComplete = (recorder, blob) => {
//                    console.log('got recorded audio');
                    this.eventHandler('data', blob);
                };

                this.recorder.onTimeout = (recorder) => {
                    console.log('timeout');
                    this.stop();
                };

                this.recorder.setOptions({
                    timeLimit: 30,
                    encodeAfterRecord: true
                });

                this.recorder.startRecording();

//                console.log('Recording started');
            });
            this.wavesurfer.microphone.on('deviceError', code => {
                this.recording = false;
                console.error('Device error: ', code);
                this.eventHandler('error', this.asString(code));
            });
        }
        return this.wavesurfer != null;
    }
    asString(code: any): string {
        const value = this.codeAsString(code);
        if (value.includes('NotAllowedError') || value.includes('PermissionDeniedError')) {
            return 'Neleidžiama';
        }
        if (value.includes('NotSupportedError')) {
            return 'Neleidžiama/Nepalaikoma';
        }
        return value;
    }

    codeAsString(code: any): string {
        if (code == null) {
            return '';
        }
        const me = typeof(MediaStreamError) !== 'undefined' ? MediaStreamError : DOMException;
        if (code instanceof me) {
            return code.name + ':' + code.message;
        }
        return code.toString();
    }
}

@Injectable()
export class MicrophoneFactory {
    create(divName: string, eventHandler: RecorderEvent): Microphone {
        return new WebSurferMicrophone(divName, eventHandler);
    }
}
