## What is loxicmd

loxicmd is command tools for loxilb. loxicmd provide the following :

- Add/Delete/Get about the service type external load-balancer 
- Get Port(interface) dump
- Get Connection track information

loxicmd aim to provide all of the configuation for the loxilb.

## How to build

1. Install package dependencies 

```
go get .
```

2. Make loxicmd

```
make
```

## How to run

1. Run loxicmd with getting lb information

```
./loxicmd get lb
```

2. Run loxicmd with getting lb information in the different API server(ex. 192.168.18.10) and ports(ex. 8099).

```
./loxicmd get lb -s 192.168.18.10 -p 8099
```

3. Run loxicmd with  getting lb information as json output format
```
./loxicmd get lb -o json
```

More information use help option!
```
./loxicmd help
```

