#!/bin/bash
function finish {
    . ${APP_DIR}/m_end.sh
}
trap finish EXIT

shift
eval "$@"
