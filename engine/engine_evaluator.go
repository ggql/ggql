package main

import (
    "errors"
    "regexp"
    "strings"
)

type Value interface{}

type Environment struct {
    Globals map[string]Value
}

type Expression interface {
    Kind() ExpressionKind
}

type ExpressionKind int

const (
    Assignment ExpressionKind = iota
    String
    Symbol
    GlobalVariable
    Number
    Boolean
    PrefixUnary
    Arithmetic
    Comparison
    Like
    Glob
    Logical
    Bitwise
    Call
    Between
    Case
    In
    IsNull
    Null
)

type AssignmentExpression struct {
    Symbol string
    Value  Expression
}

type StringExpression struct {
    ValueType StringValueType
    Value     string
}

type StringValueType int

const (
    Text StringValueType = iota
    Time
    Date
    DateTime
)

type SymbolExpression struct {
    Value string
}

type GlobalVariableExpression struct {
    Name string
}

type NumberExpression struct {
    Value float64
}

type BooleanExpression struct {
    IsTrue bool
}

type PrefixUnary struct {
    Operator PrefixUnaryOperator
    Right    Expression
}

type PrefixUnaryOperator int

const (
    Minus PrefixUnaryOperator = iota
    Bang
)

type ArithmeticExpression struct {
    Operator ArithmeticOperator
    Left     Expression
    Right    Expression
}

type ArithmeticOperator int

const (
    Plus ArithmeticOperator = iota
    Minus
    Star
    Slash
    Modulus
)

type ComparisonExpression struct {
    Operator ComparisonOperator
    Left     Expression
    Right    Expression
}

type ComparisonOperator int

const (
    Greater ComparisonOperator = iota
    GreaterEqual
    Less
    LessEqual
    Equal
    NotEqual
    NullSafeEqual
)

type LikeExpression struct {
    Input   Expression
    Pattern Expression
}

type GlobExpression struct {
    Input   Expression
    Pattern Expression
}

type LogicalExpression struct {
    Operator LogicalOperator
    Left     Expression
    Right    Expression
}

type LogicalOperator int

const (
    And LogicalOperator = iota
    Or
    Xor
)

type BitwiseExpression struct {
    Operator BitwiseOperator
    Left     Expression
    Right    Expression
}

type BitwiseOperator int

const (
    Or BitwiseOperator = iota
    And
    RightShift
    LeftShift
)

type CallExpression struct {
    FunctionName string
    Arguments    []Expression
}

type BetweenExpression struct {
    Value       Expression
    RangeStart  Expression
    RangeEnd    Expression
}

type CaseExpression struct {
    Conditions   []Expression
    Values       []Expression
    DefaultValue *Expression
}

type InExpression struct {
    Argument      Expression
    Values        []Expression
    HasNotKeyword bool
}

type IsNullExpression struct {
    Argument Expression
    HasNot   bool
}

func EvaluateExpression(env *Environment, expression Expression, titles []string, object []Value) (Value, error) {
    switch expression.Kind() {
    case Assignment:
        expr := expression.(*AssignmentExpression)
        return EvaluateAssignment(env, expr, titles, object)
    case String:
        expr := expression.(*StringExpression)
        return EvaluateString(expr)
    case Symbol:
        expr := expression.(*SymbolExpression)
        return EvaluateSymbol(expr, titles, object)
    case GlobalVariable:
        expr := expression.(*GlobalVariableExpression)
        return EvaluateGlobalVariable(env, expr)
    case Number:
        expr := expression.(*NumberExpression)
        return EvaluateNumber(expr), nil
    case Boolean:
        expr := expression.(*BooleanExpression)
        return EvaluateBoolean(expr), nil
    case PrefixUnary:
        expr := expression.(*PrefixUnary)
        return EvaluatePrefixUnary(env, expr, titles, object)
    case Arithmetic:
        expr := expression.(*ArithmeticExpression)
        return EvaluateArithmetic(env, expr, titles, object)
    case Comparison:
        expr := expression.(*ComparisonExpression)
        return EvaluateComparison(env, expr, titles, object)
    case Like:
        expr := expression.(*LikeExpression)
        return EvaluateLike(env, expr, titles, object)
    case Glob:
        expr := expression.(*GlobExpression)
        return EvaluateGlob(env, expr, titles, object)
    case Logical:
        expr := expression.(*LogicalExpression)
        return EvaluateLogical(env, expr, titles, object)
    case Bitwise:
        expr := expression.(*BitwiseExpression)
        return EvaluateBitwise(env, expr, titles, object)
    case Call:
        expr := expression.(*CallExpression)
        return EvaluateCall(env, expr, titles, object)
    case Between:
        expr := expression.(*BetweenExpression)
        return EvaluateBetween(env, expr, titles, object)
    case Case:
        expr := expression.(*CaseExpression)
        return EvaluateCase(env, expr, titles, object)
    case In:
        expr := expression.(*InExpression)
        return EvaluateIn(env, expr, titles, object)
    case IsNull:
        expr := expression.(*IsNullExpression)
        return EvaluateIsNull(env, expr, titles, object)
    case Null:
        return nil, nil
    default:
        return nil, errors.New("Invalid expression kind")
    }
}

