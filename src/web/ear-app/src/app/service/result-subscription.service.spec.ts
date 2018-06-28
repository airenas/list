import { TestBed, inject } from '@angular/core/testing';

import { ResultSubscriptionService } from './result-subscription.service';

describe('ResultSubscriptionService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [ResultSubscriptionService]
    });
  });

  it('should be created', inject([ResultSubscriptionService], (service: ResultSubscriptionService) => {
    expect(service).toBeTruthy();
  }));
});
