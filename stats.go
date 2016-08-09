package invt

type Stats struct {
	TotalTrades   int
	ProfitTrades  int
	LossTrades    int
	AverageProfit float64
	AverageLoss   float64
	MaxProfit     float64
	MaxLoss       float64
}

func NewStats() *Stats {
	return &Stats{}
}

func (s *Stats) AddTrade(pl float64) {
	s.TotalTrades += 1
	if pl > 0 {
		s.AverageProfit = (float64(s.ProfitTrades)*s.AverageProfit + pl) / (float64(s.ProfitTrades) + 1.)
		s.ProfitTrades += 1
		if s.MaxProfit < pl {
			s.MaxProfit = pl
		}
	} else {
		s.AverageLoss = (float64(s.LossTrades)*s.AverageLoss + pl) / (float64(s.LossTrades) + 1.)
		s.LossTrades += 1
		if s.MaxLoss > pl {
			s.MaxLoss = pl
		}
	}
}
