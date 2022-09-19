![build workflow](https://github.com/loxilb-io/loxicmd/actions/workflows/build.yml/badge.svg)

## What is loxicmd

loxicmd is the command-line tool for loxilb. It is equivalent of "kubectl" for loxilb. loxicmd provides the following (currently) :

- Add/Delete/Get - service type external load-balancer 
- Get Port(interface) dump used by loxilb or its docker
- Get Connection track (TCP/UDP/ICMP/SCTP) information
- Add/Delete/Get - Qos Policies

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

