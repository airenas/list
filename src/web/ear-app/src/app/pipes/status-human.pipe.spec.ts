import { StatusHumanPipe } from './status-human.pipe';

describe('StatusHumanPipe', () => {
  it('create an instance', () => {
    const pipe = new StatusHumanPipe();
    expect(pipe).toBeTruthy();
  });

  it('transforms RECEIVED', () => {
    const pipe = new StatusHumanPipe();
    const transformed = pipe.transform('RECEIVED');
    expect(transformed).not.toEqual('RECEIVED');
  });

  it('returns the same', () => {
    const pipe = new StatusHumanPipe();
    const transformed = pipe.transform('olia');
    expect(transformed).toEqual('olia');
  });
});
