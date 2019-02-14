import { TestAudioPlayerFactory } from './../utils/audio.player.specs';
import { HttpTranscriptionService, TranscriptionService } from './../service/transcription.service';
import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UploadComponent } from './upload.component';
import { MockActivatedRoute, MockSubscriptionService, MockTestService } from '../base/test.app.module';
import { TestAppModule, FileHelper, TestHelper } from '../base/test.app.module';
import { FileSizeModule } from 'ngx-filesize';
import { RouterTestingModule } from '@angular/router/testing';
import { By } from '@angular/platform-browser';
import { ParamsProviderService } from '../service/params-provider.service';
import { APP_BASE_HREF } from '@angular/common';
import { ResultSubscriptionService } from '../service/result-subscription.service';
import { ActivatedRoute } from '@angular/router';
import { AudioPlayerFactory } from '../utils/audio.player';
import { MicrophoneFactory } from '../utils/microphone';
import { TestMicrophoneFactory } from '../utils/microphone.specs';


describe('UploadComponent', () => {
  let component: UploadComponent;
  let fixture: ComponentFixture<UploadComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [UploadComponent],
      imports: [TestAppModule, FileSizeModule, RouterTestingModule.withRoutes([])],
      providers: [
        { provide: AudioPlayerFactory, useClass: TestAudioPlayerFactory },
        { provide: MicrophoneFactory, useClass: TestMicrophoneFactory }]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UploadComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should have File placeholder', async(() => {
    expect(fixture.debugElement.query(By.css('#fileInput'))
      .nativeElement.getAttribute('placeholder')).toBe('Failas');
  }));

  it('should have El Pastas placeholder', async(() => {
    expect(fixture.debugElement.query(By.css('#emailInput'))
      .nativeElement.getAttribute('placeholder')).toBe('El. paštas');
  }));

  it('should have file data when file selected', async(() => {
    component.fileChange(new FileHelper().createFakeFile());
    fixture.debugElement.query(By.css('#hiddenFileInput')).nativeElement.dispatchEvent(new Event('input'));
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      const input = fixture.debugElement.query(By.css('#fileInput'));
      const el = input.nativeElement;
      expect(el.value).toBe('file.wav');
    });
  }));

  it('should have readonly button', async(() => {
    expect(fixture.debugElement.query(By.css('#uploadButton')).nativeElement.disabled).toBe(true);
  }));

  it('should have enabled button on valid Input', async(() => {
    expect(fixture.debugElement.query(By.css('#uploadButton')).nativeElement.disabled).toBe(true);
    component.fileChange(new FileHelper().createFakeFile());
    component.email = 'olia';
    fixture.debugElement.query(By.css('#hiddenFileInput')).nativeElement.dispatchEvent(new Event('input'));
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(fixture.debugElement.query(By.css('#uploadButton')).nativeElement.disabled).toBe(false);
    });
  }));

  it('should invoke upload on click', async(() => {
    component.fileChange(new FileHelper().createFakeFile());
    component.email = 'olia';
    fixture.debugElement.query(By.css('#hiddenFileInput')).nativeElement.dispatchEvent(new Event('input'));
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      spyOn(component, 'upload');
      fixture.debugElement.query(By.css('#uploadButton')).nativeElement.click();
      expect(component.upload).toHaveBeenCalled();
    });
  }));

  it('should be play controls', async(() => {
    component.fileChange(new FileHelper().createFakeFile());
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(fixture.debugElement.query(By.css('#playAudioButton')).nativeElement.disabled).toBe(false);
      expect(fixture.debugElement.query(By.css('#stopAudioButton'))).toBeNull();
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#audioWaveDiv')))).toBe(true);
      expect(component.audioPlayer.isPlaying()).toBe(false);
    });
  }));

  it('should be stop audio button', async(() => {
    component.fileChange(new FileHelper().createFakeFile());
    fixture.detectChanges();
    fixture.debugElement.query(By.css('#playAudioButton')).nativeElement.click();
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(fixture.debugElement.query(By.css('#playAudioButton'))).toBeNull();
      expect(fixture.debugElement.query(By.css('#stopAudioButton')).nativeElement.disabled).toBe(false);
      expect(component.audioPlayer.isPlaying()).toBe(true);
    });
  }));

  it('should invoke stop audio button', async(() => {
    component.fileChange(new FileHelper().createFakeFile());
    fixture.detectChanges();
    fixture.debugElement.query(By.css('#playAudioButton')).nativeElement.click();
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(component.audioPlayer.isPlaying()).toBe(true);
      fixture.debugElement.query(By.css('#stopAudioButton')).nativeElement.click();
      fixture.detectChanges();
      fixture.whenStable().then(() => {
        expect(component.audioPlayer.isPlaying()).toBe(false);
      });
    });
  }));

  it('should be no play controls', async(() => {
    component.fileChange(null);
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(fixture.debugElement.query(By.css('#playAudioButton'))).toBeNull();
      expect(fixture.debugElement.query(By.css('#stopAudioButton'))).toBeNull();
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#audioWaveDiv')))).toBe(false);
    });
  }));

  it('should be record controls', async(() => {
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(fixture.debugElement.query(By.css('#startRecordButton')).nativeElement.disabled).toBe(false);
      expect(fixture.debugElement.query(By.css('#stopRecordButton'))).toBeNull();
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#micWaveDiv')))).toBe(false);
      expect(component.audioPlayer.isPlaying()).toBe(false);
    });
  }));

  it('should invoke record event', async(() => {
    fixture.detectChanges();
    fixture.debugElement.query(By.css('#startRecordButton')).nativeElement.click();
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(fixture.debugElement.query(By.css('#startRecordButton'))).toBeNull();
      expect(fixture.debugElement.query(By.css('#stopRecordButton')).nativeElement.disabled).toBe(false);
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#micWaveDiv')))).toBe(true);
    });
  }));

  it('should invoke select file on input click', async(() => {
    spyOn(component, 'openInput');
    fixture.debugElement.query(By.css('#fileInput')).nativeElement.click();
    expect(component.openInput).toHaveBeenCalled();
  }));

  it('should invoke select file on button click', async(() => {
    spyOn(component, 'openInput');
    fixture.debugElement.query(By.css('#selectFileButton')).nativeElement.click();
    expect(component.openInput).toHaveBeenCalled();
  }));

});

