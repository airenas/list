import math

from scripts.fix_short_silences import join_speaker, drop_silence
from scripts.rttm_to_seg import Seg


def test_join_speaker():
    do_test_join_speaker([Seg(start=0.1, dur=10, sp="a")], [Seg(start=0.1, dur=10, sp="a")])


def test_join_speaker_join():
    do_test_join_speaker([Seg(start=0.1, dur=10, sp="a"), Seg(start=10.1, dur=10, sp="a")],
                         [Seg(start=0.1, dur=20, sp="a")])


def test_join_speaker_join_several():
    do_test_join_speaker(
        [Seg(start=0.1, dur=10, sp="a"), Seg(start=10.1, dur=10, sp="a"), Seg(start=21.1, dur=5, sp="a")],
        [Seg(start=0.1, dur=26, sp="a")])


def test_join_speaker_np_join():
    do_test_join_speaker([Seg(start=0.1, dur=10, sp="a"), Seg(start=10.1, dur=10, sp="b")],
                         [Seg(start=0.1, dur=10, sp="a"), Seg(start=10.1, dur=10, sp="b")])

    do_test_join_speaker([Seg(start=0.1, dur=10, sp="a"), Seg(start=14.1, dur=10, sp="a")],
                         [Seg(start=0.1, dur=10, sp="a"), Seg(start=14.1, dur=10, sp="a")])

    do_test_join_speaker([Seg(start=0.1, dur=10, sp="a"), Seg(start=12.1, dur=40, sp="a")],
                         [Seg(start=0.1, dur=10, sp="a"), Seg(start=12.1, dur=40, sp="a")])


def test_drop_silence_start():
    do_drop_silence([Seg(start=0.1, dur=10, sp="a")], [Seg(start=0.0, dur=10.6, sp="a")])
    do_drop_silence([Seg(start=0.6, dur=10, sp="a")], [Seg(start=0.1, dur=11, sp="a")])


def test_drop_silence_between():
    do_drop_silence([Seg(start=0.0, dur=10, sp="a"), Seg(start=12.0, dur=10, sp="b")],
                    [Seg(start=0.0, dur=10.5, sp="a"), Seg(start=11.5, dur=11, sp="b")])
    do_drop_silence([Seg(start=0.0, dur=10, sp="a"), Seg(start=10.5, dur=10, sp="b")],
                    [Seg(start=0.0, dur=10.25, sp="a"), Seg(start=10.25, dur=10.75, sp="b")])


def test_drop_silence_end():
    # len is 60 in do_drop_silence
    do_drop_silence([Seg(start=0.0, dur=10, sp="a"), Seg(start=10.0, dur=49, sp="b")],
                    [Seg(start=0.0, dur=10, sp="a"), Seg(start=10.0, dur=49.5, sp="b")])
    do_drop_silence([Seg(start=0.0, dur=10, sp="a"), Seg(start=10.0, dur=49.75, sp="b")],
                    [Seg(start=0.0, dur=10, sp="a"), Seg(start=10.0, dur=50, sp="b")])
    do_drop_silence([Seg(start=0.0, dur=10, sp="a"), Seg(start=10.0, dur=51, sp="b")],
                    [Seg(start=0.0, dur=10, sp="a"), Seg(start=10.0, dur=51, sp="b")])


def do_test_join_speaker(in_data, expected):
    got = join_speaker(in_data)
    for expected, actual in zip(expected, got):
        assert expected.sp == actual.sp
        assert expected.start == actual.start
        assert expected.end == actual.end


def do_drop_silence(in_data, expected):
    got = drop_silence(in_data, 60)
    for expected, actual in zip(expected, got):
        assert expected.sp == actual.sp
        assert math.isclose(expected.start, actual.start, rel_tol=1e-9, abs_tol=0.0)
        assert math.isclose(expected.end, actual.end, rel_tol=1e-9, abs_tol=0.0)
