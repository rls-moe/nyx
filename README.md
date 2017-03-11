# nyx - Imageboard

A simple, dependency-free image board.

## Requirements

* HTTPS Capable Host
* Being able to compile and run a go binary
* Disk Space

## Installation

Simply `go get` the source code or compile
it otherwise, dependencies are included.

## Usage

By default nyx works in volatile-mode,
all changes are only stored in memory and
default credentials are setup in the 
database.

To overwrite these defaults, simply
create the file `config.yml` or specify
another file via the `-config` flag.

## Configuration

The configuration file is written in YAML.

The following is a list of options available;

* `disable_security` - Disabled some HTTPS only options for cookies or redirects
* `secret` - Secret used for User Login, CSRF and Session Management, default is `changeme`
* `listen_on` - Defaults to `:8080`, specifies on which port the HTTP server is launched. Nyx will not utilize this value otherwise so proxying is safe.
* `hosts` - A whitelists of hostnames that are allowed to be used. Nyx uses hostnames to differentiate several board collections.
* `db.file` - File to use for data storage, defaults to `:memory:` which means in-memory storage
* `site.title` - Site Title
* `site.description` - Site Description