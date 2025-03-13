# To Do

- [ ] CLI Improvements
    - [ ] Expose Web Service Managment APIs to the CLI
        - [x] List Services
        - [x] Info/Components
        - [x] Start/Stop
        - [x] Logs
        - [ ] Variable Edit
    - [ ] Usefull flags(verbose, accept-all)
    - [ ] Sync {initial, post install} services with a local machine

# Done

- [x] Tailscale Integration
    - [x] Yeet nginx (with all related shit)
    - [x] Tailscale Manager Service
    - [x] Update Service init/generator to have tailscale specific configs(root service, proxy path)
    - [x] Ability to change Web Console SELF_URL_PATH
    - [x] CLI actions of remote hosts

- [x] Docs
    - [x] Figure out a proper way to host docs


- [x] CLI Subcommands
    - [x] Integrate Cobra (Add `hybr version`, shell completions)
    - [x] Require root privileges for relevant subcommands

- [x] Nextcloud service templates

- [x] Global Notification System
    - [x] Global Notification Channel
    - [x] SSE for notifications
    - [x] Display alerts on actions/notifications

- [x] Service Edit
    - [x] View Vars / Edit Vars (trigger restart)
    - [x] Show Alerts on actions
    - [x] View Config(service.json) / Edit
    - [x] Stop The Service
    - [x] Restart(detect active log streaming stop, restart, continue)

- [x] Real-Time Event Orchestration
    - [x] SubscriptionManager
    - [x] Refactor real-time monitoring/log services to have a common pub/sub interface
    - [x] Add a subscription service to track Status/Component Statuses

- [x] Service Installtion & Persistennce
    - [x] Nginx Setup + Template
    - [x] Domain/Email - BaseConfig input
    - [x] SSL/Certbot stuff
    - [x] HTMX/Frontend Setup
    - [x] Persistent Service Installations
    - [x] Client-Side Service Info + Realtime Logs

# Before v0.1

- [ ] Use minified HTMX (check for local/prod?)
- [ ] Include install.sh checksum in README
