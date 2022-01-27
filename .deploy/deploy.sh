PIGEOMAIL_VERSION="$1"
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "PIGEOMAIL_VERSION=$PIGEOMAIL_VERSION" > "$SCRIPT_DIR"/.env
docker-compose up -f docker-compose.yml -d

exit 0