describe('UploadComponent Own Mock', () => {
  let component: UploadComponent;
  let fixture: ComponentFixture<UploadComponent>;

  it('should read File value from provider', async(() => {
    const params = new ParamsProviderService();
    params.lastSelectedFile = new FileHelper().createFakeFile();
    TestBed.configureTestingModule({
      declarations: [UploadComponent],
      imports: [TestAppModule, FileSizeModule, RouterTestingModule.withRoutes([])],
      providers: [{ provide: ParamsProviderService, useValue: params },
      { provide: APP_BASE_HREF, useValue: '/' },
      { provide: TranscriptionService, useClass: MockTestService },
      { provide: ResultSubscriptionService, useClass: MockSubscriptionService },
      { provide: AudioPlayerFactory, useClass: TestAudioPlayerFactory },
      { provide: MicrophoneFactory, useClass: TestMicrophoneFactory },
      { provide: ActivatedRoute, useClass: MockActivatedRoute }]
    })
      .compileComponents();
    fixture = TestBed.createComponent(UploadComponent);
    component = fixture.debugElement.componentInstance;
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      const input = fixture.debugElement.query(By.css('#fileInput'));
      const el = input.nativeElement;
      expect(el.value).toBe('file.wav');
    });
  }));

  it('should read email value from provider', async(() => {
    const params = new ParamsProviderService();
    params.email = 'olia';

    TestBed.configureTestingModule({
      declarations: [UploadComponent],
      imports: [TestAppModule, FileSizeModule, RouterTestingModule.withRoutes([])],
      providers: [{ provide: ParamsProviderService, useValue: params },
      { provide: APP_BASE_HREF, useValue: '/' },
      { provide: TranscriptionService, useClass: MockTestService },
      { provide: ResultSubscriptionService, useClass: MockSubscriptionService },
      { provide: AudioPlayerFactory, useClass: TestAudioPlayerFactory },
      { provide: MicrophoneFactory, useClass: TestMicrophoneFactory },
      { provide: ActivatedRoute, useClass: MockActivatedRoute }]
    })
      .compileComponents();
    fixture = TestBed.createComponent(UploadComponent);
    component = fixture.debugElement.componentInstance;
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      const input = fixture.debugElement.query(By.css('#emailInput'));
      const el = input.nativeElement;
      expect(el.value).toBe('olia');
    });
  }));
});
