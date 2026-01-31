get_footer_row() {
  echo $(($(tput lines) - 2))
}

todo_header() {
  tput cup 0 0 
  echo "--- TODO n8n ---"
}

todo_footer() {
  local row=$(get_footer_row)
  tput cup "$row" 0
  echo "[q] Quit"
}

todo_menu() {
  local selected=0
  local options=("List" "Add" "Check" "Delete" "Quit")
  local message=""
  while true; do
    tput cup 3 0
    tput el
    echo -e "\033[1;37mStatus: $message\033[0m"
    
    local menu_start_row=5
    for i in "${!options[@]}"; do
      tput cup $((menu_start_row + i)) 0
      tput el
      if [[ $i -eq $selected ]]; then
        echo -e "\033[1;31m > ${options[$i]} \033[0m"
      else
        echo "  ${options[$i]}"
      fi
    done

    todo_footer

    read -rsn1 key
    if [[ $key == $'\e' ]]; then
      read -rsn2 -t 0.1 next
      [[ $next == "[A" ]] && ((selected--))
      [[ $next == "[B" ]] && ((selected++))
      (( selected < 0 )) && selected=4
      (( selected > 4 )) && selected=0
      continue
    fi

    if [[ $key == "" ]]; then
      case $selected in
        0) message="List was selected!" ;;
        1) message="Add was selected!" ;;
        2) message="Check was selected!" ;;
        3) message="Delete was selected!" ;;
        4) return ;;
      esac
    fi

    [[ $key == "q" ]] && return
  done
}

todo_start() {
  tput civis
  tput clear
  
  todo_header
  todo_menu
  
  tput clear
  tput cnorm  
}
