#!/bin/bash

cd ..
docker build . -t aiarts:1.0 -f deployment/Dockerfile && docker tag aiarts:1.0 apulistech/aiarts:1.0
#docker push apulistech/aiarts:1.0
