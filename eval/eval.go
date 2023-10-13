package eval

import (
	"github.com/thingsme/thingscript/ast"
	"github.com/thingsme/thingscript/object"
)

var (
	NULL = &object.Null{}
)

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func isBreak(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.BREAK_OBJ
	}
	return false
}

func isReturn(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.RETURN_VALUE_OBJ
	}
	return false
}

func isTruthy(obj object.Object) bool {
	if obj == NULL {
		return false
	}
	if t, ok := obj.(*object.Boolean); ok {
		return t.Value
	}
	return true
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.VarStatement:
		return evalVarStatement(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.BreakStatement:
		return &object.Break{}
	case *ast.AssignStatement:
		val, ok := env.Get(node.Name.Value)
		if !ok {
			return object.Errorf("identifier not found: %s", node.Name.Value)
		}
		evaluated := Eval(node.Value, env)
		if isError(evaluated) {
			return evaluated
		}
		return evalAssignStatement(val, evaluated)
	case *ast.OperAssignStatement:
		left, ok := env.Get(node.Name.Value)
		if !ok {
			return object.Errorf("identifier not found: %s", node.Name.Value)
		}
		right := Eval(node.Value, env)
		if isError(right) {
			return right
		}
		evaluated := evalInfixExpression(node.Operator, left, right)
		if isError(evaluated) {
			return evaluated
		}
		return evalAssignStatement(left, evaluated)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.WhileExpression:
		return evalWhileExpression(node, env)
	case *ast.DoWhileExpression:
		return evalDoWhileExpression(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ImmediateIfExpression:
		return evalImmediateIfExpression(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return evalCallFunction(function, args)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.AccessExpression:
		return evalAccessExpression(node, env)
	case *ast.FunctionStatement:
		params := node.Parameters
		body := node.Body
		val := &object.Function{Parameters: params, Env: env, Body: body}
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.Boolean:
		return &object.Boolean{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.HashMapLiteral:
		return evalHashLiteral(node, env)
	}
	return nil
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ || rt == object.BREAK_OBJ {
				return result
			}
		}
	}
	return result
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalVarStatement(node *ast.VarStatement, env *object.Environment) object.Object {
	var evaluated object.Object
	if node.Value != nil {
		evaluated = Eval(node.Value, env)
		if isError(evaluated) {
			return evaluated
		}
	}
	if node.TypeDecl == nil {
		// infer the var type from value
		env.Set(node.Name.Value, evaluated)
	} else {
		// explicitly declare the type of the var
		if node.TypeDecl.Package == nil {
			evaluated = env.Type("", node.TypeDecl.Name.Value, evaluated)
		} else {
			evaluated = env.Type(node.TypeDecl.Package.Value, node.TypeDecl.Name.Value, evaluated)
		}
		if isError(evaluated) {
			return evaluated
		}
		env.Set(node.Name.Value, evaluated)
	}
	return nil
}

func evalCallFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		if ret := fn.Func(args...); ret != nil {
			return ret
		}
		return NULL
	default:
		return object.Errorf("not a function: %s", fn.Type())
	}
}

func evalAccessExpression(exp *ast.AccessExpression, env *object.Environment) object.Object {
	left := Eval(exp.Left, env)
	if isError(left) {
		return left
	}

	switch r := exp.Right.(type) {
	case *ast.Identifier:
		fn := left.Member(r.Value)
		if fn == nil {
			return object.Errorf("function %q not found in %q", r.Value, left.Type())
		}
		ret := fn(left)
		return ret
	case *ast.CallExpression:
		fnIdent, ok := r.Function.(*ast.Identifier)
		if !ok {
			return object.Errorf("undefined %q in %q", r.Function.String(), left.Type())
		}
		fn := left.Member(fnIdent.Value)
		if fn == nil {
			return object.Errorf("function %q not found in %q", fnIdent.Value, left.Type())
		}
		args := evalExpressions(r.Arguments, env)
		ret := fn(left, args...)
		return ret
	default:
		return object.Errorf("invalid access operator %q.(%T)", left.Type(), r)
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalWhileExpression(we *ast.WhileExpression, env *object.Environment) object.Object {
	for {
		condition := Eval(we.Condition, env)
		if isError(condition) {
			return condition
		}
		if !isTruthy(condition) {
			break
		}
		ret := Eval(we.Block, env)
		if isError(ret) {
			return ret
		}
		if isBreak(ret) {
			break
		}
		if isReturn(ret) {
			return ret
		}
	}
	return nil
}

func evalDoWhileExpression(we *ast.DoWhileExpression, env *object.Environment) object.Object {
	for {
		ret := Eval(we.Block, env)
		if isError(ret) {
			return ret
		}
		if isBreak(ret) {
			break
		}
		if isReturn(ret) {
			return ret
		}
		condition := Eval(we.Condition, env)
		if isError(condition) {
			return condition
		}
		if !isTruthy(condition) {
			break
		}
	}
	return nil
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	for n, condExpression := range ie.Condition {
		condition := Eval(condExpression, env)
		if isError(condition) {
			return condition
		}
		if isTruthy(condition) {
			return Eval(ie.Consequence[n], env)
		}
	}
	if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}
	return NULL
}

func evalImmediateIfExpression(ie *ast.ImmediateIfExpression, env *object.Environment) object.Object {
	leftVal := Eval(ie.Left, env)
	if isError(leftVal) {
		return leftVal
	}
	if leftVal != NULL {
		return leftVal
	}
	return Eval(ie.Right, env)
}

func evalIndexExpression(left object.Object, index object.Object) object.Object {
	if operFunc := left.Member("["); operFunc != nil {
		return operFunc(left, index)
	} else {
		return object.Errorf("index operation not supported: %s", left.Type())
	}
}

func evalHashLiteral(node *ast.HashMapLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return object.Errorf("unusable as hash key: %s", key.Type())
		}
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}
	return &object.HashMap{Pairs: pairs}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	if right == nil {
		return &object.Boolean{Value: true}
	}
	if t, ok := right.(*object.Boolean); ok {
		return &object.Boolean{Value: !t.Value}
	}
	return &object.Boolean{Value: false}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if node.Value == "nil" {
		return NULL
	}
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin := env.Builtin(node.Value); builtin != nil {
		return builtin
	}
	return object.Errorf("identifier not found: " + node.Value)
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	switch right.Type() {
	case object.INTEGER_OBJ:
		value := right.(*object.Integer).Value
		return &object.Integer{Value: value * -1}
	case object.FLOAT_OBJ:
		value := right.(*object.Float).Value
		return &object.Float{Value: value * -1}
	default:
		return object.Errorf("unknown operator: -%s", right.Type())
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return object.Errorf("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	if opFunc := left.Member(operator); opFunc != nil {
		if ret := opFunc(left, right); ret != nil {
			return ret
		}
	}
	return object.Errorf("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalAssignStatement(left object.Object, right object.Object) object.Object {
	if assignFunc := left.Member("="); assignFunc != nil {
		left = assignFunc(left, right)
	} else {
		left = nil
	}

	if left == nil {
		return object.Errorf("unable to set value of %T with %T", left, right)
	}
	return nil
}
