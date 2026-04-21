package crusher

type Crumb struct {
	FlowID     int
	Seq        int
	Flags      string
	PayloadLen int
	Payload    []byte
	Padding    []byte
}

func Encryption(data []byte) []byte {
	return data
}

func Crusher(encrypted []byte, flowID int) []Crumb {
	crumbs := make([]Crumb, len(encrypted))
	pad := 10

	for i := range encrypted {
		crumbs[i] = Crumb{
			FlowID:     flowID,
			Seq:        i,
			Flags:      "DATA",
			PayloadLen: pad,
			Payload:    encrypted[i : i+1],
			Padding:    encrypted[i : i+pad],
		}
	}
	return crumbs
}

func Sender() {

}
