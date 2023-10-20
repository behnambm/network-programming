# port forwarding 

## TCP 

Listens on TCP port that you can specify with `-down` flag and sends all data to the address that 
is provided by `-up` flag.

```bash 
cd tcp 
go build -o app.out main.go 
```

to test if it's working:

```bash 
sudo ./app.out -down :443 -up 185.15.59.224:443 
```

the `185.15.59.224` is the IP address of `Wikipedia` 

So now you know that we want to send all received data to the Wikipedia's server and send the response from Wikipedia to 
the client.

Now the port forwarding server is up and running

let's test if it works.

change the `/etc/hosts` file that the `wikipedia` points to localhost

My `/etc/hosts`:
```text
127.0.0.1       localhost
127.0.0.1       www.wikipedia.org
127.0.0.1       wikipedia.org
127.0.0.1       *.wikipedia.org
```

and then use `curl` to send request to wikipedia:
```bash 
curl https://www.wikipedia.org
```


## UDP 

Listens on UDP/5000 and sends all data that got from client to `1.1.1.1:53`
and then returns that response form Cloudflare to the client. 

```bash
cd udp 
go run main.go 
```

to test the port forwarding:
```bash
dig -p 5000 google.com @localhost
```


