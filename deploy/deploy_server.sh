#!/usr/bin/env bash
set -euo pipefail

# Usage: edit DEPLOY_USER, WEB_DIR if needed, then run
# Example: sudo DEPLOY_USER=ubuntu WEB_DIR=/var/www ./deploy_server.sh

DEPLOY_USER=${DEPLOY_USER:-ubuntu}
WEB_DIR=${WEB_DIR:-/var/www}
FRONT_DIR=${FRONT_DIR:-$WEB_DIR/Construction-Frontend}
BACK_DIR=${BACK_DIR:-$WEB_DIR/Construction-Backend}
FRONT_REPO=${FRONT_REPO:-https://github.com/baidar-10/Construction-Frontend.git}
BACK_REPO=${BACK_REPO:-https://github.com/baidar-10/Construction-Backend.git}

echo "[info] deploy settings: user=$DEPLOY_USER web_dir=$WEB_DIR"

# Ensure web dir exists and is owned by deploy user
sudo mkdir -p "$WEB_DIR"
sudo chown -R "$DEPLOY_USER":"$DEPLOY_USER" "$WEB_DIR"

# Move pre-existing my_project if present
if [ -d "$WEB_DIR/my_project" ]; then
  echo "[info] moving /var/www/my_project -> $FRONT_DIR"
  sudo mv "$WEB_DIR/my_project" "$FRONT_DIR" || true
  sudo chown -R "$DEPLOY_USER":"$DEPLOY_USER" "$FRONT_DIR"
fi

# Clone or update frontend
if [ -d "$FRONT_DIR/.git" ]; then
  echo "[info] frontend exists, pulling latest"
  cd "$FRONT_DIR"
  git pull || true
else
  echo "[info] cloning frontend to $FRONT_DIR"
  cd "$WEB_DIR"
  git clone "$FRONT_REPO" Construction-Frontend
  sudo chown -R "$DEPLOY_USER":"$DEPLOY_USER" "$FRONT_DIR"
fi

# Clone or update backend
if [ -d "$BACK_DIR/.git" ]; then
  echo "[info] backend exists, pulling latest"
  cd "$BACK_DIR"
  git pull || true
else
  echo "[info] cloning backend to $BACK_DIR"
  cd "$WEB_DIR"
  git clone "$BACK_REPO" Construction-Backend
  sudo chown -R "$DEPLOY_USER":"$DEPLOY_USER" "$BACK_DIR"
fi

# Build and run via docker-compose
cd "$BACK_DIR"

if ! command -v docker-compose >/dev/null 2>&1; then
  echo "[error] docker-compose not found. Install docker-compose or use docker compose (plugin)." >&2
  exit 1
fi

echo "[info] building frontend image"
sudo docker-compose build frontend

echo "[info] starting services"
sudo docker-compose up -d

echo "[info] checking containers"
sudo docker-compose ps

echo "[info] frontend logs (last 200 lines)"
sudo docker logs construction_frontend --tail 200 || true

echo "[info] quick curl check on localhost"
if curl -I http://localhost/ 2>/dev/null; then
  echo "[ok] frontend responded on localhost"
else
  echo "[warn] frontend didn't respond on localhost (check logs)"
fi

# Open ports if ufw exists
if command -v ufw >/dev/null 2>&1; then
  echo "[info] ensuring firewall allows 80 and 443"
  sudo ufw allow 80/tcp || true
  sudo ufw allow 443/tcp || true
fi

# Remove possible PAT lines from history
if [ -f "$HOME/.bash_history" ]; then
  sed -i '/ghp_/d' "$HOME/.bash_history" || true
  history -w || true
  echo "[info] removed ghp_ tokens from bash history (if present)"
fi

echo "[done] deploy script finished. If anything failed, check the output above and logs."
