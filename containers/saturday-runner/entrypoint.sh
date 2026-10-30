#!/usr/bin/env bash
# Starts the virtual display and window manager, optionally the VNC viewer,
# then execs the given command so the container's exit code is the
# command's own exit code. No process supervisor — see docs/runner-image.md.
set -euo pipefail

DISPLAY_NUM="${DISPLAY:-:99}"
XDPYINFO_TIMEOUT_SECONDS=15

start_xvfb() {
  Xvfb "$DISPLAY_NUM" -screen 0 1920x1080x24 -nolisten tcp &
}

wait_for_display() {
  local waited=0
  until xdpyinfo -display "$DISPLAY_NUM" >/dev/null 2>&1; do
    waited=$((waited + 1))
    if [ "$waited" -ge "$XDPYINFO_TIMEOUT_SECONDS" ]; then
      echo "entrypoint: Xvfb did not become ready on $DISPLAY_NUM within ${XDPYINFO_TIMEOUT_SECONDS}s" >&2
      exit 1
    fi
    sleep 1
  done
}

start_window_manager() {
  openbox &
}

generate_vnc_password_file() {
  local file="$HOME/.vnc/passwd"
  mkdir -p "$(dirname "$file")"
  local password
  password="$(python3 -c 'import secrets; print(secrets.token_hex(16))')"
  x11vnc -storepasswd "$password" "$file" >/dev/null 2>&1
  printf '%s' "$file"
}

start_viewer() {
  local passwd_file
  passwd_file="$(generate_vnc_password_file)"
  x11vnc -display "$DISPLAY_NUM" -rfbport 5900 -localhost -rfbauth "$passwd_file" -forever -shared -quiet &
  websockify --web=/usr/share/novnc 127.0.0.1:6080 localhost:5900 &
  echo "entrypoint: RUNNER_VIEWER=1 — VNC password stored at $passwd_file (noVNC at http://127.0.0.1:6080)"
}

main() {
  start_xvfb
  wait_for_display
  start_window_manager
  if [ "${RUNNER_VIEWER:-}" = "1" ]; then
    start_viewer
  fi
  exec "$@"
}

main "$@"
