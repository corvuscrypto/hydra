# Data Service Master

This master acts as sort of a load balancer for data calls

Built in Event logger as well for quick graphical analysis

# Slave Communication

## Basics
Slave communication occurs with the following protocols:

#### Slave Heartbeat
  1. Master sends a Ping request to the Slave
  2. Slave responds to Ping with an Acknowledgement

#### Slave Data Request (small amount of data)
  1. Master sends a Data request packet
  2. Slave responds with a Data response packet with the `dataLeft == 0` OR an Error packet if
    unable to complete the request

#### Slave Data Request (large amount of data)
  1. Master sends a Data request packet
  2. Slave responds with a Data response packet with `dataLeft != 0` OR an Error packet if
    unable to complete the request
  3. Master sends a Data continue request.

  If `dataLeft == -1` then we must use the same slave to continue gathering data until we receive all the data.

  If `dataLeft > 1` and the number of slaves is > 1 then the master sends requests out to other slaves to get data faster. Since data may arrive in parts out of order then we must collate.

#### Slave Status Request
  1. Master sends a Status request packet
  2. Slave responds with a Status Response packet

## Security
  All communication occurs over an encrypted channel using AES256 in GCM mode (for use as an AEAD device) initialized using a shared secret obtained via ECDHE using a cryptographically proven curve.

  To complete the AEAD device a nonce is randomly initialized before initial discovery. The nonce is initialized **by the slave**, signed and sent along with the discovery request. This nonce is designed to be different for each slave.

  NOTE: Before I get a lot of flack from security experts. The reason I decided against using the standard TLS/SSL methods is to curtail the use of certificates and instead use a shared authenticity tag for use in the optional header. This allows the Master and all Slaves to generate a new private key each time they are initialized thus ensuring ease of maintenance and also at the same time ensure unique, secure, and authenticated channels of communication per each slave.
