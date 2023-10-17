import logging
import os
import sys


def prepare_logger(_logger, _ll, _use_in_az):
    _handler = logging.StreamHandler(sys.stderr)
    if not use_in_az:
        _formatter = logging.Formatter("[%(asctime)s.%(msecs)03d] %(levelname)s - %(message)s", "%Y-%m-%d %H:%M:%S")
        _handler.setFormatter(_formatter)
    _logger.handlers.append(_handler)
    _logger.propagate = False
    _logger.setLevel(level=_ll)


ll = os.environ.get('LOG_LEVEL', 'INFO').upper()
use_in_az = bool(os.environ.get('LOG_AZ', 'FALSE').upper())

logger = logging.getLogger(__name__)

prepare_logger(logger, ll, use_in_az)
