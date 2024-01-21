package parser

type Diagnostic struct {
	label    string
	message  string
	location Location
	notes    []string
	helps    []string
	docs     string
}

func NewDiagnostic(label, message string) *Diagnostic {
	return &Diagnostic{
		label:    label,
		message:  message,
		location: Location{},
		notes:    []string{},
		helps:    []string{},
		docs:     "",
	}
}

func NewError(message string) *Diagnostic {
	return NewDiagnostic("Error", message)
}

func NewException(message string) *Diagnostic {
	return NewDiagnostic("Exception", message)
}

func (d *Diagnostic) WithLocation(location Location) *Diagnostic {
	d.location = location
	return d
}

func (d *Diagnostic) WithLocationSpan(start, end int) *Diagnostic {
	d.location = Location{start, end}
	return d
}

func (d *Diagnostic) AddNote(note string) *Diagnostic {
	d.notes = append(d.notes, note)
	return d
}

func (d *Diagnostic) AddHelp(help string) *Diagnostic {
	d.helps = append(d.helps, help)
	return d
}

func (d *Diagnostic) WithDocs(docs string) *Diagnostic {
	d.docs = docs
	return d
}

func (d *Diagnostic) Label() string {
	return d.label
}

func (d *Diagnostic) Message() string {
	return d.message
}

func (d *Diagnostic) Location() Location {
	return d.location
}

func (d *Diagnostic) Notes() []string {
	return d.notes
}

func (d *Diagnostic) Helps() []string {
	return d.helps
}

func (d *Diagnostic) Docs() string {
	return d.docs
}
