#/bin/bash
echo $(date +'%Y-%m-%d %T'): END ${tsk}
if [ -n "$tsk" ] 
then
    ${APP_DIR}/send.metric -w ${worker} -t ${tsk} -i ${id};
fi;    


