import { StatusHumanPipe } from './status-human.pipe';
import { Status } from '../api/status';

describe('StatusHumanPipe', () => {
  it('create an instance', () => {
    const pipe = new StatusHumanPipe();
    expect(pipe).toBeTruthy();
  });

  it('transforms knowns types', () => {
    const statuses: string[] = [Status.Uploaded, Status.Completed, Status.AudioConvert,
      Status.Diarization, Status.Transcription, Status.Rescore, Status.ResultMake, Status.NOT_FOUND];
    const pipe = new StatusHumanPipe();
    statuses.forEach(function (value) {
      const transformed = pipe.transform(value);
      expect(transformed).not.toEqual(value);
    });
  });

  it('transforms NOT_FOUND', () => {
    const pipe = new StatusHumanPipe();
    const transformed = pipe.transform('NOT_FOUND');
    expect(transformed).not.toEqual('NOT_FOUND');
  });

  it('returns the same', () => {
    const pipe = new StatusHumanPipe();
    const transformed = pipe.transform('olia');
    expect(transformed).toEqual('olia');
  });
});
