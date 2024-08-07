// Code generated by counterfeiter. DO NOT EDIT.
package loggingfakes

import (
	"sync"

	"github.com/gsmcwhirter/go-util/v12/logging"
)

type FakeLogger struct {
	ErrStub        func(string, error, ...interface{})
	errMutex       sync.RWMutex
	errArgsForCall []struct {
		arg1 string
		arg2 error
		arg3 []interface{}
	}
	LogStub        func(...interface{}) error
	logMutex       sync.RWMutex
	logArgsForCall []struct {
		arg1 []interface{}
	}
	logReturns struct {
		result1 error
	}
	logReturnsOnCall map[int]struct {
		result1 error
	}
	MessageStub        func(string, ...interface{})
	messageMutex       sync.RWMutex
	messageArgsForCall []struct {
		arg1 string
		arg2 []interface{}
	}
	PrintfStub        func(string, ...interface{})
	printfMutex       sync.RWMutex
	printfArgsForCall []struct {
		arg1 string
		arg2 []interface{}
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeLogger) Err(arg1 string, arg2 error, arg3 ...interface{}) {
	fake.errMutex.Lock()
	fake.errArgsForCall = append(fake.errArgsForCall, struct {
		arg1 string
		arg2 error
		arg3 []interface{}
	}{arg1, arg2, arg3})
	stub := fake.ErrStub
	fake.recordInvocation("Err", []interface{}{arg1, arg2, arg3})
	fake.errMutex.Unlock()
	if stub != nil {
		fake.ErrStub(arg1, arg2, arg3...)
	}
}

func (fake *FakeLogger) ErrCallCount() int {
	fake.errMutex.RLock()
	defer fake.errMutex.RUnlock()
	return len(fake.errArgsForCall)
}

func (fake *FakeLogger) ErrCalls(stub func(string, error, ...interface{})) {
	fake.errMutex.Lock()
	defer fake.errMutex.Unlock()
	fake.ErrStub = stub
}

func (fake *FakeLogger) ErrArgsForCall(i int) (string, error, []interface{}) {
	fake.errMutex.RLock()
	defer fake.errMutex.RUnlock()
	argsForCall := fake.errArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeLogger) Log(arg1 ...interface{}) error {
	fake.logMutex.Lock()
	ret, specificReturn := fake.logReturnsOnCall[len(fake.logArgsForCall)]
	fake.logArgsForCall = append(fake.logArgsForCall, struct {
		arg1 []interface{}
	}{arg1})
	stub := fake.LogStub
	fakeReturns := fake.logReturns
	fake.recordInvocation("Log", []interface{}{arg1})
	fake.logMutex.Unlock()
	if stub != nil {
		return stub(arg1...)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeLogger) LogCallCount() int {
	fake.logMutex.RLock()
	defer fake.logMutex.RUnlock()
	return len(fake.logArgsForCall)
}

func (fake *FakeLogger) LogCalls(stub func(...interface{}) error) {
	fake.logMutex.Lock()
	defer fake.logMutex.Unlock()
	fake.LogStub = stub
}

func (fake *FakeLogger) LogArgsForCall(i int) []interface{} {
	fake.logMutex.RLock()
	defer fake.logMutex.RUnlock()
	argsForCall := fake.logArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeLogger) LogReturns(result1 error) {
	fake.logMutex.Lock()
	defer fake.logMutex.Unlock()
	fake.LogStub = nil
	fake.logReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeLogger) LogReturnsOnCall(i int, result1 error) {
	fake.logMutex.Lock()
	defer fake.logMutex.Unlock()
	fake.LogStub = nil
	if fake.logReturnsOnCall == nil {
		fake.logReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.logReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeLogger) Message(arg1 string, arg2 ...interface{}) {
	fake.messageMutex.Lock()
	fake.messageArgsForCall = append(fake.messageArgsForCall, struct {
		arg1 string
		arg2 []interface{}
	}{arg1, arg2})
	stub := fake.MessageStub
	fake.recordInvocation("Message", []interface{}{arg1, arg2})
	fake.messageMutex.Unlock()
	if stub != nil {
		fake.MessageStub(arg1, arg2...)
	}
}

func (fake *FakeLogger) MessageCallCount() int {
	fake.messageMutex.RLock()
	defer fake.messageMutex.RUnlock()
	return len(fake.messageArgsForCall)
}

func (fake *FakeLogger) MessageCalls(stub func(string, ...interface{})) {
	fake.messageMutex.Lock()
	defer fake.messageMutex.Unlock()
	fake.MessageStub = stub
}

func (fake *FakeLogger) MessageArgsForCall(i int) (string, []interface{}) {
	fake.messageMutex.RLock()
	defer fake.messageMutex.RUnlock()
	argsForCall := fake.messageArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeLogger) Printf(arg1 string, arg2 ...interface{}) {
	fake.printfMutex.Lock()
	fake.printfArgsForCall = append(fake.printfArgsForCall, struct {
		arg1 string
		arg2 []interface{}
	}{arg1, arg2})
	stub := fake.PrintfStub
	fake.recordInvocation("Printf", []interface{}{arg1, arg2})
	fake.printfMutex.Unlock()
	if stub != nil {
		fake.PrintfStub(arg1, arg2...)
	}
}

func (fake *FakeLogger) PrintfCallCount() int {
	fake.printfMutex.RLock()
	defer fake.printfMutex.RUnlock()
	return len(fake.printfArgsForCall)
}

func (fake *FakeLogger) PrintfCalls(stub func(string, ...interface{})) {
	fake.printfMutex.Lock()
	defer fake.printfMutex.Unlock()
	fake.PrintfStub = stub
}

func (fake *FakeLogger) PrintfArgsForCall(i int) (string, []interface{}) {
	fake.printfMutex.RLock()
	defer fake.printfMutex.RUnlock()
	argsForCall := fake.printfArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeLogger) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.errMutex.RLock()
	defer fake.errMutex.RUnlock()
	fake.logMutex.RLock()
	defer fake.logMutex.RUnlock()
	fake.messageMutex.RLock()
	defer fake.messageMutex.RUnlock()
	fake.printfMutex.RLock()
	defer fake.printfMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeLogger) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ logging.Logger = new(FakeLogger)
