# backups
7574 TP1

# Build
Se utiliza una sola imagen de docker para todo, y desde el compose se cambia el command a ejecutar segun se quiera levantar un nodo o coordinador.
``` 
docker-compose up --build
```

# Demo
## WARNING: Al utilizar nc para enviar comandos, es necesario apretar CTRL+C luego de que envia los comandos, esto quiere decir que para avanzar la demo hay que apretar CTRL+C repetidas veces hasta que finalice.
``` 
./demo.sh
```

# Demo Manual
Se incluye un script para conectores por nombre a un nodo
```
# ./nc.sh <docker_network> <netcat args...>
./nc.sh backups_default server 9000 

```

