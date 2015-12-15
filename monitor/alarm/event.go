package alarm

import (
    "strings"
    "strconv"
)

type JsonProperty struct {
    key string
    value interface{}
}

type JsonHandler interface {
    ToJson() string
}

func (self *JsonProperty) ToJson() string {
    str := ""
    if self.key != "" {
        str += "\""
        str += self.key
        str += "\":"
    }
    switch self.value.(type) {
        case int:
            str += strconv.Itoa(self.value.(int))
        case float64:
            str += strconv.FormatFloat(self.value.(float64), 'f', -1, 64)
        case string:
            str += "\""
            str += self.value.(string)
            str += "\""
        case []string:
            str += "[\""
            str += strings.Join(self.value.([]string), "\",\"")
            str += "\"]"
        case map[string]string:
            sl := []string{}
            str += "{"
            for k,v := range self.value.(map[string]string) {
                p := JsonProperty{k, v}
                sl = append(sl, p.ToJson())
            }
            str += strings.Join(sl, ",")
            str += "}"
        case []int:
            str += "["
            sl := []string{}
            for _,x := range self.value.([]int) {
                sl = append(sl, string(x))
            }
            str += strings.Join(sl, ",")
            str += "]"
        case *SubExpression:
            e := self.value.(*SubExpression)
            str += e.ToJson()
        case *Metric:
            m := self.value.(*Metric)
            str += m.ToJson()
        case map[string]*Metric:
            sl := []string{}
            str += "{"
            for k,v := range self.value.(map[string]*Metric) {
                p := JsonProperty{k, v}
                sl = append(sl, p.ToJson())
            }
            str += strings.Join(sl, ",")
            str += "}"
        case map[string]*SubExpression:
            sl := []string{}
            str += "{"
            for k,v := range self.value.(map[string]*SubExpression) {
                p := JsonProperty{k, v}
                sl = append(sl, p.ToJson())
            }
            str += strings.Join(sl, ",")
            str += "}"
        case []*Metric:
            sl := []string{}
            str += "["
            for _,v := range self.value.([]*Metric) {
                p := JsonProperty{"", v}
                sl = append(sl, p.ToJson())
            }
            str += strings.Join(sl, ",")
            str += "]"
        case []JsonProperty:
            sl := []string{}
            str += "{"
            for _,v := range self.value.([]JsonProperty) {
                sl = append(sl, v.ToJson())
            }
            str += strings.Join(sl, ",")
            str += "}"
    }
    return str
}

func (self *Metric) ToJson() string {
    props := []JsonProperty{}
    props = append(props, JsonProperty{"name", self.name})
    if len(self.dimensions) > 0 {
        props = append(props, JsonProperty{"dimensions", self.dimensions})
    }
    
    ej := JsonProperty{"", props}
    return ej.ToJson()
}

func (self *SubExpression) ToJson() string {
    props := []JsonProperty{}
    props = append(props, JsonProperty{"function", self.function})
    props = append(props, JsonProperty{"metricDefinition", self.metric})
    props = append(props, JsonProperty{"period", self.period})
    props = append(props, JsonProperty{"operator", self.operator})
    props = append(props, JsonProperty{"threshold", self.threshold})
    props = append(props, JsonProperty{"periods", self.periods})
    
    ej := JsonProperty{"", props}
    return ej.ToJson()
}

func AlarmDefinitionCreatedEvent(tenantId, alarmDefId, alertName, description, expression string, 
        subAlarms map[string]*SubExpression, matchBy []string) string {
    res := ""

    name := "alarm-definition-created"
    props := []JsonProperty{}
    props = append(props, JsonProperty{"alarmDefinitionId", alarmDefId})
    props = append(props, JsonProperty{"tenantId", tenantId})
    props = append(props, JsonProperty{"alarmName", alertName})
    props = append(props, JsonProperty{"description", description})
    props = append(props, JsonProperty{"alarmExpression", expression})
    props = append(props, JsonProperty{"alarmSubExpressions", subAlarms})
    if len(matchBy) > 0 {
        props = append(props, JsonProperty{"matchBy", matchBy})
    }
    
    ej := JsonProperty{name, props}

    res += "{"
    res += ej.ToJson()
    res += "}"
    return res
}

func AlarmDeletedEvent(tenantId, alarmId string, alarmMetrics []*Metric, 
        alarmDefinitionId string, subAlarms map[string]*SubExpression) string {
    res := ""

    name := "alarm-deleted"
    props := []JsonProperty{}
    props = append(props, JsonProperty{"tenantId", tenantId})
    props = append(props, JsonProperty{"alarmId", alarmId})
    props = append(props, JsonProperty{"alarmDefinitionId", alarmDefinitionId})
    props = append(props, JsonProperty{"alarmMetrics", alarmMetrics})
    props = append(props, JsonProperty{"subAlarms", subAlarms})
    
    ej := JsonProperty{name, props}

    res += "{"
    res += ej.ToJson()
    res += "}"
    return res
}


func AlarmDefinitionDeletedEvent(alarmDefId string, subAlarmMetricDefinitions map[string]*Metric) string {
    res := ""

    name := "alarm-definition-deleted"
    props := []JsonProperty{}
    props = append(props, JsonProperty{"alarmDefinitionId", alarmDefId})
    props = append(props, JsonProperty{"subAlarmMetricDefinitions", subAlarmMetricDefinitions})
    
    ej := JsonProperty{name, props}

    res += "{"
    res += ej.ToJson()
    res += "}"
    return res
}

