todo_menu() {
  selected=0
  options=("List" "Add" "Check" "Delete" "Settings" "Quit")
  descriptions=(
    "Fetches and displays your current tasks from n8n."
    "Add a new task title, priority, and due date."
    "Toggle a task status from 'pending' to 'done'."
    "Remove a task that is no longer relevant."
    "Modify n8n workflows, logic, and integrations."
    "Exit the Todo N8N Edition TUI."
  )
  while true; do
    todo_header
    todo_footer
    select_menu
    read_keys
    perform_action
    [[ $key == "q" ]] && return
  done
}

select_menu() {
  local menu_start_row=$(get_row 16)
  for i in "${!options[@]}"; do
    tput cup $((menu_start_row + i)) 2
    tput el
    if [[ $i -eq $selected ]]; then
      echo -e "\033[1;31m > ${options[$i]} \033[0m"
    else
      echo "   ${options[$i]}"
    fi
  done
}

read_keys() {
  read -rsn1 key
  if [[ $key == $'\e' ]]; then
    read -rsn2 -t 0.1 next
    [[ $next == "[A" ]] && ((selected--))
    [[ $next == "[B" ]] && ((selected++))
    (( selected < 0 )) && selected=$length
    (( selected > "$((${#options[@]} - 1))" )) && selected=0
  fi
}

perform_action() {
  if [[ $key == "" ]]; then
    case $selected in
      0) get_todos ;;
      1) add_todo ;;
      2) check_todo ;;
      3) delete_todo ;;
      4) settings_todo ;;
      5) key="q" ;;
    esac
  fi
}

tab_switch() {
  # will switch mode between input, outputs and menu.
}
