import { Status } from '../api/status';
import { ErrorPipe } from './error.pipe';
import { ErrorCode } from '../api/error-codes';

describe('ErrorPipe', () => {
  it('create an instance', () => {
    const pipe = new ErrorPipe(false);
    expect(pipe).toBeTruthy();
  });

  it('transforms to short audio errror', () => {
    const pipe = new ErrorPipe(false);
    const transformed = pipe.transform({id: 'id', status: Status.Uploaded, errorCode: ErrorCode.TooShortAudio, error: 'olia'});
    expect(transformed).toEqual('Per trumpas įrašas');
  });

  it('transforms Wrong format error', () => {
    const pipe = new ErrorPipe(false);
    const transformed = pipe.transform({id: 'id', status: Status.Uploaded, errorCode: ErrorCode.WrongFormat, error: 'olia'});
    expect(transformed).toEqual('Blogas formatas');
  });

  it('returns the same', () => {
    const pipe = new ErrorPipe(true);
    const transformed = pipe.transform({id: 'id', status: Status.Uploaded, errorCode: ErrorCode.ServiceError, error: 'olia'});
    expect(transformed).toEqual('olia');
  });

  it('hide service error', () => {
    const pipe = new ErrorPipe(false);
    const transformed = pipe.transform({id: 'id', status: Status.Uploaded, errorCode: ErrorCode.ServiceError, error: 'olia'});
    expect(transformed).toEqual('Serviso klaida');
  });
});
