#!/bin/bash

./bot -botNum 1000 \
    -createRoomNum 100 \
    -sendInterval 100 \
    -botIdx 0 \
    -botPrefix bot0_ \
    -botDevPrefix botdev0_ \
    -domain qa.sumscope.tptest.com \
    -url https://qm-tptest.qmhost1.com \
    -concurrent 500
