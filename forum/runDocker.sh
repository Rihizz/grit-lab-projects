#!/bin/sh
# run script with "./runDocker.sh"
docker image build -f Dockerfile -t forumdock .
docker container run -p 80:80 --detach --name forum forumdock