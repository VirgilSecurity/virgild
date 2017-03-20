[![Build Status](https://travis-ci.org/VirgilSecurity/virgild.svg?branch=master)](https://travis-ci.org/VirgilSecurity/virgild)
[![Build status](https://ci.appveyor.com/api/projects/status/pdaahgsdva7bd0w3?svg=true)](https://ci.appveyor.com/project/tochka/virgild)

# VirgilD

Virgil Security, Inc., is a Manassas, Virginia, based cybersecurity company and a graduate of the Fall 2014 cohort of the MACH37 cybersecurity accelerator program.

We operate a “key management in the cloud” service in combination with open sourced libraries that are available for desktop, embedded, mobile, and cloud / web applications with support for a wide variety of modern programming languages.

Our first generation cloud-based cryptography and key management system uses a centralized trust model.  To accelerate wide-scale adoption, we propose to move to a distributed trust model.

**VirgilD** is a public key management system which is fully compatible with [Virgil cloud](https://virgilsecurity.com). VirgilD can work like a caching service of public keys and also like a standalone version. It is the core of distributed key management system. It operates in a chainable manner allowing to build decentralized trust models at any scale.

VirgilD features:
- It is  100% API compatible with the Virgil Cloud
- VirgilD instances can work as a cache to the cloud, speeding up the access to your keys..
- VirgilD instances can work as a cache to other VirgilD instances, thus forming an infinite scale trusted information database
- It has a pluggable network engine architecture. Right now it supports only HTTP(S) but we will add other protocols soon
- Local PKI

Our reference implementation is written in Go language and runs on Linux and Windows  native mode or via docker in any docker supported platform. It provides a software interface to store cryptographically validated objects as well as provide a simple validation mechanism for any data secured by the system.

VirgilD will significantly speed up the worldwide adoption of secure messaging, distributed identity management, verifiable and cryptographically protected content distribution, asset management, and many other use cases that rely on cryptographically validated security and trust.

By moving to a distributed trust model, Virgil will accelerate its ability to penetrate the market and help ensure that Virgil will be at the epicenter of Internet application security.

# Topics
* [Getting started](#getting-started)
	* [Install](#install)
		* [Hot to build](#hot-to-build)  
	* [Usage mode](#usage-mode)
		* [Cache servcie](#cache-servcie)
		* [Local PKI](#local-pki)
		* [PKI with sync mode](#pki-with-sync-mode)
	* [Check](#check)
	* [Settings](#settings)
* [API](#api)
* [Appendix A. Environment](#appendix-a-environment)
	* [Default arguments](#default-arguments)
* [Appendix B. Token based authentication](#appendix-b-token-based-authentication)
	* [Prepere](#prepere)
	* [Create token](#create-token)
	* [Get tokens](#get-tokens)
	* [Update token](#update-token)
	* [Delete token](#delete-token)

# Getting started

## Install

You can download pre-build release from the [GitHub](https://github.com/VirgilSecurity/virgild/releases) or use [docker](https://hub.docker.com/r/virgilsecurity/virgild/). Also you can compile app from source code.

### Hot to build

#### Step 1 - Get source code
``` shell
$ git clone https://github.com/VirgilSecurity/virgild
```

#### Step 2 - Build
The application support two cryptographic provider [C crypto](https://github.com/VirgilSecurity/virgil-crypto) (CGO required) and pure [Go crypto](https://gopkg.in/virgil.v4).

In first option you must install [required packages](https://github.com/VirgilSecurity/virgil-crypto#build-prerequisites) for build C crypto.

`NOTE: C crypto is not supported on Windows OS. `

Build by default
 OS       | Default crypto
 ---------|-----------------------|
 LINUX    | C crypto
 MAC OS   | C crypto
 Windows  | native crypto

``` shell
$ make
# output file virgild (virgild.exe for Windows) in root folder
```

You can manually disable use C crypto bypass C_CRYPTO=false.
``` shell
$ make C_CRYPTO=false
```


## Usage mode
Virgild can work in 3 modes.
* Cache service of global cards
* Local PKI service compatible with cloud service
* PKI service with synchronization with cloud storage

### Cache service

``` shell
$ ./virgild
```

### Local PKI
Virgild card will be generated on the first program start. All information  will be stored in *./virgild.conf*  config file so we recomended add a volume for persistence.

``` shell
$ ./virgild -mode=local
```

### PKI with sync mode
Register on [develop portal](https://developer.virgilsecurity.com) and create your application. Run app by following command where {VD_CARD_ID} and {VD_KEY} should be replaced with values from developer portal. You can use copied base64 string from developer portal or encode your private key file with base64 and supply as a command line argument

```
$ ./virgild -mode=sync -vd-card-id={CARD_ID} -vd-key={PRIVATE_KEY} -vd-key-password={PRIVATE_KEY_PASSWORD} -remote-token={REMOTE_TOKEN}
```

## Check

```
$ curl http://localhost/health/info -H 'Authorization: Basic YWRtaW46YWRtaW4=' -v
```

where [basic authentication](https://en.wikipedia.org/wiki/Basic_access_authentication) is credentials for your admin panel. For previous request used default credential (login: admin, password: admin)

## Settings
* *db* - database connection string {driver}:{connection}. Supported drivers: sqlite3, mysql, pq, mssql (by default `virgild.db`)
* *VirgilD card* - it's VirgilD card id and private key settings used in sync mode. VirgilD will sign all creating or deleting card requests which go through it. If it not set VirgilD will create card and private key on the first run and store them into local storage. You can get public key information by issuing the following curl command

```
$ curl http://localhost:8080/api/card
```

* *Authority card* - It's a card whose signature we trust. If this parameter is set up then a client's card must have signature of the authority. The parameter contains of two values: card ID card and public key
* *Auth mode* - it's authentication mode for getting access to VirgilD. It can take two values: no and local. No mode - will give you full access to VirgilD without any permissions. Local mode - provides permissions by tokens. ([Setup token based permission](#appendix-b-token-based-authentication))

Full list of parameters in [Appendix A. Environment](#appendix-a-environment).

[List of default arguments](#default-arguments)

# API
All information you can find on the [development portal](https://virgilsecurity.com/docs/services/cards/v4/cards-service)

# Appendix A. Environment

For using command line arguments (args) use prefix -

Arg | Environment variable name | Config variable name | Description
---|---|---|---
address | ADDRESS | ADDRESS | VirgilD address
 config | CONFIG | - | Path to config file
 db | DB | db |  Database connection string {driver}:{connection}. Supported drivers: sqlite3, mysql, pq, mssql
 log | LOG | log | Path to file log. 'console' is special value for print to stdout
 mode | MODE | mode | VirgilD service mode
 vd-card-id | VD_CARD_ID | vd-card-id | VirgilD card id
 vd-key | VD_KEY | vd-key | VirgilD private key
 vd-key-password | VD_KEY_PASSWORD | vd-key-password | Password for Virgild private key
 admin-enabled | ADMIN_ENABLED | admin-enabled | Enabled admin panel
 admin-login | ADMIN_LOGIN | admin_login | User name for login to admin panel
 admin-password | ADMIN_PASSWORD | admin_password | SHA256 hash of admin password
 cache | CACHE | cache | Caching duration for global cards (in seconds)
 cards-service | CARDS_SERVICE | cards-service |  Address of Cards service
 cards-ro-service | CARDS_RO_SERVICE | cards-ro-service | Address of Read only cards  service
 identity-service | IDENTITY_SERVICE | identity-service | Address of identity  service
 ra-service | RA_SERVICE | ra-service | Address of registration authority  service
 authority-card-id | AUTHIRUTY_CARD_ID | authority-card-id | Authority card id
 authority-pubkey | AUTHORITY_PUBKEY | authority-pubkey | Authority public key
 remote-token | REMOTE_TOKEN | remote-token | Token for get access to Virgil cloud
 auth-mode | AUTH_MODE | auth-mode | Authentication mode
 cache-duration | CACHE_DURATION | cache-duration | Caching duration of cards (in seconds)
 cache-size | CACHE_SIZE | cache-size | Size of cache (in megabytes)
 metrics-log-enabled | METRICS_LOG_ENABLED | metrics-log-enabled | Metrics are printing in log file
 metrics-log-interval | METRICS_LOG_INTERVAL | metrics-log-interval | Interval between flushing data to log file
 metrics-graphite-address | METRICS_GRAPHITE_ADDRESS | metrics-graphite-address | Address of graphite service where will be sending metrics (if this parameter is empty then metrics will not send)
 metrics-graphite-interval |  METRICS_GRAPHITE_INTERVAL | metrics-graphite-interval | Interval between flushing data to graphite
 metrics-graphite-prefix | METRICS_GRAPHITE_PREFIX | metrics-graphite-prefix | Prefix for VirgilD in graphite


## Default arguments

 Arg | Value
 ---|---
 address | :8080
 config | virgild.conf
 db | sqlite3:virgild.db
 log | console
 mode | cache
 admin-enabled | false
 admin-login | admin
 admin-password | admin
 cache | 3600
 cards-service | https://cards.virgilsecurity.com
 cards-ro-service | https://cards-ro.virgilsecurity.com
 identity-service | https://identity.virgilsecurity.com
 ra-service | https://ra.virgilsecurity.com
 authority-card-id | 3e29d43373348cfb373b7eae189214dc01d7237765e572db685839b64adca853
 authority-pubkey | MCowBQYDK2VwAyEAYR501kV1tUne2uOdkw4kErRRbJrc2Syaz5V1fuG+rVs=
 auth-mode | no
 cache-duration | 3600
 cache-size | 1024
 metrics-log-enabled | false
 metrics-log-interval | 1m
 metrics-graphite-interval |  1m

# Appendix B. Token based authentication

## Topic
* [Prepere](#prepere)
* [Create token](#create-token)
* [Get tokens](#get-tokens)
* [Update token](#update-token)
* [Delete token](#delete-token)

## Prepere

* Set auth-mode to local value and admin-enabled
* Restart service

_[Basic authentication](https://en.wikipedia.org/wiki/Basic_access_authentication) is credentials for your admin panel. For all following requests used default credential (login: admin, password: admin)_

## Create token

**POST /api/tokens**

``` json
{
	"permissions":{
      "get_card":true,
      "search_cards":true,
      "create_card":true,
      "revoke_card":true
   }
}
```

**Response**

``` json
{
	"token":"a707ccaabc1d2fcdad5a6cfb2487ecca7b52c53164e1ddb8ab293b0ab276391d",
    "permissions": {
    	"create_card":true,
        "get_card":true,
        "revoke_card":true,
        "search_cards":true
    }
}
```

**CURL example**
``` bash
$ curl http://localhost:8080/api/tokens -d '{"permissions":{"get_card":true,"search_cards":true,"create_card":true,"revoke_card":true}}' -H 'Authorization: Basic YWRtaW46YWRtaW4='
```

## Get tokens

**GET /api/tokens**

**Response**

``` json
[{
  "token": "a707ccaabc1d2fcdad5a6cfb2487ecca7b52c53164e1ddb8ab293b0ab276391d",
  "permissions": {
    "create_card": true,
    "get_card": true,
    "revoke_card": true,
    "search_cards": true
  }
}, ...]
```

**CURL example**
``` bash
$ curl http://localhost:8080/api/tokens -H 'Authorization: Basic YWRtaW46YWRtaW4='
```

## Update token

**PUT /api/tokens/{token_id}**

``` json
{
  "permissions": {
    "get_card": true,
    "search_cards": true,
    "create_card": true,
    "revoke_card": true
  }
}
```

**Response**

``` json
{
  "token": "{token_id}",
  "permissions": {
    "create_card": true,
    "get_card": true,
    "revoke_card": true,
    "search_cards": true
  }
}
```

**CURL example**
``` bash
$ curl -X PUT http://localhost:8080/api/tokens/a707ccaabc1d2fcdad5a6cfb2487ecca7b52c53164e1ddb8ab293b0ab276391d  -d '{"permissions":{"get_card":true,"search_cards":true,"create_card":true,"revoke_card":true}}' -H 'Authorization: Basic YWRtaW46YWRtaW4='
```
## Delete token

**DELETE /api/tokens/{token_id}**

Return status code 200 is token was removed correctly otherwise 500 (if not found return 404)

**CURL example**
``` bash
$ curl -X DELETE http://localhost:8080/api/tokens/a707ccaabc1d2fcdad5a6cfb2487ecca7b52c53164e1ddb8ab293b0ab276391d -H 'Authorization: Basic YWRtaW46YWRtaW4='
```
