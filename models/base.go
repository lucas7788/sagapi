package models

type APIExtra struct {
	Popularity      int
	Delay           int
	SuccessRate     int
	InvokeFrequency int
}

type APIDetailInstruction struct {
	DataDesc            string
	DataSource          string
	ApplicationScenario string
}
