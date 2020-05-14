import { Injectable } from '@angular/core';

export type NamedEvent = (event: string) => void;

export abstract class AudioPlayer {
  abstract loadFile(file: File);
  abstract load(url: string);
  abstract clear();
  abstract play();
  abstract pause();
  abstract destroy();
  abstract isPlaying(): boolean;
}

declare var WaveSurfer: any;

export class WebSurferAudioPlayer implements AudioPlayer {

  wavesurfer: any = null;

  constructor(private divName: string, private eventHandler: NamedEvent) {
  }
  destroy() {
    if (this.wavesurfer != null) {
      console.log('Wavesurfer destroy');
      this.wavesurfer.destroy();
      this.wavesurfer = null;
    }
  }

  loadFile(file: File) {
    this.getSurfer().loadBlob(file);
  }

  load(url: string) {
    this.getSurfer().load(url);
  }

  clear() {
    if (this.wavesurfer != null) {
      console.log('Wavesurfer clear');
      this.wavesurfer.empty();
    }
  }
  play() {
    this.getSurfer().play();
  }
  pause() {
    this.getSurfer().pause();
  }
  isPlaying(): boolean {
    return this.wavesurfer != null && this.wavesurfer.isPlaying();
  }

  getSurfer(): any {
    if (this.wavesurfer == null) {
      console.log('Wavesurfer create');
      this.wavesurfer = WaveSurfer.create({
        fillParent: true,
        scrollParent: false,
        container: this.divName,
        waveColor: 'grey',
        progressColor: 'blue',
        height: 40
      });
      this.wavesurfer.on('pause', () => { this.handle('pause'); });
      this.wavesurfer.on('play', () => { this.handle('play'); });
    }
    return this.wavesurfer;
  }

  handle(event: string) {
    if (this.eventHandler != null) {
      console.log('Wavesurfer event ' + event);
      this.eventHandler(event);
    }
  }
}

@Injectable()
export class AudioPlayerFactory {
  create(divName: string, handler: NamedEvent): AudioPlayer {
    return new WebSurferAudioPlayer(divName, handler);
  }
}
