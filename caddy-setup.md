## Caddy Setup for <project-name>.lvh.me

### Overview
Use Caddy as a developer-friendly reverse proxy for local services. The `lvh.me` domain resolves to `127.0.0.1`, so any `<anything>.lvh.me` points to your machine — no hosts file changes required.

### Install
- macOS: `brew install caddy`
- Ubuntu/Debian: `sudo apt install caddy`
- Standalone: https://caddyserver.com/docs/install

### Run Your App(s)
Start your app on the desired port(s), e.g. `:3010`, `:3011`.

### Caddyfile Examples
Create a `Caddyfile` in the repo and adjust `<project-name>` and `<port>` values.

Single service (UI only)
```
<project-name>.lvh.me {
  encode gzip
  reverse_proxy 127.0.0.1:<port-app>
}
```

Multiple services (subdomains)
```
<project-name>.lvh.me {
  encode gzip
  reverse_proxy 127.0.0.1:<port-app>
}

api.<project-name>.lvh.me {
  encode gzip
  reverse_proxy 127.0.0.1:<port-api>
}

admin.<project-name>.lvh.me {
  encode gzip
  reverse_proxy 127.0.0.1:<port-admin>
}
```

Single host, path-based split
```
<project-name>.lvh.me {
  encode gzip

  handle_path /api/* {
    reverse_proxy 127.0.0.1:<port-api>
  }

  handle /* {
    reverse_proxy 127.0.0.1:<port-app>
  }
}
```

### Start Caddy
- Foreground: `caddy run --config ./Caddyfile`
- Background (system-wide): `sudo caddy start` (uses `/etc/caddy/Caddyfile`)

### Local HTTPS (optional)
Caddy can issue locally trusted certificates:
- Trust local CA: `caddy trust`
- Or force HTTP by prefixing site addresses with `http://` (e.g., `http://<project-name>.lvh.me`).

### Access
Open `https://<project-name>.lvh.me` (or `http://` if you opted out of TLS). Update subdomains/ports to match your repository services.
