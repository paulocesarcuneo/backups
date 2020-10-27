#!/bin/bash
set -x
echoTo() {
    docker exec -it $1 /go/append.sh $2 $3
}


./nc.sh backups_default server 9000 <<EOF
register,n0,/data/,client:9001,2.
register,n1,/data/,client1:9001,3.
EOF

echoTo backups_client_1 "Agregando texto" "/data/lorem"
echoTo backups_client_1 "Agregando otro texto" "/data/ipsum"


./nc.sh backups_default server 9000 <<EOF
history,n0,/data/.
history,n1,/data/.
unregister,n1,/data/.
history,n1,/data.
EOF
echoTo backups_client1_1 "Agregando texto" "/data/lorem"
echoTo backups_client1_1 "Agregando otro texto" "/data/ipsum"

./nc.sh backups_default server 9000 <<EOF
history,n0,/data/.
history,n1,/data/.
EOF
