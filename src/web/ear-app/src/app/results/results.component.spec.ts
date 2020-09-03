import { ErrorCode } from './../api/error-codes';
import { ResultTextPipe } from './../pipes/result-text.pipe';
import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ResultsComponent, Progress } from './results.component';
import { TestAppModule, MockActivatedRoute, MockSubscriptionService } from '../base/test.app.module';
import { MockTestService, TestHelper } from '../base/test.app.module';
import { StatusHumanPipe } from '../pipes/status-human.pipe';
import { By } from '@angular/platform-browser';
import { ParamsProviderService } from '../service/params-provider.service';
import { TranscriptionService } from '../service/transcription.service';
import { ResultSubscriptionService } from '../service/result-subscription.service';
import { ActivatedRoute } from '@angular/router';
import { APP_BASE_HREF } from '@angular/common';
import { Status } from '../api/status';
import { TestParamsProviderService } from '../service/params-provider.service.spec';
import { Observable, NEVER } from 'rxjs';
import { TranscriptionResult } from '../api/transcription-result';

class TestUtil {
  static configure(providers: any[]) {
    TestBed.configureTestingModule({
      declarations: [ResultsComponent, StatusHumanPipe, ResultTextPipe],
      imports: [TestAppModule],
      providers: providers
    })
      .compileComponents();
  }
  static initProviders(params: ParamsProviderService): any[] {
    return [{ provide: ParamsProviderService, useValue: params },
    { provide: APP_BASE_HREF, useValue: '/' },
    { provide: TranscriptionService, useClass: MockTestService },
    { provide: ResultSubscriptionService, useClass: MockSubscriptionService },
    { provide: ActivatedRoute, useClass: MockActivatedRoute }];
  }
}

