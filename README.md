# CaddyAnyCable Module

## Overview

The CaddyAnyCable module integrates AnyCable with Caddy v2, enabling Caddy to handle WebSocket connections for [AnyCable](https://docs.anycable.io/) and proxy them to the AnyCable server. 

This is particularly useful for Ruby on Rails applications utilizing AnyCable for WebSocket connections in a production environment.

[Caddy](https://github.com/caddyserver/caddy) is a modern, open-source web server with a modular architecture that serves HTTP, HTTPS, and automatically obtains and renews TLS certificates.

[AnyCable](https://docs.anycable.io/) provides a more efficient way to handle WebSocket connections in Ruby on Rails applications, allowing you to offload WebSocket handling from your Rails application server to AnyCable-Go.


## Installation

To use this module, you must build Caddy with the AnyCable module included. This typically involves using `xcaddy` to custom build Caddy:

1. **Install `xcaddy`:**
   Follow the xcaddy [documentation](https://github.com/caddyserver/xcaddy?tab=readme-ov-file#install) for installation instructions.
2. **Compile Caddy with the `caddy-anycable` module:**
   Run the following command to build Caddy with the AnyCable module:

   ```bash
   xcaddy build --with github.com/evilmartians/caddy_anycable
   ```
   This command will compile a caddy binary file that includes the AnyCable module.
3. **Place the compiled Caddy binary and the Caddyfile in the same directory:**
   Move the compiled caddy binary file and your Caddyfile to the desired directory. 
4. **Run Caddy with the Caddyfile:**
   If there is a file called Caddyfile in the current directory and no other configuration is specified, Caddy will load the Caddyfile, adapt it, and run it right away. Execute the following command to start Caddy:
   
   ```bash
   ./caddy run
   ```
   
   If your Caddyfile is located in a different directory, you can specify the path:
   ```bash
   ./caddy run --config /path/to/Caddyfile
   ```

This process ensures that Caddy is running with the `AnyCable` module, and it will handle WebSocket connections as configured in your Caddyfile.

## Configuration

### Caddyfile Syntax

The Caddyfile syntax for configuring the AnyCable module is as follows:

```caddyfile
anycable {
    log_level <level>
    redis_url <url>
    http_broadcast_port <port>
}
```

- **log_level**: Sets the log level for AnyCable (e.g., debug, info, warn, error).
- **redis_url**: Specifies the URL of the Redis server used by AnyCable.
- **http_broadcast_port**: Defines the port for HTTP broadcasting.

Refer to the [AnyCable Configuration Documentation](https://docs.anycable.io/anycable-go/configuration) for additional configuration options.

In the anycable section of your Caddyfile, configure the settings directly corresponding to AnyCable configuration options, without the `--` prefix typically used in command-line settings.

You can enable [SSE mode](https://docs.anycable.io/anycable-go/sse) by adding the sse option to the block:

```caddyfile
   anycable {
      log_level <level>
      redis_url <url>
      http_broadcast_port <port>
      sse true
   }
```

You can also change the [WebSocket's path](https://docs.anycable.io/anycable-go/configuration?id=primary-settings) and [SSE's path](https://docs.anycable.io/anycable-go/sse?id=configuration) by setting the `path` and `sse_path` options in the block.
The default websocket path is `/cable`. The default SSE path is `/events`.

### Full example

Below is a complete example of a Caddyfile that integrates AnyCable:

```caddyfile
{
    order anycable before reverse_proxy
}

http://localhost:3000 {
    root * ./public
    @notStatic {
        not {
            file {
                try_files {path}
            }
        }
    }

    reverse_proxy @notStatic {
        to localhost:3100

        header_up X-Real-IP {remote_host}
        header_up X-Forwarded-Proto {scheme}
        header_up Access-Control-Allow-Origin *
        header_up Access-Control-Allow-Credentials true
        header_up Access-Control-Allow-Headers Cache-Control,Content-Type
        transport http {
            read_buffer 8192
        }
    }

    anycable {
        broadcast_adapter http
        presets broker
        rpc_host http://localhost:3100/_anycable
        log_level debug
    }

    file_server
}
```

This Caddyfile sets up a proxy server on port 3000, which proxies requests to the backend server on port 3100. 
It also configures an AnyCable server that listens for WebSocket connections at the `/cable` URL.


## Usage
Once configured, start Caddy with the Caddyfile. AnyCable-Go will handle WebSocket connections at the path `/cable`, and other requests will be managed according to other directives in your Caddyfile.

For further information on configuration and options, refer to the [AnyCable documentation](https://docs.anycable.io/anycable-go/configuration)
