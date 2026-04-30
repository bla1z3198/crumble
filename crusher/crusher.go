package crusher

var (
	crumbs       []Crumb
	padding      []byte
	completed    []byte
	payload_size uint16
)

type Crumb struct {
	FlowID  uint16
	Seq     uint16
	Flags   string
	Lost    uint16
	Payload []byte
	Padding []byte
}

type Service struct {
	Encrypted []byte
	ID        uint16
	Flg       string
	Parts     uint16
	One       uint16
}

func Crush(s *Service) []Crumb {
	crumbs = make([]Crumb, 0)
	padding = make([]byte, uint16(len(s.Encrypted))/s.Parts)

	for i := range s.Parts {
		crumbs = append(crumbs,
			Crumb{
				FlowID:  s.ID,
				Seq:     i,
				Flags:   s.Flg,
				Lost:    0,
				Payload: s.Encrypted[i+payload_size : i+payload_size+(s.One)],
				Padding: padding[0 : len(padding)-int(s.Parts)],
			})
		payload_size = s.One
	}
	return crumbs
}
