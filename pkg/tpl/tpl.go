package tpl

type Reporter struct {
	Section []ReportSection
}

type ReportSection struct {
	Name string
	Line []ReportLine
}

type ReportLine struct {
	Date     string
	Weekday  string
	Times    int
	Earliest string
	Latest   string
	Status   string
}
