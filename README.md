# Hybr

Take control of your digital world with `hybr` - a self-hosted infrastructure manager that lets you deploy and manage services with ease.

![Hybr Progress](https://github.com/rasjonell/hybr/blob/master/hybr.png)

> [!NOTE]
> 🚧 This project is under active development. Everything including this document may be changed in the future.
> Star the repository to follow our progress!

## Features

- One-command deployment of popular self-hosted services
- Web UI for monitoring and management
- Secure VPN connections using Tailscale for private network access
- Container-based isolation for all services
- Real-time service monitoring and logs
- Easy service updates and maintenance
- Comes with default services like **tt-rss** and **Nextcloud**, but is easily extendable with your own custom services.

## Quick Start

```bash
curl -sSL https://hybr.dev/install.sh | bash
```

or with `wget`

```bash
wget -qO- https://hybr.dev/install.sh | bash
```

After installation, access the web UI at http://localhost:8080 or on your tailscale network.

## Documentation

Full documentation is available at [hybr.dev/docs](https://hybr.dev/docs/intro)

## Services

Hybr comes with the following default services but can be extended with other services(refer to Adding Custom Services)
- Tiny Tiny RSS (tt-rss): *Self-hosted news feed aggregator*
- Nextcloud: *Personal Cloud Provider*

## Adding Custom Services

Hybr uses a simple service definition, so you can add any service you want. To add a new service, you need to:

1.  Create a new directory for your service in the service templates directory.
2.  Create a `docker-compose.yml` file and any other necessary configuration files for your service.
3.  Update the `service.json` file that describes the services to Hybr. See `/internal/services/templates/services.json` for more details.
4.  Run `hybr init`


### CLI Commands

The `hybr` CLI provides a set of commands for managing your services. Here's a quick overview:

- `hybr init`: Initiates a new hybr project, allowing you to select services and configure them.
    - `-f, --forceDefaults`: Use the default templates (Optional)
    - `-a, --ts-auth`: Tailscale AUTH_KEY (Optional)
- `hybr services`: Shows services info.
    - `-s, --service`: Name of the service (Optional)
- `hybr services components`: Shows docker components the services is composed of.
    - `-s, --service`: Name of the service (Required)
- `hybr services info`: Shows service information.
    - `-s, --service`: Name of the service (Required)
- `hybr services logs`: Shows docker compose logs for the service.
    - `-s, --service`: Name of the service (Required)
- `hybr services start`: Starts the service.
    - `-s, --service`: Name of the service (Required)
- `hybr services stop`: Stop the service.
    - `-s, --service`: Name of the service (Required)
- `hybr version`: Print the version number of hybr

---

## Project Status

Check out my TODOs and past progress here: [TODO.md](https://github.com/rasjonell/hybr/blob/master/TODO.md)

[Latest Release Changelog](https://github.com/rasjonell/hybr/releases/latest)

Made with ❤️ 
