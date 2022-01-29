PIGEOMAIL_VERSION="$1"
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "PIGEOMAIL_VERSION=$PIGEOMAIL_VERSION" > "$SCRIPT_DIR"/.env
docker-compose -f "$SCRIPT_DIR"/.deploy/docker-compose.yml -f "$SCRIPT_DIR"/.deploy/docker-compose.prod.yml pull
docker-compose -f "$SCRIPT_DIR"/.deploy/docker-compose.yml -f "$SCRIPT_DIR"/.deploy/docker-compose.prod.yml up -d --no-build

exit 0