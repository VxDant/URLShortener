package model

type Dataset struct {
	ValueSet string
}

func DatasetNew() *Dataset {
	return &Dataset{
		ValueSet: "abcdABCD",
	}
}
