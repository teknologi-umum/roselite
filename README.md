# Roselite

Roselite is a simple application to relay your [Uptime Kuma](https://github.com/louislam/uptime-kuma) push monitor type
(falls within the passive monitor category) for multiple applications.

## Usage

Create a file using TOML, YAML, or JSON (or JSON5) containing the configuration for Roselite. The configuration in JSON
should looks like:

```json5
{
    "monitors": [
        {
            // This is the endpoint that you can acquire from your Uptime Kuma instance
            "push_url": "https://your-uptime-kuma.com/api/push/Eq15E23yc3",
            // This is the endpoint to your private/secluded server within an internal network
            "monitor_url": "https://your-internal-endpoint.com"
        },
        // ...
    ]
}
```

Put the path to the configuration file on `CONFIGURATION_FILE_PATH` environment variable. It will be read by Roselite
and be processed. It only accepts file with extensions of `toml`, `yml`, `yaml`, `json`, and `json5`.

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
    MIT License
    
    Copyright (c) 2023 Teknologi Umum <opensource@teknologiumum.com>
    
    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:
    
    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.
    
    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.
```

See [LICENSE](./LICENSE)