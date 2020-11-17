// ref: https://dzone.com/articles/try-and-catch-in-golang
package jadesdk

type Exception interface{}

func Throw(up Exception) {
	panic(up)
}

type TryCatchBlock struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

func (tcb TryCatchBlock) Do() {
	if tcb.Try == nil {
		return
	}
	if tcb.Finally != nil {
		defer tcb.Finally()
	}
	if tcb.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcb.Catch(r)
			}
		}()
	}
	tcb.Try()
}
