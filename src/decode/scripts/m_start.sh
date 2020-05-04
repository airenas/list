#/bin/bash
echo $(date +'%Y-%m-%d %T'): START $1   
if [ -n "$tsk" ] 
then
    ${APP_DIR}/send.metric -s -w ${worker} -t ${tsk} -i ${id};
fi;  
