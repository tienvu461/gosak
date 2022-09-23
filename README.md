# gosak
a swiss army knife written in go
- curl ifconfig.me
- curl --insecure -vvI https://portal.auone.jp 2>&1 | awk 'BEGIN { cert=0 } /^\* SSL connection/ { cert=1 } /^\*/ { if (cert) print }'        
- nslookup
