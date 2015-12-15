package alarm

type Expression struct {
    expression string
    elements []interface{}
    subExpressions []*SubExpression
}

func NewExpression() *Expression {
    res := new(Expression)
    return res
}

func ExpressionOf(expr string) *Expression {
    res := new(Expression)
    res.expression = expr
    res.elements = res.Parse(res.expression)
    res.subExpressions = make([]*SubExpression, 0)
    return res
}

func (self *Expression) GetExpression() string {
    return self.expression
}

func (self *Expression) GetSubExpressions() []*SubExpression {
    if len(self.subExpressions) <= 0 {
        for _,x := range self.elements {
            switch x.(type) {
                case *SubExpression:
                    self.subExpressions = append(self.subExpressions, x.(*SubExpression))
                default:
                    break
            }
        }
    }

    return self.subExpressions
}
