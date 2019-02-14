import { Injectable } from '@angular/core';
import { NamedEvent, AudioPlayer } from './audio.player';

@Injectable()
export class TestAudioPlayerFactory {
    create(divName: string, handler: NamedEvent): AudioPlayer {
        return new TestAudioPlayer(divName, handler);
    }
}

@Injectable()
export class TestAudioPlayer implements AudioPlayer {

    playing = false;

    constructor(private divName: string, private eventHandler: NamedEvent) {
    }

    loadFile(file: File) {
    }

    clear() {
    }

    play() {
        this.playing = true;
    }

    pause() {
        this.playing = false;
    }

    isPlaying(): boolean {
        return this.playing;
    }
}
