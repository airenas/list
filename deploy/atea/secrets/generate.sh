#/bin/bash
cp secrets.yml.template secrets.yml
### rabbitmq
user=list
user64="$(echo -n $user | base64)"
echo "user=$user : $user64 "
sed -i "s/{{rabbit-mq-user}}/$user64/" secrets.yml

pass="$(pwgen 20 1)"
pass64="$(echo -n $pass | base64)"
echo "user=$pass : $pass64 "

sed -i "s/{{rabbit-mq-pass}}/$pass64/" secrets.yml

### mongo
user=list
user64="$(echo -n $user | base64)"
echo "user=$user : $user64 "
sed -i "s/{{mongo-user}}/$user64/" secrets.yml

pass="$(pwgen 20 1)"
pass64="$(echo -n $pass | base64)"
echo "user=$pass : $pass64 "

sed -i "s/{{mongo-pass}}/$pass64/" secrets.yml  