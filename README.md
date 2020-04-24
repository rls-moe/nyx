# nyx - Imageboard

A simple, dependency-free image board.

![screenshot](/screenshot.png)

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

## Hostnames

Nyx seperates hostnames as distinct namespaces with their own content.

At the moment you cannot port content between namespaces.

## Administration

The administration panel is available under `/admin/` (don't forget the trailing slash), the default login is `admin` with password `admin`. It is recommended to add a new administrator and delete the default ID.

Here you can add boards, set board rules and start a database cleanup.

The cleanup will remove entries older than 7 days, deleted and orphaned threads and replies.

Once logged in as administrator you can also delete posts on the setup boards
or mark your own posts as special (though still anonymous)

## Posting

Posts are limited to 10k characters and uploads to 10MB including Base64 overhead (realistically you should be able to upload a 7MB file)

Nyx includes a system called "Trollthrottle". It will rate content based on how well it compresses, it's length and number of lines and the occurence of (currently fixed) keywords.

The end result is a spam score and a captcha probability, both displayed along
each post. The captcha probability specifies how like it is that a user
will fail a captcha despite having entered it correctly. This is capped at 99%, which means only 1 out of 100 correct solutions will be accepted.

This systems does not stop all trolls but will make it harder for people to post spam by forcing them to do more work.

## TripCodes

Tripcodes are non-traditional, they are calculated as the first 8 bytes of the
Blake2b Hashsum of the entered Code in Base64 Encoding.

Tripcodes do not offer a guarantee that a user is who they say they are as the
codes can be trivially cracked even on a mobile device.

## Configuration

The configuration file is written in YAML.

The following is a list of options available (and supported);

* `secret` - Secret used for User Login, CSRF and Session Management, default is `changeme`
* `listen_on` - Defaults to `:8080`, specifies on which port the HTTP server is launched. Nyx will not utilize this value otherwise so proxying is safe.
* `hosts` - A whitelists of hostnames that are allowed to be used. Nyx uses hostnames to differentiate several board collections.
* `db.file` - File to use for data storage, defaults to `:memory:` which means in-memory storage
* `site.title` - Site Title
* `site.description` - Site Description
* `captcha` - Captcha Mode, currently only `internal` is supported

The config accepts other options but these may not be supported.

## Asked Questions

* **Where is the demo?**

* It's available here: https://nyx.demo.bjphoster.com (it resets every 6 hours, it's just to "try and buy" the software, )

* **Where and how is data stored?**

* Nyx uses a KV-Stored called BuntDB. It features in-memory storage and JSON-indexing, Nyx currently only takes advantage of the in-memory storage however. It enables Nyx to operate in "demo mode", leaving no trace on a computer if it's not explicitly configured. Inside the KV Store all data is encoded with JSON.

* **Why no javascript?**

* The lack of Javascript makes Nyx a simple webservice. No public API or complicated authentication. Just Vanilla Go Templates and clean endpoints. I also suck at JS so that helped.

* **Why make an image board?**

* One, there isn't much software out there for this, Wakaba and similar software runs mostly on perl and not on the shared hosting service I use. Nyx is designed to run everywhere without dependencies and without configuration. Secondly, I wanted to see if I could make an imageboard over a weekend, which I succeeded in doing (as you can see).

* **How can I migrate/convert/upgrade the database?**

* The database doesn't migrate old entries, I think that for this application it's more important that new threads and replies work, old data will simply be parse with fallbacks and phased out over time. Currently there is also no way to export or migrate the database, even just between hostnames.

## ToDo

* Moderation needs to be worked on, currently only global administrators work
* Video Files might be implemented, though they will most likely require some extra work on the database side, as streaming video from JSON encoded entries will be anything but efficient. The first step would be to enable Nyx to stream content without loading it in memory.
* MSGPACK encoding instead of JSON
* More efficient KV usage, images and thumbnails should be moved into seperate entries so that decoding thread and reply data is less costly.
* Implement Host-Whitelisting more efficiently
* Manual Hostname Whitelisting (only enable admin panels globally)
* Global links (kinda like WordPress: "home URL" to be hosted in a subdirectory path, and "board root URL" to "return home" from inside a board)
* Admin panel refactoring
    * Password change, instead of deleting an entry and re-creating with a different password
    * Board deletion
    * Modern-ish admin panel restyling
* Database Tool for the following functions:
    * Online & Offline Management
    * Manage Administrators and Moderators
    * Manage Hostnames
    * Manage Boards, Threads and Replies
    * Export / Dump Database
    * Import / Upgrade Database
