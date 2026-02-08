get_row() {
  echo $(($(tput lines) - $1))
}

todo_footer() {
  local row=$(get_row 4)
  local desc_row=$(get_row 7)
  tput cup "$desc_row" 2
  tput el
  echo -e "\e[2m${descriptions[$selected]}\e[0m"
  tput cup "$row" 2
  echo "[q] Quit"
}

search() {
  # searches records
}

update() {
  # will update the todo
}
