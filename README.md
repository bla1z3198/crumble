# Crumble v1t(testing accuracy of divide into parts)
## About
Conditionally, origin packet divides into crumbs, which have 
a random number of random sizes pieces of, previously, encrypted
data.
## How it works?
Payload are destroying into crumbs. Every crumb have a random part of payload,
and service info for collect this crumbs into a full cookie. Each crumb have flowID
and seq position.
## Structure
connector.go - main .go file
crusher - crush into crumbs
wrapper - wrap crumbs into byte slice
randomizer - random numbers for crumbs size
encryptor - package for encryption
collector - in a next commits
## Crumb struct
**Crumb** - is a main structure which contains payload and service info, as follows:
```
type Crumb struct {
	FlowID uint16
	Seq uint16
	Flags string
	Lost uint16
	Payload []byte
	Padding []byte
}
```
**FlowID** is an ID of current data flow.
**Seq** is a position of crumb in a current flow id.
**Flags** is an info about crumb. Fake or not fake.
**Lost** is a number of crumbs in a flow, to request lost crumbles (if this need).
**Payload** is a an data in array of bytes.
**Padding** is an trash in array of bytes.
## Crusher
Firstly, crusher define the bounds of encrypted parts (from, to). Then append crumb to array of
crumbs. Later crusher checking remain of encrypted, because not always he can divide data into
integer number of parts. And add "extra crumb" if it need, return array of crumbs ([]Crumb).
## Wrapper
This package just wraps crumbs into byte slice.
## Randomizer
This package generate random sizes for parts of data.
## Encryptor
This package encrypts data.
