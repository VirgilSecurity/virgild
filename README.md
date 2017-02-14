# VirgilD

Virgil Security, Inc., is a Manassas, Virginia, based cybersecurity company and a graduate of the Fall 2014 cohort of the MACH37 cybersecurity accelerator program. 

We operate a “key management in the cloud” service in combination with open sourced libraries that are available for desktop, embedded, mobile, and cloud / web applications with support for a wide variety of modern programming languages. 

Our first generation cloud-based cryptography and key management system uses a centralized trust model.  To accelerate wide-scale adoption, we propose to move to a distributed trust model.

**VirgilD** Is the core of distributed key management system. It operates in a chainable manner allowing to build decentralized trust models at any scale.
#VirgilD features:
- It is  100% API compatible with the cloud
- VirgilD instances can work as a cache to the cloud, speeding up the access to your keys..
- VirgilD instances can work as a cache to other VirgilD instances, thus forming an infinite scale trusted information database
- It has a pluggable network engine architecture. Right now it supports only HTTP(S) but we will add other protols soon

Our reference implementation is written in Go language and runs on Linux and Windows servers. It provides a software interface to store cryptographically validated objects as well as provide a simple validation mechanism for any data secured by the system.

We also achieved **~200x speed up** in basic operations when using VirgilD (Lower - better)

![Benchmark](https://habrastorage.org/files/e16/45f/dca/e1645fdcadc34feb953473622ec0c95d.png)
Time is in milliseconds


This will significantly speed up the worldwide adoption of secure messaging, distributed identity management, verifiable and cryptographically protected content distribution, asset management, and many other use cases that rely on cryptographically validated security and trust.

By moving to a distributed trust model, Virgil will accelerate its ability to penetrate the market and help ensure that Virginia will be at the epicenter of Internet application security.
