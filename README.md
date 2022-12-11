# canary-replay

### Install

Go [here](https://github.com/BRUHItsABunny/canary-replay/releases), download the binary for your OS and extract the archive.

### Usage

```
Usage:
  -multi string
        path to the directory containing the sub directories with each a request.hcy and request.json file
  -proxy string
        proto://ip:port format for proxy
  -single string
        path to directory containing request.hcy and request.json
  -timeout int
        timeout in seconds (default 30)
  -version
        This argument will print the current version data and exit
```

So if we imagine the following tree output for the exported HTTPCanary requests inside of subdirectory `a`:

```
> tree /F
Folder PATH listing
Volume serial number is XXXXXXXXX
C:.
├───1
│       request.hcy
│       request.json
│       request_body.bin
│       response.hcy
│       response.json
│       response_body.bin
│
├───2
│       request.hcy
│       request.json
│       request_body.bin
│       response.hcy
│       response.json
│       response_body.bin
│
└───3
        request.hcy
        request.json
        request_body.bin
        response.hcy
        response.json
        response_body.bin
```

In order to ONLY execute the request inside `./a/1/`:

```canaryreplay -single "a/1"```

In order to execute ALL the requests inside `./a/`:

```canaryreplay -multi "a"```

### About 

This program is used to replay exported requests from the [HTTPCanary app](https://play.google.com/store/apps/details?id=com.guoshi.httpcanary) ([premium](https://play.google.com/store/apps/details?id=com.guoshi.httpcanary.premium)) especially useful if you captured some requests that cannot be properly analyzed in HTTPCanary but can be in another proxy like Burp or Charles.