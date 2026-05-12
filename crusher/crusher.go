package crusher

var (
	crumbs       []Crumb
	padding      []byte
	completed    []byte
	payload_size uint16
	from         uint16
	to           uint16
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
		if i == 0 {
			from = 0
			to = s.One + 1
		} else {
			from = i*s.One + 1
			to = (i+1)*s.One + 1
		}
		crumbs = append(crumbs,
			Crumb{
				FlowID:  s.ID,
				Seq:     i,
				Flags:   s.Flg,
				Lost:    s.Parts + 1,
				Payload: s.Encrypted[from:to],
				Padding: padding,
			})
	}
	if len(s.Encrypted) > int(to) {
		crumbs = append(crumbs,
			Crumb{
				FlowID:  s.ID,
				Seq:     s.Parts,
				Flags:   "LAST",
				Lost:    s.Parts + 1,
				Payload: s.Encrypted[(s.Parts*s.One)+1:],
				Padding: padding,
			})
	}
	return crumbs
}
