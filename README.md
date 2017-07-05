# ClickToCall #

This is simple http server that receives 2 phone numbers, token, and trying to connect this 2 numbers using Asterisk AMI Originate.

## Usage ##

1. Install Asterisk and open AMI access for container
2. Start docker container
```
docker run -d -p 9090:9090 \
    --name clicktocall --restart always \
    -e AMIHOST=127.0.0.1:5038 \
    -e AMILOGIN=admin \
    -e AMIPASSWORD=verysecret \
    -e CHANNEL=SIP/12345 \
    -e CONTEXT=outgoing-context \
    -e TOKEN=randomtoken geri4/clicktocall
```
3. Try to place call:
```
curl "127.0.0.1:9090/?token=randomtoken&phone1=89120000001&phone2=89120000002"
```