package parser

type Diagnostic struct {
	Label    string
	Message  string
	Location Location
	Notes    []string
	Helps    []string
	Docs     string
}

func NewDiagnostic(label, message string) *Diagnostic {
	return &Diagnostic{
		Label:    label,
		Message:  message,
		Location: Location{},
		Notes:    []string{},
		Helps:    []string{},
		Docs:     "",
	}
}

func NewError(message string) *Diagnostic {
	return NewDiagnostic("Error", message)
}

func NewException(message string) *Diagnostic {
	return NewDiagnostic("Exception", message)
}

func (d *Diagnostic) WithLocation(location Location) *Diagnostic {
	d.Location = location
	return d
}

func (d *Diagnostic) WithLocationSpan(start, end int) *Diagnostic {
	d.Location = Location{start, end}
	return d
}

func (d *Diagnostic) AddNote(note string) *Diagnostic {
	d.Notes = append(d.Notes, note)
	return d
}

func (d *Diagnostic) AddHelp(help string) *Diagnostic {
	d.Helps = append(d.Helps, help)
	return d
}

func (d *Diagnostic) WithDocs(docs string) *Diagnostic {
	d.Docs = docs
	return d
}

func (d *Diagnostic) DiaLabel() string {
	return d.Label
}

func (d *Diagnostic) DiaMessage() string {
	return d.Message
}

func (d *Diagnostic) DiaLocation() Location {
	return d.Location
}

func (d *Diagnostic) DiaNotes() []string {
	return d.Notes
}

func (d *Diagnostic) DiaHelps() []string {
	return d.Helps
}

func (d *Diagnostic) DiaDocs() string {
	return d.Docs
}
