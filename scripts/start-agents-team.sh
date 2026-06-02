#!/bin/zsh

CMD="${1#--}"  # strip leading --, so --claude → claude, --codex → codex
if [ -z "$CMD" ]; then
  echo "Usage: $0 --claude | --codex"
  exit 1
fi

AGENTS_DIR="./agents"
TEAMMATES=""

# --- Interactive selection ---
echo "=== OmniPixel Agents Team Launcher ==="
echo ""
echo "Available teammates:"

for file in "$AGENTS_DIR"/*.md; do
  [ -e "$file" ] || continue
  file_info="- ${$(basename "$file"):r}:"$'\n'"$(<"$file")"
  TEAMMATES+="$file_info"
done

TEAMMATES=${TEAMMATES%$'\n'}
PROMPT=$'I need create an agent team to build the project in parallel. Please create teammates:\n'"${TEAMMATES}"$'\n\nRead CLAUDE.md for project context. Each teammate should read .claude/rules/ before starting any work.'

case "$CMD" in
  claude) FLAG="--dangerously-skip-permissions" ;;
  codex)  FLAG="--dangerously-bypass-approvals-and-sandbox" ;;
  *)      echo "Unknown command: $CMD"; exit 1 ;;
esac

exec "$CMD" "$FLAG" "$PROMPT"