func EvaluateAssignment(env *Environment, expr *AssignmentExpression, titles []string, object []Value) (Value, error) {
    value, err := EvaluateExpression(env, expr.Value, titles, object)
    if err != nil {
        return nil, err
    }
    env.Globals[expr.Symbol] = value
    return value, nil
}

func EvaluateString(expr *StringExpression) (Value, error) {
    switch expr.ValueType {
    case Text:
        return expr.Value, nil
    case Time:
        return expr.Value, nil
    case Date:
        return expr.Value, nil
    case DateTime:
        return expr.Value, nil
    default:
        return nil, errors.New("Invalid string value type")
    }
}

func EvaluateSymbol(expr *SymbolExpression, titles []string, object []Value) (Value, error) {
    for index, title := range titles {
        if expr.Value == title {
            return object[index], nil
        }
    }
    return nil, errors.New("Invalid column name")
}

func EvaluateGlobalVariable(env *Environment, expr *GlobalVariableExpression) (Value, error) {
    value, ok := env.Globals[expr.Name]
    if ok {
        return value, nil
    }
    return nil, errors.New("Global variable not found")
}

func EvaluateNumber(expr *NumberExpression) Value {
    return expr.Value
}

func EvaluateBoolean(expr *BooleanExpression) Value {
    return expr.IsTrue
}

func EvaluatePrefixUnary(env *Environment, expr *PrefixUnary, titles []string, object []Value) (Value, error) {
    rhs, err := EvaluateExpression(env, expr.Right, titles, object)
    if err != nil {
        return nil, err
    }
    switch expr.Operator {
    case Minus:
        if f, ok := rhs.(float64); ok {
            return -f, nil
        }
        return nil, errors.New("Invalid prefix unary operation")
    case Bang:
        if b, ok := rhs.(bool); ok {
            return !b, nil
        }
        return nil, errors.New("Invalid prefix unary operation")
    default:
        return nil, errors.New("Invalid prefix unary operator")
    }
}

func EvaluateArithmetic(env *Environment, expr *ArithmeticExpression, titles []string, object []Value) (Value, error) {
    lhs, err := EvaluateExpression(env, expr.Left, titles, object)
    if err != nil {
        return nil, err
    }
    rhs, err := EvaluateExpression(env, expr.Right, titles, object)
    if err != nil {
        return nil, err
    }
    switch expr.Operator {
    case Plus:
        return lhs.(float64) + rhs.(float64), nil
    case Minus:
        return lhs.(float64) - rhs.(float64), nil
    case Star:
        return lhs.(float64) * rhs.(float64), nil
    case Slash:
        return lhs.(float64) / rhs.(float64), nil
    case Modulus:
        return int(lhs.(float64)) % int(rhs.(float64)), nil
    default:
        return nil, errors.New("Invalid arithmetic operator")
    }
}

