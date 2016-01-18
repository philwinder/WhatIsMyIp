Use this service to return your cluster-wide IP address.

## Problem
We often struggle with attempting to determine the correct publishable IP address for a service. This occurs when operating in isolated network modes. It is also difficult to guess which network adaptor IP Address to use. The solution is to attempt to ping a service that resides within the cluster, on a different machine.

## How to use
```
$ docker run -d -p 8080:80 philwinder/whatismyip
$ curl $HOST:8080
{"Ip":"192.168.99.1","Port":"57630","ForwardedIp":""}
```
The Ip field returns the Ip address from the request. If behind a proxy, the ForwardedIp field will be present.
