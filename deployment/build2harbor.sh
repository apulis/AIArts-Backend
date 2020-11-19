#!/bin/bash

cd ..
docker build . -t aiarts:1.0 -f deployment/Dockerfile
docker tag aiarts:1.0 harbor.sigsus.cn:8443/sz_gongdianju/apulistech/dlworkspace_aiarts-backend:1.0
docker push harbor.sigsus.cn:8443/sz_gongdianju/apulistech/dlworkspace_aiarts-backend:1.0


