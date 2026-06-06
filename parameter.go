package fleurform

type ParameterMetadata struct {
	Name     string
	Required bool
}

func (p *ParameterMetadata) String() string {
	s := p.Name

	if !p.Required {
		s += " (optional)"
	}

	return s
}
