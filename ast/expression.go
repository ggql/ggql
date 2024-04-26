package ast

type ExpressionKind int

const (
	ExprAssignment ExpressionKind = iota
	ExprString
	ExprSymbol
	ExprGlobalVariable
	ExprNumber
	ExprBoolean
	ExprPrefixUnary
	ExprArithmetic
	ExprComparison
	ExprLike
	ExprGlob
	ExprLogical
	ExprBitwise
	ExprCall
	ExprBetween
	ExprCase
	ExprIn
	ExprIsNull
	ExprNull
)

type Expression interface {
	Kind() ExpressionKind
	ExprType(scope *Environment) DataType
	AsAny() interface{}
	IsConst() bool
}

type AssignmentExpression struct {
	Symbol string
	Value  Expression
}

func (e *AssignmentExpression) Kind() ExpressionKind {
	return ExprAssignment
}

func (e *AssignmentExpression) ExprType(scope *Environment) DataType {
	return e.Value.ExprType(scope)
}

func (e *AssignmentExpression) AsAny() interface{} {
	return e
}

func (e *AssignmentExpression) IsConst() bool {
	return false
}

type StringValueType int

const (
	StringValueText StringValueType = iota
	StringValueTime
	StringValueDate
	StringValueDateTime
)

type StringExpression struct {
	Value     string
	ValueType StringValueType
}

func (e *StringExpression) Kind() ExpressionKind {
	return ExprString
}

func (e *StringExpression) ExprType(scope *Environment) DataType {
	switch e.ValueType {
	case StringValueText:
		return Text{}
	case StringValueTime:
		return Time{}
	case StringValueDate:
		return Date{}
	case StringValueDateTime:
		return DateTime{}
	default:
		return Undefined{}
	}
}

func (e *StringExpression) AsAny() interface{} {
	return e
}

func (e *StringExpression) IsConst() bool {
	return true
}

type SymbolExpression struct {
	Value string
}

func (e *SymbolExpression) Kind() ExpressionKind {
	return ExprSymbol
}

func (e *SymbolExpression) ExprType(scope *Environment) DataType {
	if scope.Contains(e.Value) {
		return scope.Scopes[e.Value]
	}

	if buf, ok := TablesFieldsTypes[e.Value]; ok {
		return buf
	}

	return Undefined{}
}

func (e *SymbolExpression) AsAny() interface{} {
	return e
}

func (e *SymbolExpression) IsConst() bool {
	return false
}

type GlobalVariableExpression struct {
	Name string
}

func (e *GlobalVariableExpression) Kind() ExpressionKind {
	return ExprGlobalVariable
}

func (e *GlobalVariableExpression) ExprType(scope *Environment) DataType {
	if scope.Contains(e.Name) {
		return scope.Scopes[e.Name]
	}

	if buf, ok := TablesFieldsTypes[e.Name]; ok {
		return buf
	}

	return Undefined{}
}

func (e *GlobalVariableExpression) AsAny() interface{} {
	return e
}

func (e *GlobalVariableExpression) IsConst() bool {
	return false
}

type NumberExpression struct {
	Value Value
}

func (e *NumberExpression) Kind() ExpressionKind {
	return ExprNumber
}

func (e *NumberExpression) ExprType(scope *Environment) DataType {
	return e.Value.DataType()
}

func (e *NumberExpression) AsAny() interface{} {
	return e
}

func (e *NumberExpression) IsConst() bool {
	return true
}

type BooleanExpression struct {
	IsTrue bool
}

func (e *BooleanExpression) Kind() ExpressionKind {
	return ExprBoolean
}

func (e *BooleanExpression) ExprType(scope *Environment) DataType {
	return Boolean{}
}

func (e *BooleanExpression) AsAny() interface{} {
	return e
}

func (e *BooleanExpression) IsConst() bool {
	return true
}

type PrefixUnaryOperator int

const (
	POMinus PrefixUnaryOperator = iota
	POBang
)

type PrefixUnary struct {
	Right Expression
	Op    PrefixUnaryOperator
}

func (e *PrefixUnary) Kind() ExpressionKind {
	return ExprPrefixUnary
}

func (e *PrefixUnary) ExprType(scope *Environment) DataType {
	if e.Op == POBang {
		return Boolean{}
	} else {
		return Integer{}
	}
}

func (e *PrefixUnary) AsAny() interface{} {
	return e
}

func (e *PrefixUnary) IsConst() bool {
	return false
}

type ArithmeticOperator int

const (
	AOPlus ArithmeticOperator = iota
	AOMinus
	AOStar
	AOSlash
	AOModulus
)

type ArithmeticExpression struct {
	Left     Expression
	Operator ArithmeticOperator
	Right    Expression
}

func (e *ArithmeticExpression) Kind() ExpressionKind {
	return ExprArithmetic
}

func (e *ArithmeticExpression) ExprType(scope *Environment) DataType {
	lhs := e.Left.ExprType(scope)
	rhs := e.Right.ExprType(scope)

	if lhs.IsInt() && rhs.IsInt() {
		return Integer{}
	}

	return Float{}
}

func (e *ArithmeticExpression) AsAny() interface{} {
	return e
}

func (e *ArithmeticExpression) IsConst() bool {
	return false
}