describe('ResultsComponent', () => {
  let component: ResultsComponent;
  let fixture: ComponentFixture<ResultsComponent>;

  beforeEach(async(() => {

    TestBed.configureTestingModule({
      declarations: [ResultsComponent, StatusHumanPipe, ResultTextPipe],
      imports: [TestAppModule]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ResultsComponent);
    component = fixture.debugElement.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });

  it('should have Transkripcijos ID placeholder', async(() => {
    const compiled = fixture.debugElement.nativeElement;
    expect(compiled.querySelector('input').getAttribute('placeholder')).toBe('Transkripcijos ID');
  }));

  it('should have input value from transcription id component', async(() => {
    component.transcriptionId = 'id1';
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      const input = fixture.debugElement.query(By.css('#transcriptionIDInput'));
      const el = input.nativeElement;
      expect(el.value).toBe('id1');

      el.value = 'olia';
      el.dispatchEvent(new Event('input'));
      expect(fixture.componentInstance.transcriptionId).toBe('olia');
    });
  }));
  it('should have readonly button', async(() => {
    expect(fixture.debugElement.query(By.css('#updateButton')).nativeElement.disabled).toBe(true);
  }));

  it('should have enabled button on transcription id selected', async(() => {
    expect(fixture.debugElement.query(By.css('#updateButton')).nativeElement.disabled).toBe(true);
    component.transcriptionId = 'test';
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(fixture.debugElement.query(By.css('#updateButton')).nativeElement.disabled).toBe(false);
    });
  }));

  it('should invoke update on click', async(() => {
    component.transcriptionId = 'test';
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      spyOn(component, 'refresh');
      fixture.debugElement.query(By.css('#updateButton')).nativeElement.click();
      expect(component.refresh).toHaveBeenCalled();
    });
  }));

  it('should have no progress bar', async(() => {
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(fixture.debugElement.query(By.css('#progressBar'))).toBe(null);
    });
  }));

  it('should have progress status bar', async(() => {
    component.progress = new Progress();
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(fixture.debugElement.query(By.css('#progressBar')).nativeElement).not.toBe(null);
    });
  }));

  it('should have progress status bar set from result', async(() => {
    const r = { status: 'Status', id: '1', error: '', recognizedText: '', progress: 10 };
    component.onResult(r);
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(fixture.debugElement.query(By.css('#progressBar')).nativeElement).not.toBe(null);
    });
  }));

  it('should have status set from result', async(() => {
    const r = { status: 'Status', id: '1', error: '', recognizedText: '', progress: 10 };
    component.onResult(r);
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(fixture.debugElement.query(By.css('#status')).nativeElement).not.toBe(null);
    });
  }));

  it('should have no status set from result', async(() => {
    const r = { status: Status.Completed, id: '1', error: '', recognizedText: '', progress: 10 };
    component.onResult(r);
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(fixture.debugElement.query(By.css('#status'))).toBe(null);
    });
  }));

  it('should have no progress status bar set from result', async(() => {
    const r = { status: Status.Completed, id: '1', error: '', recognizedText: '', progress: 10 };
    component.onResult(r);
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(fixture.debugElement.query(By.css('#progressBar'))).toBe(null);
    });
  }));

  it('should have no audio control on start', async(() => {
    fixture.whenStable().then(() => {
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#audioWaveDiv')))).toBe(false);
    });
  }));

  it('should have no audio control on no ID', async(() => {
    const r = { status: Status.NOT_FOUND, id: 'x', error: '', recognizedText: '', progress: 0 };
    component.onResult(r);
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#audioWaveDiv')))).toBe(false);
    });
  }));

  it('should have audio control', async(() => {
    const r = { status: Status.Completed, id: 'x', error: '', recognizedText: '', progress: 0 };
    component.onResult(r);
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#audioWaveDiv')))).toBe(true);
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#playAudioButton')))).toBe(true);
    });
  }));

  it('should have audio download button', async(() => {
    const r = { status: Status.Completed, id: 'oliaID', error: '', recognizedText: '', progress: 0 };
    component.onResult(r);
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#audioDownloadButton')))).toBe(true);
      expect(fixture.debugElement.query(By.css('#audioDownloadButton')).nativeElement.href).toContain('oliaID');
    });
  }));

  it('should hide play audio button', async(() => {
    const r = { status: Status.Completed, id: 'oliaID', error: '', recognizedText: '', progress: 0 };
    component.onResult(r);
    fixture.detectChanges();
    fixture.debugElement.query(By.css('#playAudioButton')).nativeElement.click();
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(fixture.debugElement.query(By.css('#playAudioButton'))).toBeNull();
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#stopAudioButton')))).toBe(true);
    });
  }));

  it('should have download menu', async(() => {
    const r = { status: Status.Completed, id: 'x', error: '', recognizedText: '', progress: 0 };
    component.onResult(r);
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#dwnMenu')))).toBe(true);
    });
  }));

  it('should show download buttons', async(() => {
    const r = { status: Status.Completed, id: 'x', error: '', recognizedText: '', progress: 0 };
    component.onResult(r);
    fixture.detectChanges();
    component.menuTrigger.openMenu();
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      const dfs = ['dfResult', 'dfResultFinal', 'dfLat', 'dfLatGz', 'dfN10', 'dfLatRescore', 'dfLatRescoreGz', 'dfWebVTT'];
      dfs.forEach(element => {
        expect(TestHelper.Visible(fixture.debugElement.query(By.css('#' + element)))).toBe(true);
      });
    });
  }));

  it('should call download', async(() => {
    const r = { status: Status.Completed, id: 'iddddd', error: '', recognizedText: '', progress: 0 };
    component.onResult(r);
    fixture.detectChanges();
    component.menuTrigger.openMenu();
    let arg: string;
    const dwnSpy = spyOn(component.fileKeeper, 'download').and.callFake(function (a: string) { arg = a; });
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      fixture.debugElement.query(By.css('#dfResult')).nativeElement.click();
      fixture.detectChanges();
      expect(dwnSpy).toHaveBeenCalledTimes(1);
      expect(arg).toBe('result.txt');
    });
  }));

  it('should call download - WebVTT', async(() => {
    const r = { status: Status.Completed, id: 'iddddd', error: '', recognizedText: '', progress: 0 };
    component.onResult(r);
    fixture.detectChanges();
    component.menuTrigger.openMenu();
    let arg: string;
    const dwnSpy = spyOn(component.fileKeeper, 'download').and.callFake(function (a: string) { arg = a; });
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      fixture.debugElement.query(By.css('#dfWebVTT')).nativeElement.click();
      fixture.detectChanges();
      expect(dwnSpy).toHaveBeenCalledTimes(1);
      expect(arg).toBe('webvtt.txt');
    });
  }));

  it('should call download - LatRestored', async(() => {
    const r = { status: Status.Completed, id: 'iddddd', error: '', recognizedText: '', progress: 0 };
    component.onResult(r);
    fixture.detectChanges();
    component.menuTrigger.openMenu();
    let arg: string;
    const dwnSpy = spyOn(component.fileKeeper, 'download').and.callFake(function (a: string) { arg = a; });
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      fixture.debugElement.query(By.css('#dfLatRescore')).nativeElement.click();
      fixture.detectChanges();
      expect(dwnSpy).toHaveBeenCalledTimes(1);
      expect(arg).toBe('lat.restored.txt');
    });
  }));

  it('should have hidden file buttons when in progress', async(() => {
    const r = { status: Status.Transcription, id: 'id', error: '', recognizedText: '', progress: 0 };
    component.onResult(r);
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#resultFileButton')))).toBe(false);
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#latticeFileButton')))).toBe(false);
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#latticeTxtFileButton')))).toBe(false);
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#nBestFileButton')))).toBe(false);
    });
  }));

  it('should have hidden file buttons on start', async(() => {
    fixture.whenStable().then(() => {
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#resultFileButton')))).toBe(false);
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#latticeFileButton')))).toBe(false);
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#latticeTxtFileButton')))).toBe(false);
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#nBestFileButton')))).toBe(false);
    });
  }));

  it('should have hidden editor button on start', async(() => {
    fixture.whenStable().then(() => {
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#openEditorButton')))).toBe(false);
    });
  }));

  it('should have hidden editor button when in progress', async(() => {
    const r = { status: Status.Transcription, id: 'id', error: '', recognizedText: '', progress: 0 };
    component.onResult(r);
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#openEditorButton')))).toBe(false);
    });
  }));

  it('should display editor button when completed', async(() => {
    const r = { status: Status.Completed, id: 'id', error: '', recognizedText: '', progress: 0 };
    component.onResult(r);
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(TestHelper.Visible(fixture.debugElement.query(By.css('#openEditorButton')))).toBe(true);
    });
  }));

  it('should display error', async(() => {
    const r = { status: Status.Rescore, id: 'x', error: 'olia', errorCode: ErrorCode.ServiceError, progress: 0 };
    component.onResult(r);
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(fixture.debugElement.query(By.css('#errorDiv')).nativeElement.innerText).toContain('Serviso klaida');
    });
  }));

  it('should show detailed error', async(() => {
    const r = { status: Status.Rescore, id: 'x', error: 'olia', errorCode: ErrorCode.ServiceError, progress: 0 };
    component.onResult(r);
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      for (let i = 0; i <= 5; i++) {
        fixture.debugElement.query(By.css('#errorDiv')).nativeElement.click();
      }
      fixture.detectChanges();
      fixture.whenStable().then(() => {
        expect(fixture.debugElement.query(By.css('#errorDiv')).nativeElement.innerText).toContain('olia');
      });
    });
  }));

  it('should click error details', async(() => {
    const r = { status: Status.Rescore, id: 'x', error: 'olia', errorCode: ErrorCode.ServiceError, progress: 0 };
    component.onResult(r);
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      spyOn(component, 'showErrDetails');
      fixture.debugElement.query(By.css('#errorDiv')).nativeElement.click();
      expect(component.showErrDetails).toHaveBeenCalled();
    });
  }));

  it('should click editor button', async(() => {
    const r = { status: Status.Completed, id: 'iddddd', error: '', recognizedText: '', progress: 0 };
    component.onResult(r);
    fixture.detectChanges();
    const editorSpy = spyOn(component, 'openEditor').and.callFake(function () { });
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      fixture.debugElement.query(By.css('#openEditorButton')).nativeElement.click();
      fixture.detectChanges();
      expect(editorSpy).toHaveBeenCalledTimes(1);
    });
  }));

});

