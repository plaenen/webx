package utils

import (
	"fmt"
	"strings"
)

// DatastarExpression builds composable multi-statement Datastar expressions.
type DatastarExpression struct {
	statements []string
}

// NewExpression creates a new expression builder.
func NewExpression() *DatastarExpression {
	return &DatastarExpression{
		statements: make([]string, 0),
	}
}

// Statement adds a raw statement to the expression.
func (e *DatastarExpression) Statement(stmt string) *DatastarExpression {
	if stmt != "" {
		e.statements = append(e.statements, stmt)
	}
	return e
}

// SetSignal adds a signal assignment: $signal = value
func (e *DatastarExpression) SetSignal(signal, value string) *DatastarExpression {
	return e.Statement(fmt.Sprintf("$%s = %s", signal, value))
}

// Conditional adds a ternary expression.
func (e *DatastarExpression) Conditional(condition, trueExpr, falseExpr string) *DatastarExpression {
	if falseExpr == "" {
		falseExpr = "null"
	}
	return e.Statement(fmt.Sprintf("%s ? %s : %s", condition, trueExpr, falseExpr))
}

// Build joins all statements with "; " and returns the final expression.
func (e *DatastarExpression) Build() string {
	if len(e.statements) == 0 {
		return ""
	}
	return strings.Join(e.statements, "; ")
}

// BuildConditional creates a standalone ternary expression.
func BuildConditional(condition, trueExpr, falseExpr string) string {
	if falseExpr == "" {
		falseExpr = "null"
	}
	return fmt.Sprintf("%s ? %s : %s", condition, trueExpr, falseExpr)
}
