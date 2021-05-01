package misc

import (
  "errors"
  "fmt"
  "go/ast"
  "go/parser"
  "go/token"
  "reflect"
  "strconv"
  "strings"
)

func ParseAndEval(exp string) (float64, error) {
  tree, err := parser.ParseExpr(exp)
  if err != nil {
    return 0, err
  }
  return eval(tree)
}

func eval(tree ast.Expr) (float64, error) {
  switch n := tree.(type) {
    case *ast.BasicLit:
      if n.Kind == token.INT {
        var i int64

        if strings.HasPrefix(n.Value, "0x") || strings.HasPrefix(n.Value, "0X") {
          i, _ = strconv.ParseInt(n.Value[2:], 16, 64)
        } else if strings.HasPrefix(n.Value, "0b") || strings.HasPrefix(n.Value, "0B") {
          i, _ = strconv.ParseInt(n.Value[2:], 2, 64)
        } else if strings.HasPrefix(n.Value, "0") {
          i, _ = strconv.ParseInt(n.Value[1:], 8, 64)
        } else {
          i, _ = strconv.ParseInt(n.Value, 10, 64)
        }

        return float64(i), nil
      } else if n.Kind == token.FLOAT {
        i, _ := strconv.ParseFloat(n.Value, 64)
        return float64(i), nil
      }
      return unsupported(n.Kind)
    case *ast.BinaryExpr:
      switch n.Op {
      case token.ADD, token.SUB, token.MUL, token.QUO:
      default:
        return unsupported(n.Op)
      }
      x, err := eval(n.X)
      if err != nil {
        return 0, err
      }
      y, err := eval(n.Y)
      if err != nil {
        return 0, err
      }
      switch n.Op {
      case token.ADD:
        return x + y, nil
      case token.SUB:
        return x - y, nil
      case token.MUL:
        return x * y, nil
      case token.QUO:
        return x / y, nil
      }
    case *ast.ParenExpr:
      return eval(n.X)
  }

  return unsupported(reflect.TypeOf(tree))
}

func unsupported(i interface{}) (float64, error) {
  return 0.0, errors.New(fmt.Sprintf("%v unsupportedported", i))
}
