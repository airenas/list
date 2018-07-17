import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ResultsComponent } from './results.component';
import { TestAppModule } from '../base/test.app.module';
import { StatusHumanPipe } from '../pipes/status-human.pipe';
import { By } from '@angular/platform-browser';

describe('ResultsComponent', () => {
  let component: ResultsComponent;
  let fixture: ComponentFixture<ResultsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ResultsComponent, StatusHumanPipe ],
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
      const input = fixture.debugElement.query(By.css('input'));
      const el = input.nativeElement;
      expect(el.value).toBe('id1');

      el.value = 'olia';
      el.dispatchEvent(new Event('input'));
      expect(fixture.componentInstance.transcriptionId).toBe('olia');
    });
  }));
});
