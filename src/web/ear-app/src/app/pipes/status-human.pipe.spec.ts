import { StatusHumanPipe } from './status-human.pipe';

describe('StatusHumanPipe', () => {
  it('create an instance', () => {
    const pipe = new StatusHumanPipe();
    expect(pipe).toBeTruthy();
  });

  it('transforms ADDED', () => {
    const pipe = new StatusHumanPipe();
    const transformed = pipe.transform('ADDED');
    expect(transformed).not.toEqual('ADDED');
  });

  it('returns the same', () => {
    const pipe = new StatusHumanPipe();
    const transformed = pipe.transform('olia');
    expect(transformed).toEqual('olia');
  });
});
