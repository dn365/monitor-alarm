package alarm

import (
	"errors"
	//"fmt"
    "strings"
    "strconv"
	//"os"

	bwlerrors "github.com/bobappleyard/bwl/errors"
	"github.com/bobappleyard/bwl/lexer"
)

const (
    _ = iota

    OR
    AND

    MAX
    MIN
    AVG
    COUNT
    SUM

    PRM_SEPARATOR
    INTEGER
    DECIMAL
    REL_OPER
    TXT

    PERIOD
    REPEAT
    THRESHOLD
    NS
    DMS
    PRM_PER
    DMS_PER
    FUNC_PER

    DMS_LIST
    FUNC_NAME
    FUNC_DEF
    METRIC_DEF
    FUNC_EXPRESSION
    EXPRESSION

    SPACE
)

type TokenSet map[int]string

type MatchedData struct {
    final int
    data string
}


var baseToken = TokenSet{
    PRM_SEPARATOR: ",",
    INTEGER: "[1-9]\\d*",
    DECIMAL: "\\-?[0-9]+(\\.[0-9]+)?",
    REL_OPER: "[<>=]|>=|<=|!=",
    TXT: "\\a[a-zA-Z0-9_.\\\\-]+",
    SPACE: "\\s",
}

var boolOperToken = TokenSet{
    OR: "[oO][rR]",
    AND: "[aA][nN][dD]",
}

var funcToken = TokenSet{ 
    MAX: "[mM][aA][xX]", 
    MIN: "[mM][iI][nN]", 
    COUNT: "[cC][oO][uU][nN][tT]", 
    AVG: "[aA][vV][gG]", 
    SUM: "[sS][uU][mM]",
}

var expressionBaseToken = TokenSet{
    PERIOD: "[1-9]\\d+",
    REPEAT: "times \\d+",
    THRESHOLD: baseToken[DECIMAL],
    NS: "\\a+(\\.[a-zA-Z0-9_.\\\\-]+)+",
    DMS: baseToken[TXT] + "=" + baseToken[TXT],
    PRM_PER: baseToken[PRM_SEPARATOR],
    DMS_PER: "[}{]",
    FUNC_PER: "[)(]",
}


func GetExpression(key int) string {
    dms_list := expressionBaseToken[DMS] + "(\\," + expressionBaseToken[DMS] + ")*" 
    metric_def := expressionBaseToken[NS] + "(\\{" + dms_list + "\\})?"

    funcNames := make([]string, 0, len(funcToken))
    for _,x := range funcToken {
        funcNames = append(funcNames,x)
    }
    func_param := metric_def + "(," + expressionBaseToken[PERIOD] + ")?"

    func_name := strings.Join(funcNames, "|")
    func_def := "(" + func_name + ")\\(" + func_param + "\\)"
    func_expression := func_def + "(" + baseToken[REL_OPER] + ")" + expressionBaseToken[THRESHOLD] + "(" + expressionBaseToken[REPEAT] + ")?"
    expression := metric_def + "(" + baseToken[REL_OPER] + ")" + expressionBaseToken[THRESHOLD]

    switch key {
        case OR, AND:
            return boolOperToken[key]
        case FUNC_NAME:
            return func_name
        case MAX, MIN, AVG, COUNT, SUM:
            return funcToken[key]
        case PRM_SEPARATOR, INTEGER, DECIMAL, REL_OPER, TXT, SPACE:
            return baseToken[key]
        case DMS_LIST: 
            return dms_list
        case METRIC_DEF: 
            return metric_def
        case FUNC_DEF:
            return func_def
        case FUNC_EXPRESSION:
            return func_expression
        case EXPRESSION:
            return expression
        case PERIOD, REPEAT, THRESHOLD, NS, DMS, PRM_PER, DMS_PER, FUNC_PER:
            return expressionBaseToken[key]
        default:
            return ""
    }
    return ""
}

func parse_string(str string, tokens []int) []MatchedData {
    l := lexer.New()
    res := make([]MatchedData, 0)

	for _, x := range tokens {
		l.ForceRegex(GetExpression(x), nil).SetFinal(x)
	}

    l.StartString(str)
    for !l.Eof() {
        f := l.Next()
        if f == -1 {
            bwlerrors.Fatal(errors.New("failed to match"))
        }
		//fmt.Printf("%d (%2d): %#v\n", f, l.Pos(), l.String())
        res = append(res, MatchedData{f, l.String()})
    }
    return res
}

func parse_func(str string) (string, *Metric, int) {
    tokens := []int{
        FUNC_NAME,
        DMS,
        PERIOD,
        NS,
        PRM_PER,
        DMS_PER,
        FUNC_PER,
        SPACE,
    }

    f_name, m_name, period := "", "", 0
    dms_list := map[string]string{}

	for _,x := range parse_string(str, tokens) {
        if x.final == NS {
            m_name = x.data
        } else if x.final == FUNC_NAME {
            f_name = x.data
        } else if x.final == PERIOD {
            period,_ = strconv.Atoi(x.data)
        } else if x.final == DMS {
            dp := strings.Split(x.data, "=")
            dms_list[dp[0]] = dp[1]
        } else {
            continue
        }
    }
	//fmt.Printf("%#v %#v %#v, %#v\n", f_name, m_name, dms_list, period)
    return f_name, NewMetric(m_name, dms_list), period
}

func ParseSubExpression(str string) *SubExpression {
    DEFAULT_PERIOD := 60
    DEFAULT_REPEAT := 1

    tokens := []int{
        FUNC_DEF,
        METRIC_DEF,
        THRESHOLD,
        REL_OPER,
        REPEAT,
        SPACE,
    }

    fun, rel, threshold, period, repeat := "", "", 0.0, 0, 0
    var metric *Metric

	for _,x := range parse_string(str, tokens) {
        if x.final == FUNC_DEF || x.final == METRIC_DEF {
            fun, metric, period = parse_func(x.data)
        } else if x.final == REL_OPER {
            rel = x.data
        } else if x.final == THRESHOLD {
            threshold,_ = strconv.ParseFloat(x.data, 64)
        } else if x.final == REPEAT {
            sl := strings.Split(x.data, " ")
            repeat,_ = strconv.Atoi(sl[1])
        }
    }

    if fun == "" {
        if rel == "<" || rel == "<=" {
            fun = "min"
        } else if rel == ">" || rel == ">=" {
            fun = "max"
        }
    }

    if period <= 0 {
        period = DEFAULT_PERIOD
    }

    if repeat <= 0 {
        repeat = DEFAULT_REPEAT
    }

    return NewSubExpression(fun, metric, rel, threshold, period, repeat)
}

func (self *Expression) Parse(str string) []interface{} {
    tokens := []int{
        OR,
        AND,
        EXPRESSION,
        FUNC_EXPRESSION,
        SPACE,
    }

    //str := "system.iostate.usage{hostname=node-01,devicename=xvda}<80 AND system.mem.usage>=85 or avg(system.iostat.rs,120)>500"
	for _,x := range parse_string(str, tokens) {
        if x.final == FUNC_EXPRESSION || x.final == EXPRESSION {
            self.elements = append(self.elements, ParseSubExpression(x.data))
        } else if x.final == OR || x.final == AND {
            self.elements = append(self.elements, strings.ToUpper(x.data))
        } else {
            continue
        }
	}

    return self.elements
}

