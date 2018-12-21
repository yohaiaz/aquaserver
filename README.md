
curl -X POST localhost:12345/files -H 'Content-Type: application/json' -d '{"file": "as.docx", "size": 400}'


curl -sL https://localhost:443/files -H 'Content-Type: application/json' -d '{"file": "as.docx", "size": 400}' -k


openssl genrsa -out server.key 2048

openssl ecparam -genkey -name secp384r1 -out server.key

openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650