func EvaluateComparison(env *Environment, expr *ComparisonExpression, titles []string, object []Value) (Value, error) {
    lhs, err := EvaluateExpression(env, expr.Left, titles, object)
    if err != nil {
        return nil, err
    }
    rhs, err := EvaluateExpression(env, expr.Right, titles, object)
    if err != nil {
        return nil, err
    }
    leftType := getType(lhs)
    comparisonResult := 0
    if leftType == "int" {
        if lhs.(int) > rhs.(int) {
            comparisonResult = 1
        } else if lhs.(int) < rhs.(int) {
            comparisonResult = -1
        }
    } else if leftType == "float64" {
        if lhs.(float64) > rhs.(float64) {
            comparisonResult = 1
        } else if lhs.(float64) < rhs.(float64) {
            comparisonResult = -1
        }
    } else if leftType == "bool" {
        if lhs.(bool) == rhs.(bool) {
            comparisonResult = 0
        } else {
            comparisonResult = 1
        }
    } else {
        if lhs.(string) > rhs.(string) {
            comparisonResult = 1
        } else if lhs.(string) < rhs.(string) {
            comparisonResult = -1
        }
    }
    if expr.Operator == NullSafeEqual {
        if isNull(lhs) && isNull(rhs) {
            return 1, nil
        } else if isNull(lhs) || isNull(rhs) {
            return 0, nil
        } else if comparisonResult == 0 {
            return 1, nil
        } else {
            return 0, nil
        }
    }
    switch expr.Operator {
    case Greater:
        return comparisonResult > 0, nil
    case GreaterEqual:
        return comparisonResult >= 0, nil
    case Less:
        return comparisonResult < 0, nil
    case LessEqual:
        return comparisonResult <= 0, nil
    case Equal:
        return comparisonResult == 0, nil
    case NotEqual:
        return comparisonResult != 0, nil
    default:
        return nil, errors.New("Invalid comparison operator")
    }
}

func EvaluateLike(env *Environment, expr *LikeExpression, titles []string, object []Value) (Value, error) {
    rhs, err := EvaluateExpression(env, expr.Pattern, titles, object)
    if err != nil {
        return nil, err
    }
    pattern := "^" + strings.ToLower(rhs.(string))
    pattern = strings.ReplaceAll(pattern, "%", ".*")
    pattern = strings.ReplaceAll(pattern, "_", ".")
    regex, err := regexp.Compile(pattern)
    if err != nil {
        return nil, err
    }
    lhs, err := EvaluateExpression(env, expr.Input, titles, object)
    if err != nil {
        return nil, err
    }
    return regex.MatchString(strings.ToLower(lhs.(string))), nil
}

func EvaluateGlob(env *Environment, expr *GlobExpression, titles []string, object []Value) (Value, error) {
    rhs, err := EvaluateExpression(env, expr.Pattern, titles, object)
    if err != nil {
        return nil, err
    }
    pattern := "^" + regexp.QuoteMeta(rhs.(string))
    pattern = strings.ReplaceAll(pattern, "*", ".*")
    pattern = strings.ReplaceAll(pattern, "?", ".")
    regex, err := regexp.Compile(pattern)
    if err != nil {
        return nil, err
    }
    lhs, err := EvaluateExpression(env, expr.Input, titles, object)
    if err != nil {
        return nil, err
    }
    return regex.MatchString(lhs.(string)), nil
}

func EvaluateLogical(env *Environment, expr *LogicalExpression, titles []string, object []Value) (Value, error) {
    lhs, err := EvaluateExpression(env, expr.Left, titles, object)
    if err != nil {
        return nil, err
    }
    if expr.Operator == And && !lhs.(bool) {
        return false, nil
    }
    if expr.Operator == Or && lhs.(bool) {
        return true, nil
    }
    rhs, err := EvaluateExpression(env, expr.Right, titles, object)
    if err != nil {
        return nil, err
    }
    switch expr.Operator {
    case And:
        return lhs.(bool) && rhs.(bool), nil
    case Or:
        return lhs.(bool) || rhs.(bool), nil
    case Xor:
        return lhs.(bool) != rhs.(bool), nil
    default:
        return nil, errors.New("Invalid logical operator")
    }
}

