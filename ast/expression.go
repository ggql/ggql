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
	ExpressionKind() ExpressionKind
	ExprType(scope *Environment) DataType
	AsAny() interface{}
}

func IsConst(expr Expression) bool {
	switch expr.ExpressionKind() {
	case ExprNumber, ExprBoolean, ExprString:
		return true
	default:
		return false
	}
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

func (s *StringExpression) ExpressionKind() ExpressionKind {
	return ExprString
}

func (s *StringExpression) ExprType(scope *Environment) DataType {
	switch s.ValueType {
	case StringValueText:
		return Text
	case StringValueTime:
		return Time
	case StringValueDate:
		return Date
	case StringValueDateTime:
		return DateTime
	default:
		return Undefined
	}
}

func (s *StringExpression) AsAny() interface{} {
	return s
}

type SymbolExpression struct {
	Value string
}

func (s *SymbolExpression) ExpressionKind() ExpressionKind {
	return ExprSymbol
}

func (s *SymbolExpression) ExprType(scope *Environment) DataType {
	if scope.Contains(s.Value) {
		return scope.Scopes[s.Value].Clone()
	}

	if typ, ok := tablesFieldsTypes[s.Value]; ok {
		return typ.Clone()
	}

	return Undefined
}

func (s *SymbolExpression) AsAny() interface{} {
	return s
}

type GlobalVariableExpression struct {
	Name string
}

func (s *GlobalVariableExpression) ExpressionKind() ExpressionKind {
	return ExprGlobalVariable
}

func (s *GlobalVariableExpression) ExprType(scope *Environment) DataType {
	if scope.Contains(s.Name) {
		return scope.Scopes[s.Name].Clone()
	}

	if typ, ok := tablesFieldsTypes[s.Name]; ok {
		return typ.Clone()
	}

	return Undefined
}

func (s *GlobalVariableExpression) AsAny() interface{} {
	return s
}

type NumberExpression struct {
	Value Value
}

func (s *NumberExpression) ExpressionKind() ExpressionKind {
	return ExprNumber
}

func (s *NumberExpression) ExprType(scope *Environment) DataType {
	return s.Value.DataType()
}

func (s *NumberExpression) AsAny() interface{} {
	return s
}

type BooleanExpression struct {
	IsTrue bool
}

func (s *BooleanExpression) ExpressionKind() ExpressionKind {
	return ExprBoolean
}

func (s *BooleanExpression) ExprType(scope *Environment) DataType {
	return Boolean
}

func (s *BooleanExpression) AsAny() interface{} {
	return s
}

type PrefixUnaryOperator int

const (
	Minus PrefixUnaryOperator = iota
	Bang
)

type PrefixUnary struct {
	Right Expression
	Op    PrefixUnaryOperator
}

func (s *PrefixUnary) ExpressionKind() ExpressionKind {
	return ExprPrefixUnary
}

func (s *PrefixUnary) ExprType(scope *Environment) DataType {
	if s.Op == Bang {
		return Boolean
	} else {
		return Integer
	}
}

func (s *PrefixUnary) AsAny() interface{} {
	return s
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

func (s *ArithmeticExpression) ExpressionKind() ExpressionKind {
	return ExprArithmetic
}

func (s *ArithmeticExpression) ExprType(scope *Environment) DataType {
	lhs := s.Left.ExprType(scope)
	rhs := s.Right.ExprType(scope)

	if lhs.IsInt() && rhs.IsInt() {
		return Integer
	}

	return Float
}

func (s *ArithmeticExpression) AsAny() interface{} {
	return s
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

func (s *ComparisonExpression) ExpressionKind() ExpressionKind {
	return ExprComparison
}

func (s *ComparisonExpression) ExprType(scope *Environment) DataType {
	if s.Operator == CONullSafeEqual {
		return Integer
	} else {
		return Boolean
	}
}

func (s *ComparisonExpression) AsAny() interface{} {
	return s
}

type LikeExpression struct {
	Input   Expression
	Pattern Expression
}

func (s *LikeExpression) ExpressionKind() ExpressionKind {
	return ExprLike
}

func (s *LikeExpression) ExprType(scope *Environment) DataType {
	return Boolean
}

func (s *LikeExpression) AsAny() interface{} {
	return s
}

type GlobExpression struct {
	Input   Expression
	Pattern Expression
}

func (s *GlobExpression) ExpressionKind() ExpressionKind {
	return ExprGlob
}

func (s *GlobExpression) ExprType(scope *Environment) DataType {
	return Boolean
}

func (s *GlobExpression) AsAny() interface{} {
	return s
}

type LogicalOperator int

const (
	Or LogicalOperator = iota
	And
	Xor
)

type LogicalExpression struct {
	Left     Expression
	Operator LogicalOperator
	Right    Expression
}

func (s *LogicalExpression) ExpressionKind() ExpressionKind {
	return ExprLogical
}

func (s *LogicalExpression) ExprType(scope *Environment) DataType {
	return Boolean
}

func (s *LogicalExpression) AsAny() interface{} {
	return s
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

func (s *BitwiseExpression) ExpressionKind() ExpressionKind {
	return ExprBitwise
}

func (s *BitwiseExpression) ExprType(scope *Environment) DataType {
	return Integer
}

func (s *BitwiseExpression) AsAny() interface{} {
	return s
}

type CallExpression struct {
	FunctionName  string
	Arguments     []Expression
	IsAggregation bool
}

func (s *CallExpression) ExpressionKind() ExpressionKind {
	return ExprCall
}

func (s *CallExpression) ExprType(scope *Environment) DataType {
	prototype := Prototypes[s.FunctionName]
	return prototype.Result.Clone()
}

func (s *CallExpression) AsAny() interface{} {
	return s
}

type BetweenExpression struct {
	Value      Expression
	RangeStart Expression
	RangeEnd   Expression
}

func (s *BetweenExpression) ExpressionKind() ExpressionKind {
	return ExprBetween
}

func (s *BetweenExpression) ExprType(scope *Environment) DataType {
	return Boolean
}

func (s *BetweenExpression) AsAny() interface{} {
	return s
}

type CaseExpression struct {
	Conditions   []Expression
	Values       []Expression
	DefaultValue Expression
	ValuesType   DataType
}

func (s *CaseExpression) ExpressionKind() ExpressionKind {
	return ExprCase
}

func (s *CaseExpression) ExprType(scope *Environment) DataType {
	return Boolean
}

func (s *CaseExpression) AsAny() interface{} {
	return s
}

type InExpression struct {
	Argument   Expression
	Values     []Expression
	ValuesType DataType
}

func (s *InExpression) ExpressionKind() ExpressionKind {
	return ExprIn
}

func (s *InExpression) ExprType(scope *Environment) DataType {
	return s.ValuesType.Clone()
}

func (s *InExpression) AsAny() interface{} {
	return s
}

type IsNullExpression struct {
	Argument Expression
	HasNot   bool
}

func (s *IsNullExpression) ExpressionKind() ExpressionKind {
	return ExprIsNull
}

func (s *IsNullExpression) ExprType(scope *Environment) DataType {
	return Boolean
}

func (s *IsNullExpression) AsAny() interface{} {
	return s
}

type NullExpression struct{}

func (s *NullExpression) ExpressionKind() ExpressionKind {
	return ExprNull
}

func (s *NullExpression) ExprType(scope *Environment) DataType {
	return Null
}

func (s *NullExpression) AsAny() interface{} {
	return s
}
