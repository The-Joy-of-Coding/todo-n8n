todo_header() {
  local RED=$'\e[1;31m'
  local WHITE=$'\e[1;37m'
  local logo="
${RED}████████╗  ███████╗  ██████╗   ███████╗
${RED}╚══██╔══╝ ██╔════██╗ ██╔══██╗ ██╔════██╗
${RED}   ██║    ██║    ██║ ██║  ██║ ██║    ██║
${WHITE}   ██║    ██║    ██║ ██║  ██║ ██║    ██║
${WHITE}   ██║    ╚███████╔╝ ██████╔╝ ╚███████╔╝
${WHITE}   ╚═╝     ╚══════╝  ╚═════╝   ╚══════╝
    ${RED}N 8 N  ${WHITE}E D I T I O N"
  local start_row=1
  local start_col=2
  local current_row=$start_row
  while IFS= read -r line || [[ -n "$line" ]]; do
    tput cup "$current_row" "$start_col"
    printf "%s\n" "$line"
    ((current_row++))
  done <<< "$logo"
  printf "${NC}"
}
