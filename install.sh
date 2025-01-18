#!/bin/bash

clear

if [ -t 1 ]; then
    NC=$(tput sgr0)
    BOLD=$(tput bold)
    RED=$(tput setaf 1)
    BLUE=$(tput setaf 4)
    GREEN=$(tput setaf 2)
    YELLOW=$(tput setaf 3)
else
    NC=""
    RED=""
    BLUE=""
    BOLD=""
    GREEN=""
    YELLOW=""
fi

print_red() {
    printf "%s%s%s\n" "$RED" "$1" "$NC"
}

print_green() {
    printf "%s%s%s\n" "$GREEN" "$1" "$NC"
}

print_yellow() {
    printf "%s%s%s\n" "$YELLOW" "$1" "$NC"
}

print_yellow_no_newline() {
    printf "%s%s%s" "$YELLOW" "$1" "$NC"
}

print_blue() {
    printf "%s%s%s\n" "$BLUE" "$1" "$NC"
}


print_bold() {
    printf "%s%s%s\n" "$BOLD" "$1" "$NC"
}

print_prompt() {
    printf "\n${YELLOW}"
    read -r -p "$1 " answer < /dev/tty
    printf "\n\n${NC}"
}

detect_package_manager() {
    local pkg_manager=""

    if command -v apt-get >/dev/null 2>&1; then
        pkg_manager="apt-get"
    elif command -v pacman >/dev/null 2>&1; then
        pkg_manager="pacman"
    elif command -v xbps-install >/dev/null 2>&1; then
        pkg_manager="xbps-install"
    else
        print_red "No supported package manager found (only apt-get, pacman, and xbps-install are supported)"
        return 1
    fi

    echo "$pkg_manager"
    return 0
}

check_command() {
    command -v "$1" >/dev/null 2>&1
}

get_package_names() {
    case $1 in
        apt-get)
            echo "golang-go docker.io nginx"
            ;;
        pacman)
            echo "go docker nginx"
            ;;
        xbps-install)
            echo "go docker nginx"
            ;;
    esac
}

get_install_command() {
    case $1 in
        apt-get)
            echo "apt-get install -y"
            ;;
        pacman)
            echo "pacman -S --noconfirm"
            ;;
        xbps-install)
            echo "xbps-install -y"
            ;;
    esac
}

pkg_manager=$(detect_package_manager)
ret_val=$?

if [ $ret_val -ne 0 ]; then
    print_red "Error: Unable to detect package manager"
    exit 1
fi

missing_deps=""
deps="go docker nginx"

for dep in $deps; do
    if ! check_command "$dep"; then
        case "$missing_deps" in
            "") missing_deps="$dep" ;;
            *) missing_deps="$missing_deps $dep" ;;
        esac
    fi
done

if [ -n "$missing_deps" ]; then
    print_yellow "The following dependencies are missing:"
    for dep in $missing_deps; do
      print_bold "$(printf "\t")* $dep"
    done
    print_prompt "Would you like to install them? [y/N]"
    
    case $answer in
        [Yy]*)
            install_cmd=$(get_install_command "$pkg_manager")
            packages=$(get_package_names "$pkg_manager")
            
            if [ "$pkg_manager" = "apt-get" ]; then
                print_blue "Updating package list..."
                sudo apt-get update
            fi
            
            print_blue "Installing packages: $(printf "%s%s%s" "$GREEN" "$packages" "$NC")"
            if ! eval "sudo $install_cmd $packages"; then
                print_red "Error: Failed to install packages"
                exit 1
            fi
            print_green "Dependencies installed successfully"
            ;;
        *)
            print_yellow "Installation cancelled"
            exit 1
            ;;
    esac
else
    print_green "All dependencies are already installed"
fi
