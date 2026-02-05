get_footer_row() {
  echo $(($(tput lines) - 2))
}

todo_footer() {
  local row=$(get_footer_row)
  tput cup "$row" 0
  echo "[q] Quit"
}
