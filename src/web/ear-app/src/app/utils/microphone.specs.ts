import { Injectable } from '@angular/core';
import { NamedEvent, AudioPlayer } from './audio.player';
import { Microphone } from './microphone';

export class TestMicrophone implements Microphone {
    recording = false;
    constructor(private divName: string) {
    }

    stop() {
        this.recording = false;
    }

    start() {
        this.recording = true;
    }
}

@Injectable()
export class TestMicrophoneFactory {
    create(divName: string, handler: NamedEvent): Microphone {
        return new TestMicrophone(divName);
    }
}
