# gosak

a swiss army knife written in go

- curl ifconfig.me
- curl --insecure -vvI https://portal.auone.jp 2>&1 | awk 'BEGIN { cert=0 } /^\* SSL connection/ { cert=1 } /^\*/ { if (cert) print }'
- nslookup

# How to add new tool

- cd to project root
- Run the following command:

```
cobra-cli add new-tool-to-gosak
```
