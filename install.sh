#!/usr/bin/env bash

set -e

function getDistributionTag {
    local ext=""
    if [[ `uname -m` == "386" ]]; then
        machine="386"
    else
        machine="amd64"
    fi

    if [[ `uname -s` == "Darwin" ]]; then
        kernel="darwin"
    elif [[ `uname -s` =~ CYGWIN*|MINGW32*|MSYS* ]]; then
        kernel="windows"
        ext=".exe"
    else
        kernel="linux"
    fi

    echo "captain_${kernel}_${machine}${ext}"
}

CAPTAIN_DIR=$HOME/.captain
CAPTAIN_BIN_DIR=$CAPTAIN_DIR/bin
CAPTAIN_BINARIES_DIR=$CAPTAIN_DIR/binaries
CAPTAIN_CURRENT_VERSION_URL=$(curl -sS https://raw.githubusercontent.com/harbur/captain/master/VERSION)
CAPTAIN_CURRENT_VERSION_PATH="${CAPTAIN_BINARIES_DIR}/captain-${CAPTAIN_CURRENT_VERSION_URL}"
CAPTAIN_DISTRIBUTION=$(getDistributionTag)


echo "Creating folders in ${CAPTAIN_DIR}"
mkdir -p $CAPTAIN_BIN_DIR $CAPTAIN_BINARIES_DIR

echo "Start downloading Captain ${CAPTAIN_CURRENT_VERSION_URL}"
curl -sSL https://github.com/harbur/captain/releases/download/${CAPTAIN_CURRENT_VERSION_URL}/${CAPTAIN_DISTRIBUTION} > ${CAPTAIN_CURRENT_VERSION_PATH}
ln -snf ${CAPTAIN_CURRENT_VERSION_PATH} "${CAPTAIN_BIN_DIR}/captain"
chmod +x "${CAPTAIN_BIN_DIR}/captain"

echo "Captain ${CAPTAIN_CURRENT_VERSION_URL} installed"
echo ""
echo "IMPORTANT: Add ${CAPTAIN_BIN_DIR} in your path. E.g.:"
echo "export PATH=${CAPTAIN_BIN_DIR}:\$PATH"

