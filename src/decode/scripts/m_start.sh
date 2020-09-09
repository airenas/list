#/bin/bash
echo $(date +'%Y-%m-%d %T'): START $1   
if [ -n "$tsk" ] 
then
    ${SEND_METRIC_APP} -s -w ${worker} -t ${tsk} -i ${id};
fi;  
