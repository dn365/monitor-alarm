package main

import (
	"fmt"
	"github.com/WH-Wang/monitor/alarm"
	"strconv"
)

func main() {
	str := "system.iostate.usage{hostname=node-01,devicename=xvda}<80.2 AND system.mem.usage>=85 or avg(system.iostat.rs,120)>500"
	e := alarm.ExpressionOf(str)
	subAlarms := map[string]*alarm.SubExpression{}

	ses := e.GetSubExpressions()
	for i, x := range ses {
		id := "id#" + strconv.Itoa(i)
		subAlarms[id] = x
	}
	event := alarm.AlarmDefinitionCreatedEvent("abcdefg", "123456", "test alarm", "this is a test alarm", e.GetExpression(), subAlarms, []string{})
	fmt.Println(event)

	fmt.Println("")

	metrics := []*alarm.Metric{
		alarm.NewMetric("metric-01", map[string]string{"key1": "val1", "key2": "val2"}),
		alarm.NewMetric("metric-02", map[string]string{"key1": "val1", "key2": "val2"}),
	}

	alarmMetrics := map[string][]*alarm.Metric{
		"alarm-id-01": metrics,
	}

	subExpressions := map[string]*alarm.SubExpression{
		"sub-expression-01": alarm.SubExpressionOf("avg(system.mem.buffers{hostname=node01}) > 800"),
		"sub-expression-02": alarm.SubExpressionOf("min(system.load.1{hostname=node01}) < 5.0"),
	}

	for k, x := range alarmMetrics {
		fmt.Println(alarm.AlarmDeletedEvent("tenant-id-01", k, x, "alarm-def-id-01", subExpressions))
	}

	fmt.Println("")

	metricMaps := map[string]*alarm.Metric{
		"metric-01": alarm.NewMetric("metric-01", map[string]string{"key1": "val1", "key2": "val2"}),
		"metric-02": alarm.NewMetric("metric-02", map[string]string{"key1": "val1", "key2": "val2"}),
	}

	fmt.Println(alarm.AlarmDefinitionDeletedEvent("alarm-def-id-01", metricMaps))
}
