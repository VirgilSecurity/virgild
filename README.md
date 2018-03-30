[![Build Status](https://travis-ci.org/VirgilSecurity/virgild.svg?branch=master)](https://travis-ci.org/VirgilSecurity/virgild)
[![Build status](https://ci.appveyor.com/api/projects/status/pdaahgsdva7bd0w3/branch/master?svg=true)](https://ci.appveyor.com/project/tochka/virgild/branch/master)

# VirgilD

Virgil Security, Inc., is a Manassas, Virginia, based cybersecurity company and a graduate of the Fall 2014 cohort of the MACH37 cybersecurity accelerator program.

We operate a “key management in the cloud” service in combination with open sourced libraries that are available for desktop, embedded, mobile, and cloud / web applications with support for a wide variety of modern programming languages.

Our first generation cloud-based cryptography and key management system uses a centralized trust model.  To accelerate wide-scale adoption, we propose to move to a distributed trust model.

**VirgilD** is a public key management system which is fully compatible with [Virgil cloud](https://virgilsecurity.com). VirgilD can work like a caching service of public keys. It is the core of distributed key management system. It operates in a chainable manner allowing to build decentralized trust models at any scale.

VirgilD features:
- It is  100% API compatible with the Virgil Cloud
- VirgilD instances can work as a cache to the cloud, speeding up the access to your keys..
- VirgilD instances can work as a cache to other VirgilD instances, thus forming an infinite scale trusted information database
- It has a pluggable network engine architecture. Right now it supports only HTTP(S) but we will add other protocols soon

Our reference implementation is written in Go language and runs on Linux and Windows  native mode or via docker in any docker supported platform. It provides a software interface to store cryptographically validated objects as well as provide a simple validation mechanism for any data secured by the system.

VirgilD will significantly speed up the worldwide adoption of secure messaging, distributed identity management, verifiable and cryptographically protected content distribution, asset management, and many other use cases that rely on cryptographically validated security and trust.

By moving to a distributed trust model, Virgil will accelerate its ability to penetrate the market and help ensure that Virgil will be at the epicenter of Internet application security.

# Topics
* [Getting started](#getting-started)
	* [Install](#install)
		* [Hot to build](#hot-to-build)  
	* [Usage mode](#usage-mode)
		* [Cache servcie](#cache-servcie)		
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
* Cache service for cards

### Cache service

``` shell
$ ./virgild
```


# API
All information you can find on the [development portal](https://developer.virgilsecurity.com/docs/api-reference/card-service/)

# Appendix A. Environment

For using command line arguments (args) use prefix -

Arg | Environment variable name | Config variable name | Description
---|---|---|---
 address | ADDRESS | address | VirgilD address
 https-enabled | HTTPS_ENABLED | https-enabled | Enable HTTPS mode
 https-certificate | HTTPS_CERTIFICATE | https-certificate | The path of the certificate file.
 https-private-key | HTTPS_PRIVATE_KEY | https-private-key | The path of private key file.
 config | CONFIG | - | Path to config file
 logger-type | LOGGER_TYPE | logger-type | Logger type (enum: file)
 logger-file-output | LOGGER_FILE_OUTPUT | logger-file-output | Path to log file ('-' - special parameter for colsole output)
 cache-type | CACHE_TYPE | cache-type | Cache type (enum: mem)
 cache-mem-duration | CACHE_DURATION | cache-duration | Cache duration
 cache-mem-size | CACHE_SIZE | cache-size | Cache size (mb)
 card-raservice | CARD_RASERVICE | card-raservice | Addres of Registration authority
 card-raservice | CARD_CARDSSERVICE | card-cardsservice | Addres of Cards service


## Default arguments

 Arg | Value
 ---|---
 address | :8080
 https-enabled | false
 config | virgild.conf
 logger-type | file
 logger-file-output | -
 cache-type | mem
 cache-mem-duration | 1h
 cache-mem-size | 1024
 card-raservice | https://ra.virgilsecurity.com
 card-raservice | https://cards.virgilsecurity.com
 identity-service | https://identity.virgilsecurity.com
