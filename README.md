# Apulis AIArts Backend

## Quick Start
### configuration
* public datasets path: `storage/dataset/storage`
* private datasets path: `work/user/storage`
* Update docker images into local dockerhub 

  `sudo vim /etc/hosts`
  
   add 127.0.0.1  harbor.sigsus.cn
 
  ```shell 
  sudo vim /etc/docker/daemon.json
  ```
    ```
  {
    "registry-mirrors": [],
    "insecure-registries": [
     "https://harbor.sigsus.cn:8443"
    ],
    "debug": true,
    "experimental": false
  }
  ```
* restart docker and login 
  ```shell 
  sudo systemctl  restart docker 
  sudo docker login harbor.sigsus.cn:8443
  ```

 * push builded image into harbor
    ```shell 
    cd deployment
    ./build2harbor.sh
    ```
   
