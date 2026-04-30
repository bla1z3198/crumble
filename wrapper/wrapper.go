package wrapper

import (
	"crusher/crusher"
)

var (
	bytes []byte
)

func Wrap(crumb crusher.Crumb) []byte {
	size := 25 + 12 +
		len(crumb.Payload) +
		len(crumb.Padding)
	bytes = make([]byte, size)

	copy(bytes[0:24], "GET /ya.ru HTTP/1.1")

	bytes[25] = byte(crumb.FlowID >> 8)
	bytes[26] = byte(crumb.FlowID)

	bytes[27] = byte(crumb.Seq >> 8)
	bytes[28] = byte(crumb.Seq)

	copy(bytes[29:33], []byte(crumb.Flags))

	bytes[34] = byte(crumb.Lost >> 8)
	bytes[35] = byte(crumb.Lost)

	copy(bytes[36:36+len(crumb.Payload)],
		crumb.Payload)
	copy(bytes[37+len(crumb.Payload):],
		crumb.Padding)

	return bytes
}
