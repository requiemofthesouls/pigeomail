docker build . -t ghcr.io/requiemofthesouls/gogen:latest
echo $CR_PAT | docker login ghcr.io -u requiemofthesouls --password-stdin
docker push ghcr.io/requiemofthesouls/gogen:latest