#/bin/bash
echo $(date +'%Y-%m-%d %T'): END ${tsk}
if [ -n "$tsk" ] 
then
    ${SEND_METRIC_APP} -w ${worker} -t ${tsk} -i ${id};
fi;    


