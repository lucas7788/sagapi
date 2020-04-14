package models

type ApiExtra struct {
	Popularity      int
	Delay           int
	SuccessRate     int
	InvokeFrequency int
}

type ApiDetailInstruction struct {
	DataDesc            string
	DataSource          string
	ApplicationScenario string
}
