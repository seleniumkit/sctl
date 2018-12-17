# SCTL
This repository contains source code for simple Selenium quota management binary.

## Building
To build the code:

1. Ensure you have [Golang](http://golang.org/) 1.11 and above.
2. Checkout this source tree: ```$ git clone https://github.com/seleniumkit/sctl.git```
3. Build as usually: ```$ go build```

## Running
## Generating quota files
Use **sctl generate** command. When run without arguments it tries to read **input.json** file in current directory and outputs XML files to the same directory:
```
$ sctl generate
```
You may want to adjust input file and output directory like the following:
```
$ sctl generate --inputFile /path/to/input.json --outputDirectory /path/to/output/directory
```
If you want to view what will be outputted without actually creating XML files use "dry run" mode:
```
$ sctl generate --dryRun
```

## Showing quota statistics
Use **sctl stat** command:
```
$ sctl stat --inputFile /path/to/input.json
```
Additionally you can only show information for one quota:
```
$ sctl stat --inputFile /path/to/input.json --quotaName test-quota
```

## Input file format
See [test-data/input.json](test-data/input.json) for full example. In the **hosts** section of the file we define a set of named host lists with regions:
```
  "hosts": {
    "cloud": {
      "region-a": {
        "selenium-cloud-a-[1:3].example.com": {
          "port": 4444,
          "count": 1
        }
      },
      "region-b": {
        "selenium-cloud-b-[1:40].example.com": {
          "ports": "444[5:8]",
          "count": 2
        }
      }
    }
  }
```
Here **cloud** is a free-form host group name, **region-a** is a free-form region (data center) name and **selenium-cloud-a-[1:20].example.com** is a short notation for a group of hosts:
```
selenium-cloud-a-1.example.com
selenium-cloud-a-2.example.com
selenium-cloud-a-3.example.com
```
In **quota** section we define quota names, browser names, their versions and use names defined in **hosts** section to refer to a group of hosts:
```
  "quota": {
    "test-quota": {
      "firefox" : {
        "defaultVersion": "33.0",
        "defaultPlatform": "LINUX",
        "versions": {
          "33.0": "cloud",
          "42.0@LINUX": "cloud"
        }
      }
    }
  }
```
To specify a platform use `@`-notation, e.g. `42.0@LINUX`. Here **test-quota** is free-form name of the quota, **firefox** is the browser name. Finally **versions** section contains a mapping of browser version to host group name, e.g. **firefox 33.0** will correspond to all hosts defined in **cloud** hosts group.
In **aliases** section we define aliases for quota blocks from **quota** section. For each defined alias quota contents will be copied to a separate file with new name.

Cloud provider attributes `username` and `password` can be included in the input file:
```
  "hosts": {
    "cloud-provider": {
      "provider-1" : {
        "cloud-provider-1.com": {
          "port": 4444,
          "count": 1,
          "username": "user1",
          "password": "Password1"
        }
      }
    }
  }
```

To specify VNC proxying settings - use `vnc` attribute as follows:
```
  "vnc-hosts": {
    "some-dc" : {
      "selenoid-host.example.com": {
        "port": 4444,
        "count": 1,
        "vnc": "selenoid"
      }
    }
  }
```
When `vnc` equals to `selenoid` then VNC will be set to `ws://host:port/vnc`, otherwise it will be set to specified value with `$hostName` placeholder is replaced by respective host name:
```
  "vnc-hosts": {
    "some-dc" : {
      "some-host-[1:5].example.com": {
        "port": 4444,
        "count": 5,
        "vnc": "vnc://$hostName:5900"
      }
    }
  }
```