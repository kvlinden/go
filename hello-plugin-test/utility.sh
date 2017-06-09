#!/bin/bash

COMMAND="$1"
APP_NAME=hello-plugin
APP_SOURCE_PATH=$GOPATH\\src\\github.com\\kvlinden\\go\\$APP_NAME
TEST_SOURCE_PATH=$GOPATH\\src\\github.com\\kvlinden\\go\\$APP_NAME-test
GOCMS_PATH=$GOPATH\\src\\github.com\\gocms-io\\gocms
DESTINATION_PATH=$GOCMS_PATH\\content\\plugins\\$APP_NAME

if [ "$COMMAND" = "test-standalone" ];
        then
            echo "starting server (stand-alone)..."
            cd $APP_SOURCE_PATH
            go build
            ./hello-plugin.exe &

            echo "testing application..."
            cd $TEST_SOURCE_PATH
            go build
            ./hello-plugin-test.exe

            # kill the server in the background
            ps -ef | grep hello-plugin | grep -v grep | awk '{print $2}' | xargs kill -9
fi
if [ "$COMMAND" = "test-plugin" ];
        then
            echo "starting server (as goCMS plugin)..."
            cd $APP_SOURCE_PATH
            go build
            rm -rf $DESTINATION_PATH
            mkdir $DESTINATION_PATH
            cp $APP_NAME.exe $DESTINATION_PATH\\.
            cp manifest.json $DESTINATION_PATH\\.
            cd $GOCMS_PATH
            go build
            ./gocms.exe &

            echo "testing application..."
            cd $TEST_SOURCE_PATH
            go build
            ./hello-plugin-test.exe

            # kill the server in the background
            ps -ef | grep gocms | grep -v grep | awk '{print $2}' | xargs kill -9
fi
