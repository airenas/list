import { TestBed, inject } from '@angular/core/testing';

import { ParamsProviderService } from './params-provider.service';

describe('ParamsProviderService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [ParamsProviderService]
    });
  });

  it('should be created', inject([ParamsProviderService], (service: ParamsProviderService) => {
    expect(service).toBeTruthy();
  }));
});
