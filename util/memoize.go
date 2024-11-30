package util

type MemoFunction[Input comparable, Output any] func(in Input, rec func(in Input) Output) Output

// Given a function of the form f(n, k) where k is a continuation representing the recursive call,
// memoize it.
func Memoize[Input comparable, Output any](fn MemoFunction[Input, Output]) func(in Input) Output {
	memo := make(map[Input]Output)

	return func(in Input) Output {
		return recurser(memo, fn, in)
	}
}

func recurser[Input comparable, Output any](
	memo map[Input]Output,
	fn MemoFunction[Input, Output],
	in Input) Output {
	val, ok := memo[in]
	if ok {
		return val
	}
	val = fn(in, func(i Input) Output {
		return recurser(memo, fn, i)
	})
	memo[in] = val
	return val
}
