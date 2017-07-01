# GoSSLTerm

Golang reverse proxy to terminate SSL in the network namespace of a given container.
Even though it could be use genericly. 


## Example

The service `qnib/plain-httpcheck` serves a webservice on `:8080`.<br>
**Note**, that we are serving `:8081` instead, which will be provided by the reverse proxy.
```bash
$ docker service create --name http --replicas=1 -p 8081:8081 qnib/plain-httpcheck                                                                                                                      git:(master|●106
s0jgo83hsjghp0atptwellhbe
```

Now we spin up the container to join the task. In the example below we are using the container_id of the last container started for the service `http`.

```bash
$ docker run -ti --rm --network=container:$(docker ps -qlf label=com.docker.swarm.service.name=http) \
             -e GOSSLTERM_BACKEND_ADDR=127.0.0.1:8080 -e GOSSLTERM_FRONTEND_ADDR=:8081 qnib/$(basename ${PWD})
[II] qnib/init-plain script v0.4.28
> execute CMD 'gosslterm'
2017/07/01 14:56:46 Load cert '/opt/qnib/ssl/cert.pem' and key '/opt/qnib/ssl/key.pem'
2017/07/01 14:56:46 Create http.Server on ':8081'
```

Nice, now we have 

```bash
[negroni] 2017-07-01T14:56:50Z | 200 | 	 19.844831ms | 127.0.0.1:8081 | GET /pi
```
