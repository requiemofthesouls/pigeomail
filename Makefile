BINARY_NAME=pigeomail
DOMAIN=pigeomail.ddns.net

certs:
	openssl req -newkey rsa:2048 -sha256 -nodes -keyout .deploy/cert.key -x509 -days 365 -out .deploy/cert.pem -subj "/C=US/ST=New York/L=Brooklyn/O=Example Brooklyn Company/CN=${DOMAIN}"
	chmod 775 .deploy/cert.key
	chmod 775 .deploy/cert.pem

