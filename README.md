evil-proxy
==========

A TCP proxy that can be configured to simulate evil networks

Development State
-----------------

This code is alpha quality and not yet ready for prime time. Please feel free to
contribute, but don't try to use it yet.


Design
------

There are three primary structures. The first is the 'Packet' structure. It
represents a TCP Packet. If other connection protocols are implemented in the
future it could be changed to an interface.

The second is the 'Pipe' interface. A pipe is unidirectional and is used to
model the properties of a network , such as latency, bandwidth, etc.  The
unidirectional nature of a pipe allows for modeling of non-symmetric networks.
Any 'Pipe' implementation except for 'BasicPipe' should be built to forward
packets onto another 'Pipe' object. This will allow for composing a variety of
network simulations.

The third is the 'Connection' interface. A connection is a bidirectional and is
used to model network protocols, such as TCP. TCP will be initially implemented
by incrementally implementing parts of it using many different 'Connection's.
For example one for the three way handshake, another for the TCP window
management, etc. While this will not likely be the most performant way to
implement TCP it will allow each piece to be independently developed and tested.
Once TCP has been fully implemented profiling can be used to guide refactoring.
