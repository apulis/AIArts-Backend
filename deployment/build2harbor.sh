#!/bin/bash

cd ..
sudo docker build . -t aiarts:1.0 -f deployment/Dockerfile
sudo docker tag aiarts:1.0 harbor.sigsus.cn:8443/sz_gongdianju/apulistech/dlworkspace_aiarts-backend:1.0
sudo docker push harbor.sigsus.cn:8443/sz_gongdianju/apulistech/dlworkspace_aiarts-backend:1.0


