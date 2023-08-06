package hotpot

type Hotpot struct {
}

func NewHotpot() (*Hotpot, error) {
	info := Hotpot{}
	return &info, nil
}

func (h *Hotpot) BlackLists() ([]string, error) {
	return []string{}, nil
}

func (h *Hotpot) SaveRegistInfo(isRegisted bool) error {
	return nil
}
