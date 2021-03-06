# Data Service Master

This master acts as sort of a load balancer for data calls

Built in Event logger as well for quick graphical analysis

# Configuration

Configuration files are expected to be in YAML format. This program checks for configuration files in the same directory
it is being run from and checks the following paths in order:

 - ./config.yml
 - ./config.yaml

Full config specs to come!

# Slave Communication

## Basics
Slave communication occurs with the following protocols:

#### Slave Heartbeat
  1. Master sends a Ping request to the Slave
  2. Slave responds to Ping with an Acknowledgement

#### Slave Discovery
  1. Slave sends a Discovery request to Master along with its unique ID
  2. Master sends its public key and waits for the slave's public key so both can create a shared secret.
  3. After secure symmetric encryption is obtained a cryptographically random nonce is generated and a Challenge packet
  is sent.
  4. The Slave then combines the nonce (as salt) to the pre-determined secret agreed upon out of band and hashes this,
  sending the hash back as a Challenge Response.
  5. If the hash matches what is expected, then an Accept packet is sent, otherwise a Rejected packet is sent.

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
  All communication occurs over an encrypted channel using AES256 in GCM mode (for use as an AEAD device) initialized using a shared secret obtained via ECDHE using a cryptographically proven curve. After symmetric encryption is implemented, a
  nonce is generated by the Master. This is then sent in the form of a Challenge Packet. In order to successfully complete
  the challenge and be accepted by the Master as a trusted data source, the slave must send back this nonce and a
  pre-determined secret (determined by the admin). If either is out of sync, then the Challenge is failed and the slave
  is rejected as a new member of the network.

  NOTE: Before I get a lot of flack from security experts. The reason I decided against using the standard TLS/SSL methods is to curtail the use of certificates and instead use a shared pre-determined secret. This allows the Master and all Slaves to generate a new private key each time they are initialized thus ensuring ease of maintenance and also at the same time ensure unique, secure, and authenticated channels of communication per each slave.

  AND YES THERE IS CURRENTLY SUSCEPTIBILITY TO A MAN-IN-THE-MIDDLE ATTACK! But don't fret, I'm fixin' this soon. Just
  let me get through the basic protocol development.
