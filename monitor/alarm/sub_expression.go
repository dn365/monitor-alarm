package alarm

import "strings"

type SubExpression struct {
    function, operator string
    period, periods int
    threshold float64
    metric *Metric
}

func NewSubExpression(
        function string, 
        metricDefinition *Metric, 
        operator string, 
        threshold float64, 
        period int, 
        periods int) *SubExpression {
    res := new(SubExpression)
    res.function = strings.ToUpper(function)
    res.metric = metricDefinition
    switch operator {
    case "<":
        res.operator = "lt"
    case ">":
        res.operator = "gt"
    case "<=":
        res.operator = "lte"
    case ">=":
        res.operator = "gte"
    default:
        res.operator = operator
    }
    res.threshold = threshold
    res.period = period
    res.periods = periods
    return res
}

func SubExpressionOf(expr string) *SubExpression {
    res := ParseSubExpression(expr)
    return res
}

