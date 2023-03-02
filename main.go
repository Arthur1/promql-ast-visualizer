package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/prometheus/prometheus/promql/parser"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	exprStr := scanner.Text()

	expr, err := parse(exprStr)
	if err != nil {
		log.Fatalln(err)
	}

	if err := printExpr(expr, 0); err != nil {
		log.Fatalln(err)
	}
}

func parse(exprStr string) (parser.Expr, error) {
	expr, err := parser.ParseExpr(exprStr)
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func printExpr(expr parser.Expr, level int) error {
	switch e := expr.(type) {
	case *parser.BinaryExpr:
		printfWithIndent(level, "%v %s\n", color.GreenString("BinaryExpr"), e.String())
		if err := printExpr(e.LHS, level+1); err != nil {
			return err
		}
		if err := printExpr(e.RHS, level+1); err != nil {
			return err
		}
	case *parser.UnaryExpr:
		printfWithIndent(level, "%v %s\n", color.GreenString("UnaryExpr"), e.String())
		if err := printExpr(e.Expr, level+1); err != nil {
			return err
		}
	case *parser.AggregateExpr:
		printfWithIndent(level, "%v %s\n", color.GreenString("AggregateExpr"), e.String())
		if err := printExpr(e.Expr, level+1); err != nil {
			return err
		}
	case *parser.Call:
		printfWithIndent(level, "%v %s\n", color.GreenString("Call"), e.String())
		for _, v := range e.Args {
			if err := printExpr(v, level+1); err != nil {
				return err
			}
		}
	case *parser.MatrixSelector:
		printfWithIndent(level, "%v %s\n", color.GreenString("MatrixSelector"), e.String())
		if err := printExpr(e.VectorSelector, level+1); err != nil {
			return err
		}
	case *parser.SubqueryExpr:
		printfWithIndent(level, "%v %s\n", color.GreenString("SubqueryExpr"), e.String())
		if err := printExpr(e.Expr, level+1); err != nil {
			return err
		}
	case *parser.NumberLiteral:
		printfWithIndent(level, "%v %s\n", color.GreenString("NumberLiteral"), e.String())
	case *parser.ParenExpr:
		printfWithIndent(level, "%v %s\n", color.GreenString("ParenExpr"), e.String())
		if err := printExpr(e.Expr, level+1); err != nil {
			return err
		}
	case *parser.StringLiteral:
		printfWithIndent(level, "%v %s\n", color.GreenString("StringLiteral"), e.String())
	case *parser.StepInvariantExpr:
		printfWithIndent(level, "%v %s\n", color.GreenString("StepInvariantExpr"), e.String())
	case *parser.VectorSelector:
		printfWithIndent(level, "%v %s\n", color.GreenString("VectorSelector"), e.String())
	default:
		return fmt.Errorf("unsupported expr: %s", expr.String())
	}

	return nil
}

func printfWithIndent(level int, format string, a ...any) {
	fmt.Print(strings.Repeat(" ", level))
	fmt.Printf(format, a...)
}
