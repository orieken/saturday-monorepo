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

# Xvfb puts the pointer at the screen's centre, under the browser windows. A
# window opening or closing under it sends that position to the page below as
# a mouse move, in the middle of a test's drag. Park it in a corner instead.
park_pointer() {
  /usr/local/bin/x-pointer.py park >/dev/null
}

# x11vnc keeps the password encoded in passwd; password.txt keeps it readable
# for the person watching (scripts/ui-suite-in-chrome.sh --watch prints it).
generate_vnc_password_file() {
  local file="$HOME/.vnc/passwd"
  mkdir -p "$(dirname "$file")"
  local password
  password="$(python3 -c 'import secrets; print(secrets.token_hex(16))')"
  x11vnc -storepasswd "$password" "$file" >/dev/null 2>&1
  (umask 077 && printf '%s\n' "$password" > "$HOME/.vnc/password.txt")
  printf '%s' "$file"
}

# RUNNER_VIEWER_BIND: where noVNC listens inside the container. 127.0.0.1 by
# default; --watch sets 0.0.0.0 so `docker run -p 127.0.0.1:6080:6080` can
# reach it, still from the host's loopback only.
# RUNNER_VIEWER_VIEW_ONLY=1: the viewer sends no mouse or keys to the display,
# which would land in a running test like any other stray input.
start_viewer() {
  local passwd_file
  passwd_file="$(generate_vnc_password_file)"
  local bind="${RUNNER_VIEWER_BIND:-127.0.0.1}"
  local view_only=()
  if [ "${RUNNER_VIEWER_VIEW_ONLY:-}" = "1" ]; then
    view_only=(-viewonly)
  fi
  x11vnc -display "$DISPLAY_NUM" -rfbport 5900 -localhost -rfbauth "$passwd_file" -forever -shared -quiet "${view_only[@]}" &
  websockify --web=/usr/share/novnc "$bind:6080" localhost:5900 &
  echo "entrypoint: RUNNER_VIEWER=1 — VNC password stored at $passwd_file (noVNC on $bind:6080)"
}

main() {
  start_xvfb
  wait_for_display
  start_window_manager
  park_pointer
  if [ "${RUNNER_VIEWER:-}" = "1" ]; then
    start_viewer
  fi
  exec "$@"
}

main "$@"
