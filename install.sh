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
            echo "golang-go docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin"
            ;;
        pacman)
            echo "go docker docker-compose"
            ;;
        xbps-install)
            echo "go docker docker-compose"
            ;;
    esac
}

install_apt_repo() {
  # Add Docker's official GPG key:
  sudo apt-get update
  sudo apt-get install ca-certificates curl
  sudo install -m 0755 -d /etc/apt/keyrings
  sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
  sudo chmod a+r /etc/apt/keyrings/docker.asc

  # Add the repository to Apt sources:
  echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
    $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
    sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

  sudo apt-get update
}

get_install_command() {
    case $1 in
        apt-get)
            install_apt_repo
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
    deps="go docker docker-compose tailscale"

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
                tailscale_cmd=""
                if echo "$missing_deps" | grep -q "tailscale"; then
                    tailscale_cmd="curl -fsSL https://tailscale.com/install.sh | sudo sh"
                fi


                if [ "$pkg_manager" = "apt-get" ]; then
                    print_blue "Updating package list..."
                    sudo -p "Please enter your password to update packages: " apt-get update
                fi

                print_blue "Installing packages: $(printf "%s%s%s" "$GREEN" "$packages" "$NC")"

                if ! eval "sudo $install_cmd $packages"; then
                    print_red "Error: Failed to install packages"
                    return 1
                fi

                if [ -n "$tailscale_cmd" ]; then
                    print_blue "Installing Tailscale..."
                    if ! eval "$tailscale_cmd"; then
                        print_red "Error: Failed to install Tailscale"
                        return 1
                    fi
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
    local hybr_server_url="https://github.com/$repo/releases/download/$version/hybr-server"

    if ! curl -s -L -o "$temp_dir/hybr" "$hybr_cli_url"; then
        print_red "Failed to download CLI binary"
        return 1
    fi


    if ! curl -s -L -o "$temp_dir/hybr-console" "$hybr_server_url"; then
        print_red "Failed to download Web Console binary"
        return 1
    fi

    chmod +x "$temp_dir/hybr"
    chmod +x "$temp_dir/hybr-console"
    mkdir -p "$install_dir"

    echo "TEMP DIR IS:"
    ls "$temp_dir"

    if ! sudo -p "Please enter your password to install hybr CLI: " mv "$temp_dir/hybr" "$install_dir/"; then
      print_red "Failed to install hybr CLI binary"
      rm -rf "$temp_dir"
      return 1
    fi

    if ! sudo -p "Please enter your password to install hybr Web Console: " mv "$temp_dir/hybr-console" "$install_dir/"; then
      print_red "Failed to install hybr Web Console binary"
      rm -rf "$temp_dir"
      return 1
    fi

    rm -rf "$temp_dir"

    if ! sudo "$install_dir/hybr" init --forceDefaults --no-ssl; then
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
    # figure out
    # sudo usermod -aG docker $USER && newgrp docker
    if download_and_run "v0.0.1"; then
        read -p "Do you want to start Hybr Web Console? [Y/n] " answer

        if [[ "$answer" == "" || "$answer" == "y" || "$answer" == "Y" ]]; then
          HYBR_CONSOLE_HOST="/_hybr" /usr/local/bin/hybr-console > /dev/null 2>&1 &
          
          read -p "Do you want to expose \`hybr-console: localhost:8080\` on your tailnet? [Y/n] " expose_answer
          
          if [[ "$expose_answer" == "" || "$expose_answer" == "y" || "$expose_answer" == "Y" ]]; then
            sudo tailscale serve --bg --set-path="_hybr" 8080
          fi
        fi

        echo "Done!"

        exit 0
    else
        print_red "Failed to setup application"
        exit 1
    fi
else
  print_red "Cannot Proceed"
fi