func EvaluateBitwise(env *Environment, expr *BitwiseExpression, titles []string, object []Value) (Value, error) {
    lhs, err := EvaluateExpression(env, expr.Left, titles, object)
    if err != nil {
        return nil, err
    }
    rhs, err := EvaluateExpression(env, expr.Right, titles, object)
    if err != nil {
        return nil, err
    }
    switch expr.Operator {
    case Or:
        return lhs.(int) | rhs.(int), nil
    case And:
        return lhs.(int) & rhs.(int), nil
    case RightShift:
        if rhs.(int) >= 64 {
            return nil, errors.New("Attempt to shift right with overflow")
        }
        return lhs.(int) >> rhs.(int), nil
    case LeftShift:
        if rhs.(int) >= 64 {
            return nil, errors.New("Attempt to shift left with overflow")
        }
        return lhs.(int) << rhs.(int), nil
    default:
        return nil, errors.New("Invalid bitwise operator")
    }
}

func EvaluateCall(env *Environment, expr *CallExpression, titles []string, object []Value) (Value, error) {
    function := FUNCTIONS[expr.FunctionName]
    if function == nil {
        return nil, errors.New("Function not found")
    }
    arguments := make([]Value, len(expr.Arguments))
    for i, arg := range expr.Arguments {
        value, err := EvaluateExpression(env, arg, titles, object)
        if err != nil {
            return nil, err
        }
        arguments[i] = value
    }
    return function(arguments), nil
}

func EvaluateBetween(env *Environment, expr *BetweenExpression, titles []string, object []Value) (Value, error) {
    value, err := EvaluateExpression(env, expr.Value, titles, object)
    if err != nil {
        return nil, err
    }
    rangeStart, err := EvaluateExpression(env, expr.RangeStart, titles, object)
    if err != nil {
        return nil, err
    }
    rangeEnd, err := EvaluateExpression(env, expr.RangeEnd, titles, object)
    if err != nil {
        return nil, err
    }
    return value.(int) >= rangeStart.(int) && value.(int) <= rangeEnd.(int), nil
}

func EvaluateCase(env *Environment, expr *CaseExpression, titles []string, object []Value) (Value, error) {
    conditions := expr.Conditions
    values := expr.Values
    for i := 0; i < len(conditions); i++ {
        condition, err := EvaluateExpression(env, conditions[i], titles, object)
        if err != nil {
            return nil, err
        }
        if condition.(bool) {
            return EvaluateExpression(env, values[i], titles, object)
        }
    }
    if expr.DefaultValue != nil {
        return EvaluateExpression(env, *expr.DefaultValue, titles, object)
    }
    return nil, errors.New("Invalid case statement")
}

func EvaluateIn(env *Environment, expr *InExpression, titles []string, object []Value) (Value, error) {
    argument, err := EvaluateExpression(env, expr.Argument, titles, object)
    if err != nil {
        return nil, err
    }
    for _, valueExpr := range expr.Values {
        value, err := EvaluateExpression(env, valueExpr, titles, object)
        if err != nil {
            return nil, err
        }
        if argument == value {
            return !expr.HasNotKeyword, nil
        }
    }
    return expr.HasNotKeyword, nil
}

func EvaluateIsNull(env *Environment, expr *IsNullExpression, titles []string, object []Value) (Value, error) {
    argument, err := EvaluateExpression(env, expr.Argument, titles, object)
    if err != nil {
        return nil, err
    }
    isNull := isNull(argument)
    if expr.HasNot {
        return !isNull, nil
    }
    return isNull, nil
}

func getType(v Value) string {
    switch v.(type) {
    case int:
        return "int"
    case float64:
        return "float64"
    case bool:
        return "bool"
    case string:
        return "string"
    default:
        return ""
    }
}

func isNull(v Value) bool {
    return v == nil
}

var FUNCTIONS = map[string]func([]Value) Value{}

func main() {
    // Your code here
}
