#!/usr/bin/env bash
set -euo pipefail

# guard_commit.sh — wrapper to prevent commits/pushes on protected branches
# Usage: agents/scripts/guard_commit.sh git commit -m "message"

if ! command -v git >/dev/null 2>&1; then
  echo "git is required" >&2
  exit 2
fi

branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
if [ -z "$branch" ]; then
  echo "not a git repository or detached HEAD" >&2
  exit 2
fi

case "$branch" in
  main|master)
    echo "Refusing to run git on protected branch '$branch'. Create a feature branch first." >&2
    exit 2
    ;;
  *)
    # Forward the command to git
    if [ "$#" -eq 0 ]; then
      echo "No git command provided. Example: agents/scripts/guard_commit.sh git commit -m \"msg\"" >&2
      exit 2
    fi
    git "$@"
    ;;
esac
