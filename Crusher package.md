Crusher package contains `crusher.go` file, which crush 
packets a random number of pieces of random sizes.
## About
Conditionally, origin packet divides into crumbs, which have 
a random number of random sizes pieces of, previously, encrypted
data.
## How it works?
Payload are destroying into crumbs. Every crumb have a random part of payload,
and service info for collect this crumbs into a full cookie. Each crumb have flowID
and seq position.
![[raw to encrypted.png]]
## Crumb struct
**Crumb** - is a main structure which contains payload and service info, as follows:

```
type Crumb struct {
	FlowID int
	Seq int
	Flags string
	PayloadLen int
	Payload []byte
	Padding []byte
}
```

**FlowID** is an ID of current data flow.  
**Seq** is a position of crumb in a current flow id.
**Flags** is an info about crumb. Fake or not fake.
**Payloadlen** is a size of payload field.
**Payload** is a an data in array of bytes.
**Padding** is an trash in array of bytes.

For example crusher can create a random number of flows (goroutines), and each crumb
have an unique **FlowID** and **Seq**. So, for ident correct sequence of **crumbs**, **collector** is need
this service info (FlowID and Seq). Then


![[encrypted to crumbs.png]]
## Random system
Crusher uses random numbers in a **Flags** and **PayloadLen** define. With **90%** chance next flag
will be "**DATA**", and with chance **10%** next flag will be "**FAKE**". **Payloadlen** now are full random
integer number (it will have an adaptive random in next versions). **Padding** is a []byte array, which will be added to the end. It's length is also full random. This array always contains trash info.
## Double-encryption
After creating all **crumbs**, every of them will be encrypted in **Ecryption** function, and then send to destination **IP**. 
## Send
All previously encrypted **crumbs** will be sent over **TCP**.