type ComparisonOperator int

const (
	COGreater ComparisonOperator = iota
	COGreaterEqual
	COLess
	COLessEqual
	COEqual
	CONotEqual
	CONullSafeEqual
)

type ComparisonExpression struct {
	Left     Expression
	Operator ComparisonOperator
	Right    Expression
}

func (e *ComparisonExpression) Kind() ExpressionKind {
	return ExprComparison
}

func (e *ComparisonExpression) ExprType(scope *Environment) DataType {
	if e.Operator == CONullSafeEqual {
		return Integer{}
	} else {
		return Boolean{}
	}
}

func (e *ComparisonExpression) AsAny() interface{} {
	return e
}

func (e *ComparisonExpression) IsConst() bool {
	return false
}

type LikeExpression struct {
	Input   Expression
	Pattern Expression
}

func (e *LikeExpression) Kind() ExpressionKind {
	return ExprLike
}

func (e *LikeExpression) ExprType(scope *Environment) DataType {
	return Boolean{}
}

func (e *LikeExpression) AsAny() interface{} {
	return e
}

func (e *LikeExpression) IsConst() bool {
	return false
}

type GlobExpression struct {
	Input   Expression
	Pattern Expression
}

func (e *GlobExpression) Kind() ExpressionKind {
	return ExprGlob
}

func (e *GlobExpression) ExprType(scope *Environment) DataType {
	return Boolean{}
}

func (e *GlobExpression) AsAny() interface{} {
	return e
}

func (e *GlobExpression) IsConst() bool {
	return false
}

type LogicalOperator int

const (
	LOOr LogicalOperator = iota
	LOAnd
	LOXor
)

type LogicalExpression struct {
	Left     Expression
	Operator LogicalOperator
	Right    Expression
}

func (e *LogicalExpression) Kind() ExpressionKind {
	return ExprLogical
}

func (e *LogicalExpression) ExprType(scope *Environment) DataType {
	return Boolean{}
}

func (e *LogicalExpression) AsAny() interface{} {
	return e
}

func (e *LogicalExpression) IsConst() bool {
	return false
}

type BitwiseOperator int

const (
	BOOr BitwiseOperator = iota
	BOAnd
	BORightShift
	BOLeftShift
)

type BitwiseExpression struct {
	Left     Expression
	Operator BitwiseOperator
	Right    Expression
}

func (e *BitwiseExpression) Kind() ExpressionKind {
	return ExprBitwise
}

func (e *BitwiseExpression) ExprType(scope *Environment) DataType {
	return Integer{}
}

func (e *BitwiseExpression) AsAny() interface{} {
	return e
}

func (e *BitwiseExpression) IsConst() bool {
	return false
}

type CallExpression struct {
	FunctionName  string
	Arguments     []Expression
	IsAggregation bool
}

func (e *CallExpression) Kind() ExpressionKind {
	return ExprCall
}

func (e *CallExpression) ExprType(scope *Environment) DataType {
	prototype := Prototypes[e.FunctionName]

	return prototype.Result
}

func (e *CallExpression) AsAny() interface{} {
	return e
}

func (e *CallExpression) IsConst() bool {
	return false
}

type BetweenExpression struct {
	Value      Expression
	RangeStart Expression
	RangeEnd   Expression
}

func (e *BetweenExpression) Kind() ExpressionKind {
	return ExprBetween
}

func (e *BetweenExpression) ExprType(scope *Environment) DataType {
	return Boolean{}
}

func (e *BetweenExpression) AsAny() interface{} {
	return e
}

func (e *BetweenExpression) IsConst() bool {
	return false
}

type CaseExpression struct {
	Conditions   []Expression
	Values       []Expression
	DefaultValue Expression
	ValuesType   DataType
}

func (e *CaseExpression) Kind() ExpressionKind {
	return ExprCase
}

func (e *CaseExpression) ExprType(scope *Environment) DataType {
	return e.ValuesType
}

func (e *CaseExpression) AsAny() interface{} {
	return e
}

func (e *CaseExpression) IsConst() bool {
	return false
}

type InExpression struct {
	Argument      Expression
	Values        []Expression
	ValuesType    DataType
	HasNotKeyword bool
}

func (e *InExpression) Kind() ExpressionKind {
	return ExprIn
}

func (e *InExpression) ExprType(scope *Environment) DataType {
	return e.ValuesType
}

func (e *InExpression) AsAny() interface{} {
	return e
}

func (e *InExpression) IsConst() bool {
	return false
}

type IsNullExpression struct {
	Argument Expression
	HasNot   bool
}

func (e *IsNullExpression) Kind() ExpressionKind {
	return ExprIsNull
}

func (e *IsNullExpression) ExprType(scope *Environment) DataType {
	return Boolean{}
}

func (e *IsNullExpression) AsAny() interface{} {
	return e
}

func (e *IsNullExpression) IsConst() bool {
	return false
}

type NullExpression struct{}

func (e *NullExpression) Kind() ExpressionKind {
	return ExprNull
}

func (e *NullExpression) ExprType(scope *Environment) DataType {
	return Null{}
}

func (e *NullExpression) AsAny() interface{} {
	return e
}

func (e *NullExpression) IsConst() bool {
	return false
}
