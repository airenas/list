#!/bin/bash
function finish {
    . ${SCRIPTS_DIR}/m_end.sh
}
trap finish EXIT

shift
eval "$@"
