import { TestBed, async, ComponentFixture } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { AppComponent } from './app.component';
describe('AppComponent', () => {

  let fixture: ComponentFixture<AppComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        RouterTestingModule
      ],
      declarations: [
        AppComponent
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(AppComponent);
  }));
  it('should create the app', async(() => {
    const component = fixture.debugElement.componentInstance;
    expect(component).toBeTruthy();
  }));
  it(`should have navigation to upload`, async(() => {
    const compiled = fixture.debugElement.nativeElement;
    expect(compiled.querySelector('div[routerlink="/upload"]')).toBeTruthy();
  }));

  it(`should have navigation to result`, async(() => {
    const compiled = fixture.debugElement.nativeElement;
    expect(compiled.querySelector('div[routerlink="/results"]')).toBeTruthy();
  }));
});
