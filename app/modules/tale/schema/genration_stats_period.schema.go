package schema

type GenerationPeriodStatsSchema struct {
	DateFrom string `query:"date_from"`
	DateTo   string `query:"date_to"`
}
