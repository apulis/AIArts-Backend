#!/bin/bash

set -x

# 镜像名和代码分支配置
image_name=aiarts-backend
tag_name=rc0
branch_name=v1.2.0
image_fullname=${image_name}:${branch_name}-${tag_name}

# 推送到harbor
harbor_path=harbor.sigsus.cn:8443/library/apulistech/
harbor_fullname=${harbor_path}${image_fullname}
local_harbor_fullname=harbor.sigsus.cn:8443/sz_gongdianju/apulistech/dlworkspace_aiarts-backend:1.0

# 拷贝到其它机器
remote_ip=219.133.167.42
tar_name=`sed 's|:|_|g' <<< ${image_fullname}`.tar


########################################################
function prebuild(){
    git pull origin $branch_name
}

function build(){

    cd ..

    docker build . -t ${image_fullname} -f deployment/Dockerfile
    docker tag ${image_fullname} ${harbor_fullname}
}

function postbuild() {

    docker push ${harbor_fullname}
    docker save ${harbor_fullname} > ${tar_name}


    ## 开发环境
    scp -P 52080 ${tar_name} root@${remote_ip}:/tmp

    ## 导入镜像包并重新打印标签、push 本地harbor等等
    cmd="cd /tmp; docker load < ${tar_name};
             docker tag ${harbor_fullname} ${local_harbor_fullname};
             docker push ${local_harbor_fullname};
             cd /home/dlwsadmin/DLWorkspace/YTung/src/ClusterBootstrap;
             ./deploy.py kubernetes stop ${image_name};
             ./deploy.py kubernetes start ${image_name};
    "
    echo $cmd
    ssh -p 52080 root@${remote_ip} $cmd



    ## 测试环境
    #scp -P 5122 ${new_image}.tar root@10.31.3.106:/tmp
}


function main(){

    prebuild
    build
    postbuild
}

main $*


