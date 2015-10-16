#!/bin/bash

APP_NAME=redis-interval-work-queue
TMP_DIR=/tmp/$APP_NAME/
IMAGE_NAME=local/$APP_NAME

init() {
  rm -rf $TMP_DIR/ \
   && mkdir -p $TMP_DIR/
}

build() {
  docker build -t $IMAGE_NAME .
}

run() {
  docker run --rm \
    --volume $TMP_DIR:/export/ \
    $IMAGE_NAME \
      cp $APP_NAME /export
}

copy() {
  cp $TMP_DIR/$APP_NAME .
}

panic() {
  local message=$1
  echo $message
  exit 1
}

main() {
  init  || panic "init failed"
  build || panic "build failed"
  run   || panic "run failed"
  copy  || panic "copy failed"
}
main
