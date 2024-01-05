package page

type Data struct {
	Title    string
	CSSPath  string
	HTMXPath string
}

type DataMutator interface {
	SetPageData(d *Data)
}

func mutateData(d Data, with any) Data {
	if m, ok := with.(DataMutator); ok {
		m.SetPageData(&d)
	}

	return d
}
