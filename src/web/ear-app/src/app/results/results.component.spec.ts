import { TranscriptionResult } from './../api/transcription-result';
import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ResultsComponent, Progress } from './results.component';
import { TestAppModule, MockActivatedRoute, MockSubscriptionService, MockTestService } from '../base/test.app.module';
import { StatusHumanPipe } from '../pipes/status-human.pipe';
import { By } from '@angular/platform-browser';
import { ParamsProviderService } from '../service/params-provider.service';
import { TranscriptionService } from '../service/transcription.service';
import { ResultSubscriptionService } from '../service/result-subscription.service';
import { ActivatedRoute } from '@angular/router';
import { APP_BASE_HREF } from '@angular/common';

describe('ResultsComponent', () => {
  let component: ResultsComponent;
  let fixture: ComponentFixture<ResultsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ResultsComponent, StatusHumanPipe],
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
    const r = { status: 'Status', id: '1', error: '', recognizedText: '', progress: 10};

    component.onResult(r);
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(fixture.debugElement.query(By.css('#progressBar')).nativeElement).not.toBe(null);
      expect(fixture.debugElement.query(By.css('#progressBar')).nativeElement.value).toBe(10);
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

    TestBed.configureTestingModule({
      declarations: [ResultsComponent, StatusHumanPipe],
      imports: [TestAppModule],
      providers: [ParamsProviderService, { provide: APP_BASE_HREF, useValue: '/' },
        { provide: TranscriptionService, useClass: MockTestService },
        { provide: ResultSubscriptionService, useClass: MockSubscriptionService },
        { provide: ActivatedRoute, useClass: MockActivatedRouteInternal }]
    })
      .compileComponents();
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
    const params = new ParamsProviderService();
    params.lastId = 'id1';
    TestBed.configureTestingModule({
      declarations: [ResultsComponent, StatusHumanPipe],
      imports: [TestAppModule],
      providers: [{ provide: ParamsProviderService, useValue: params },
      { provide: APP_BASE_HREF, useValue: '/' },
      { provide: TranscriptionService, useClass: MockTestService },
      { provide: ResultSubscriptionService, useClass: MockSubscriptionService },
      { provide: ActivatedRoute, useClass: MockActivatedRoute }]
    })
      .compileComponents();
    fixture = TestBed.createComponent(ResultsComponent);
    component = fixture.debugElement.componentInstance;
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      const input = fixture.debugElement.query(By.css('#transcriptionIDInput'));
      const el = input.nativeElement;
      expect(el.value).toBe('id1');
    });
  }));
});

