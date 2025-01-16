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

check_and_install_dependencies() {
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
                    sudo -p "Please enter your password to update packages: " apt-get update
                fi

                print_blue "Installing packages: $(printf "%s%s%s" "$GREEN" "$packages" "$NC")"
                if ! eval "sudo $install_cmd $packages"; then
                    print_red "Error: Failed to install packages"
                    return 1
                fi
                print_green "Dependencies installed successfully"
                return 0
                ;;
            *)
                print_yellow "Installation cancelled"
                return 1
                ;;
        esac
    else
        print_green "All dependencies are already installed"
        return 0
    fi
}

get_latest_release() {
    local api_response
    local repo="rasjonell/hybr"

    print_blue "Fetching latest release version..."
    api_response=$(curl -s "https://api.github.com/repos/$repo/releases/latest")

    if [ $? -ne 0 ]; then
        print_red "Failed to fetch release information"
        return 1
    fi

    echo "$api_response" | grep -Po '"tag_name": "\K.*?(?=")'
}

download_and_run() {
    local repo="rasjonell/hybr"
    local version="$1"
    local install_dir="$2"

    if [ -z "$version" ]; then
        version=$(get_latest_release "$repo")
        if [ $? -ne 0 ]; then
            print_red "Failed to get latest version"
            return 1
        fi
    fi

    if [ -z "$install_dir" ]; then
        install_dir="/usr/local/bin"
        if [ $? -ne 0 ]; then
            print_red "Failed to get latest version"
            return 1
        fi
    fi

    print_blue "Using version: $version"
    print_blue "Downloading & installing hybr CLI..."

    local temp_dir=$(mktemp -d)
    local hybr_cli_url="https://github.com/$repo/releases/download/$version/hybr"

    if ! curl -s -L -o "$temp_dir/hybr" "$hybr_cli_url"; then
        print_red "Failed to download CLI binary"
        return 1
    fi

    chmod +x "$temp_dir/hybr"
    mkdir -p "$install_dir"

    if ! sudo -p "Please enter your password to install hybr: " mv "$temp_dir/hybr" "$install_dir/"; then
      print_red "Failed to install hybr binary"
      rm -rf "$temp_dir"
      return 1
    fi

    rm -rf "temp_dir"

    if ! "$install_dir/hybr"; then
        print_red "hybr CLI application failed to run"
        return 1
    fi
}

pkg_manager=$(detect_package_manager)
ret_val=$?

if [ $ret_val -ne 0 ]; then
    print_red "Error: Unable to detect package manager"
    exit 1
fi

if check_and_install_dependencies; then
    if download_and_run "0.0.0"; then
        exit 0
    else
        print_red "Failed to setup application"
        exit 1
    fi
else
  print_red "Cannot Proceed"
fi
