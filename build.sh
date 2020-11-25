#/usr/bin/bash

image="harbor.sigsus.cn:8443/sz_gongdianju/apulistech/dlworkspace_aiarts-backend:1.0"

function checkExit() {

  if [ $? -ne 0 ]; then
    echo $1
    exit -1
  fi
}

force_build=0
publish=0
action=build


help_str="\n
   --force  force rebuild docker\n
   --dist   push image to remote harbor\n
   --help
   "

function checkArgs(){
    local TEMP=`getopt -o fhp --long force,dist,help -n error.sh -- "$@"`
    if [ $? -ne 0 ];then 
      echo "invalid arguments  --force --dist --help"
      exit -1
    fi
    eval set -- "$TEMP"
    while true
    do
        case "$1" in 
        --force) 
            force_build=1
            shift
            echo "force build"
            ;;
        --dist)
             publish=1
             shift
             ;;
        --help)
            echo -e $help_str
            shift
            exit 0
            ;;
        --)
            shift
            break
            ;;
        *)
            echo "internal error !!!"
            exit -1
            ;;
        esac
    done
}

function checkAction(){

    local prefix=${1:0:2}
    if [ "$prefix" != "" -a "$prefix" != "--" ];then
      action=$1
      shift
    fi
    echo "action is:$action"
    checkArgs $@
}

checkAction $@


if [ $force_build -eq 1 ];then
   docker build -t "$image" . -f deployment/Dockerfile
   checkExit "docker build failed !!!"
fi

if [ $publish -eq 1 ];then
   docker push "$image"
   checkExit "docker publish failed !!!"
fi
cwd=`pwd`
server_args=(
     ./AIArtsBackend
)
volumes=(
     -v "$cwd/config.yaml:/root/config.yaml" --network host
)
ports="-p 9000:9000"
case $action in
     debug)
       docker run -it --rm $ports ${volumes[@]} "$image" bash
       ;;
     run)
       docker run -d --rm $ports ${volumes[@]} "$image"  ${server_args[@]}
       ;;
     start)
       docker run  --rm $ports ${volumes[@]} "$image"  ${server_args[@]}
       ;;
     stop)                                  
       docker stop `docker ps|grep "$image"|awk {'print $1'}`
       ;;
     build)
       ;;
     *)
       docker run $ports "$image" ${server_args[@]}
       ;;
esac
