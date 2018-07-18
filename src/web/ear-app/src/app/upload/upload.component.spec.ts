import { HttpTranscriptionService, TranscriptionService } from './../service/transcription.service';
import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UploadComponent } from './upload.component';
import { TestAppModule, FileHelper, MockActivatedRoute, MockSubscriptionService, MockTestService } from '../base/test.app.module';
import { FileSizeModule } from 'ngx-filesize';
import { RouterTestingModule } from '@angular/router/testing';
import { By } from '@angular/platform-browser';
import { ParamsProviderService } from '../service/params-provider.service';
import { ResultsComponent } from '../results/results.component';
import { APP_BASE_HREF } from '@angular/common';
import { ResultSubscriptionService } from '../service/result-subscription.service';
import { ActivatedRoute } from '@angular/router';
import { StatusHumanPipe } from '../pipes/status-human.pipe';


describe('UploadComponent', () => {
  let component: UploadComponent;
  let fixture: ComponentFixture<UploadComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [UploadComponent],
      imports: [TestAppModule, FileSizeModule, RouterTestingModule.withRoutes([])]
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

  it('should have enabled button on File selected', async(() => {
    expect(fixture.debugElement.query(By.css('#uploadButton')).nativeElement.disabled).toBe(true);
    component.fileChange(new FileHelper().createFakeFile());
    fixture.debugElement.query(By.css('#hiddenFileInput')).nativeElement.dispatchEvent(new Event('input'));
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(fixture.debugElement.query(By.css('#uploadButton')).nativeElement.disabled).toBe(false);
    });
  }));

  it('should invoke upload on click', async(() => {
    component.fileChange(new FileHelper().createFakeFile());
    fixture.debugElement.query(By.css('#hiddenFileInput')).nativeElement.dispatchEvent(new Event('input'));
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      spyOn(component, 'upload');
      fixture.debugElement.query(By.css('#uploadButton')).nativeElement.click();
      expect(component.upload).toHaveBeenCalled();
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

  it('should read value from provider', async(() => {
    const params = new ParamsProviderService();
    params.lastSelectedFile = new FileHelper().createFakeFile();
    TestBed.configureTestingModule({
      declarations: [ UploadComponent ],
      imports: [TestAppModule, FileSizeModule, RouterTestingModule.withRoutes([])],
      providers: [{provide: ParamsProviderService, useValue: params},
      { provide: APP_BASE_HREF, useValue: '/' },
      { provide: TranscriptionService, useClass: MockTestService },
      { provide: ResultSubscriptionService, useClass: MockSubscriptionService },
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
});
