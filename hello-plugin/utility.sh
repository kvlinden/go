#!/bin/bash

COMMAND="$1"
APP_NAME=hello-plugin
SOURCE_PATH=$GOPATH\\src\\github.com\\kvlinden\\go\\$APP_NAME
GOCMS_PATH=$GOPATH\\src\\github.com\\gocms-io\\gocms
DESTINATION_PATH=$GOCMS_PATH\\content\\plugins\\$APP_NAME

if [ "$COMMAND" = "deploy" ];
        then
            echo "(re)building application..."
            cd $SOURCE_PATH
            go build

            echo "(re)deploying application..."
            rm -rf $DESTINATION_PATH
            mkdir $DESTINATION_PATH
            cp $APP_NAME.exe $DESTINATION_PATH\\.
            cp manifest.json $DESTINATION_PATH\\.

            echo "running application..."
            cd $GOCMS_PATH
            go build
            ./gocms.exe
fi
