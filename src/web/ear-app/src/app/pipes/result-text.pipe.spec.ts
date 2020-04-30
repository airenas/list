import { ResultTextPipe } from './result-text.pipe';

describe('ResultTextPipe', () => {
  it('create an instance', () => {
    const pipe = new ResultTextPipe();
    expect(pipe).toBeTruthy();
  });

  it('transforms NULL', () => {
    const pipe = new ResultTextPipe();
    const transformed = pipe.transform(null);
    expect(transformed).toBeNull();
  });

  it('returns the same', () => {
    const pipe = new ResultTextPipe();
    const transformed = pipe.transform('olia');
    expect(transformed).toEqual('olia');
  });

  it('changes new line symbol', () => {
    const pipe = new ResultTextPipe();
    const transformed = pipe.transform('olia\nolia');
    expect(transformed).toEqual('olia\n  olia');
  });
  it('changes several new line symbols', () => {
    const pipe = new ResultTextPipe();
    const transformed = pipe.transform('olia\nolia\nooo');
    expect(transformed).toEqual('olia\n  olia\n  ooo');
  });
});
