# Roselite

Roselite is a simple application to relay your [Uptime Kuma](https://github.com/louislam/uptime-kuma) push monitor type
(falls within the passive monitor category) for multiple applications.

## Usage

Create a file using TOML, YAML, or JSON containing the configuration for Roselite. The configuration in JSON
should look like:

```json5
{
    // This "monitors" block is required
    "monitors": [
        {
            // Monitor type specifies what kind of thing you're monitoring
            // and help the program to know what's the best way of reaching the target.
            // Available monitor types are: "HTTP", "ICMP"
            "monitor_type": "HTTP",
            // This is the endpoint that you can acquire from your Uptime Kuma instance
            "push_url": "https://your-uptime-kuma.com/api/push/Eq15E23yc3",
            // This is the endpoint to your private/secluded server within an internal network
            "monitor_target": "https://your-internal-endpoint.com"
        },
        // ...
    ],
    // This "error_reporting" block is optional. It's useful to have it when you have Sentry
    // on your environment. So you can report bugs to us.
    "error_reporting": {
        "sentry_dsn": "https://***@ingest.sentry.io/***"
    },
    // This "server" block is only required if you start the Roselite as a server, not as agent.
    "server": {
        "listen_address": "127.0.0.1:8321"
    }
}
```

Put the path to the configuration file on `CONFIGURATION_FILE_PATH` environment variable. It will be read by Roselite
and be processed. It only accepts file with extensions of `toml`, `yml`, `yaml`, `json`, and `json5`.
For more details, see [conf.example.yml](./conf.example.toml) file.

That's it. Now you can run your own Roselite and looks at moving heartbeats on your Uptime Kuma instance.

## Contributing

See [contributing guide](./CONTRIBUTING.md)

## Where do "Roselite" comes from?

We have a convention of having project names being taken from
the [list of minerals on Wikipedia](https://en.wikipedia.org/wiki/List_of_minerals).
As this one happens to be some project that's probably a bit niche, I'd like to think of having it named after
our convention. You can read more about Roselite [on Wikipedia](https://en.wikipedia.org/wiki/Roselite).

## License

```
Copyright (C) 2025  Teknologi Umum <opensource@teknologiumum.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
```

See [LICENSE](./LICENSE)