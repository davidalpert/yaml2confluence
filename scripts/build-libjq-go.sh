#!/bin/bash

OUT_DIR="`pwd`/out"
WORK_DIR="`pwd`/work"
JQ_DIR="$WORK_DIR/jq"
LIBJQ_DIR="$WORK_DIR/libjq"

if [ -d "$WORK_DIR" ]
then
    rm -rf "$WORK_DIR"
fi

mkdir $WORK_DIR

if [ -d "$OUT_DIR" ]
then
    rm -rf "$OUT_DIR"
fi

mkdir $OUT_DIR


git clone git@github.com:flant/libjq-go.git $LIBJQ_DIR

# checkout recommended jq commit
$LIBJQ_DIR/scripts/jq-build/checkout.sh 839316e $JQ_DIR
# apply macosx build fixes
git --git-dir $JQ_DIR/.git format-patch -1 e660003abf9bdb9f9e6959d5ebe0a536862960e7 -o $JQ_DIR/patches
git --git-dir $JQ_DIR/.git format-patch -1 77417c1335a12c4ceef469caf38c0cbfb6315b45 -o $JQ_DIR/patches
git apply --directory=$JQ_DIR --unsafe-paths $JQ_DIR/patches/*.patch 

$LIBJQ_DIR/scripts/jq-build/build-unix.sh $JQ_DIR $OUT_DIR

rm -rf "`pwd`/libjq"
mkdir "`pwd`/libjq"
cp -R $OUT_DIR/libjq/lib "`pwd`/libjq"
cp -R $OUT_DIR/libjq/include "`pwd`/libjq"
rm -rf "$WORK_DIR"
rm -rf "$OUT_DIR"