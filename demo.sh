#!/bin/bash
set -x
echoTo() {
    docker exec -it $1 /go/append.sh $2 $3
}


./nc.sh backups_default coordinator 9000 <<EOF
register,node1,/data/,9001,2.
register,node2,/data/,9001,3.
EOF

echoTo backups_client_1 "Agregando texto" "/data/lorem"
echoTo backups_client_1 "Agregando otro texto" "/data/ipsum"


./nc.sh backups_default coordinator 9000 <<EOF
history,node1,/data/.
history,node2,/data/.
unregister,node2,/data/.
history,node2,/data.
EOF
echoTo backups_client1_1 "Agregando texto" "/data/lorem"
echoTo backups_client1_1 "Agregando otro texto" "/data/ipsum"

./nc.sh backups_default coordinator 9000 <<EOF
history,node1,/data/.
history,node2,/data/.
EOF
