package jq

// Evaluator evaluates jq queries.
type Evaluator interface {
	Evaluate(string, interface{}) ([]interface{}, error)
}

// QueryEvaluator compiles queries and evaluates JSON input.
type QueryEvaluator struct {
	compiler Compiler
}

// NewQueryEvaluator returns a new QueryEvaluator using compiler.
func NewQueryEvaluator(compiler Compiler) *QueryEvaluator {
	return &QueryEvaluator{
		compiler: compiler,
	}
}

// Evaluate compiles the raw query string rawQuery and evaluates input.
// It returns a slice of evaluation results.
// If the query fails to compile, or input fails to evaluate, it errors.
func (e *QueryEvaluator) Evaluate(rawQuery string, input interface{}) ([]interface{}, error) {
	code, err := e.compiler(rawQuery)
	if err != nil {
		return nil, &QueryEvaluationError{Err: err}
	}
	var results []interface{}
	iter := code.Run(input)
	for {
		result, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := result.(error); ok {
			return nil, &QueryEvaluationError{Err: err}
		}
		results = append(results, result)
	}
	return results, nil
}