describe('ResultsComponent Own Mock', () => {
  let component: ResultsComponent;
  let fixture: ComponentFixture<ResultsComponent>;

  it('should read transcription ID from route', async(() => {
    class MockActivatedRouteInternal {
      snapshot = { paramMap: new Map([['id', 'id1']]) };
    }

    const params = new TestParamsProviderService();
    params.setTranscriptionID('id2');
    let providers = TestUtil.initProviders(params);
    providers = providers.concat({ provide: ActivatedRoute, useClass: MockActivatedRouteInternal });
    TestUtil.configure(providers);

    fixture = TestBed.createComponent(ResultsComponent);
    component = fixture.debugElement.componentInstance;
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      const input = fixture.debugElement.query(By.css('#transcriptionIDInput'));
      const el = input.nativeElement;
      expect(el.value).toBe('id1');
    });
  }));

  it('should ignore provider value', async(() => {
    class MockActivatedRouteInternal {
      snapshot = { paramMap: new Map([['id', 'id1']]) };
    }

    const params = new TestParamsProviderService();
    params.setTranscriptionID('id2');
    let providers = TestUtil.initProviders(params);
    providers = providers.concat({ provide: ActivatedRoute, useClass: MockActivatedRouteInternal });
    TestUtil.configure(providers);

    fixture = TestBed.createComponent(ResultsComponent);
    component = fixture.debugElement.componentInstance;
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      const input = fixture.debugElement.query(By.css('#transcriptionIDInput'));
      const el = input.nativeElement;
      expect(el.value).toBe('id1');
    });
  }));

  it('should read value from provider', async(() => {
    const params = new TestParamsProviderService();
    params.setTranscriptionID('id1');
    TestUtil.configure(TestUtil.initProviders(params));

    fixture = TestBed.createComponent(ResultsComponent);
    component = fixture.debugElement.componentInstance;
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      const input = fixture.debugElement.query(By.css('#transcriptionIDInput'));
      const el = input.nativeElement;
      expect(el.value).toBe('id1');
    });
  }));

  it('should stop playing on destroy', async(() => {
    const params = new TestParamsProviderService();
    TestUtil.configure(TestUtil.initProviders(params));

    fixture = TestBed.createComponent(ResultsComponent);
    component = fixture.debugElement.componentInstance;
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      component.playAudio();
      expect(component.audioPlayer.isPlaying()).toEqual(true);
      fixture.destroy();
      fixture.whenStable().then(() => {
        expect(component.audioPlayer.isPlaying()).toEqual(false);
      });
    });
  }));

  it('should unsubscribe websocket on destroy', async(() => {
    const params = new TestParamsProviderService();
    class MockResultSubscriptionService implements ResultSubscriptionService {
      connect(): Observable<TranscriptionResult> {
        return NEVER;
      }
      send(id: string): void {
      }
    }

    let providers = TestUtil.initProviders(params);
    providers = providers.concat({ provide: ResultSubscriptionService, useClass: MockResultSubscriptionService });
    TestUtil.configure(providers);

    fixture = TestBed.createComponent(ResultsComponent);
    component = fixture.debugElement.componentInstance;
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(component.resultSubscription.closed).toEqual(false);
      fixture.destroy();
      fixture.whenStable().then(() => {
        expect(component.resultSubscription.closed).toEqual(true);
      });
    });
  }));
});

