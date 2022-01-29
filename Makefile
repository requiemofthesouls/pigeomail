BINARY_NAME=pigeomail

certs:
	openssl req -newkey rsa:2048 -sha256 -nodes -keyout .deploy/cert.key -x509 -days 365 -out .deploy/cert.pem -subj "/C=US/ST=New York/L=Brooklyn/O=Example Brooklyn Company/CN=pigeomail.ddns.net"

