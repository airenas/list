import { Injectable } from '@angular/core';
import { NamedEvent, AudioPlayer } from './audio.player';

export class TestAudioPlayer implements AudioPlayer {
    playing = false;

    constructor(private divName: string, private eventHandler: NamedEvent) {
    }
    destroy() {
        this.playing = false;
    }

    loadFile(file: File) {
    }

    load(url: string) {
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

@Injectable()
export class TestAudioPlayerFactory {
    create(divName: string, handler: NamedEvent): AudioPlayer {
        return new TestAudioPlayer(divName, handler);
    }
}
