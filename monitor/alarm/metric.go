package alarm

type Metric struct {
    name string
    dimensions map[string]string
}

func NewMetric(name string, dms map[string]string) *Metric {
    res := new(Metric)
    res.name = name
    res.dimensions = dms
    return res
}
