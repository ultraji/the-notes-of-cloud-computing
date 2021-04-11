package singleton

import "sync"

// 饿汉
type singleton1 struct{}

var instance1 *singleton1

func init() {
	instance1 = &singleton1{}
}

func GetInstance1() *singleton1 {
	return instance1
}

// 懒汉（带锁）
type singleton2 struct{}

var instance2 *singleton2
var mu2 sync.Mutex

func GetInstance2() *singleton2 {
	mu2.Lock()
	defer mu2.Unlock()
	if instance2 == nil {
		instance2 = &singleton2{}
	}
	return instance2
}

// 双重检测
type singleton3 struct{}

var instance3 *singleton3
var mu3 sync.Mutex

func GetInstance3() *singleton3 {
	if instance3 == nil { // <-- Not yet perfect. since it's not fully atomic
		mu3.Lock()
		defer mu3.Unlock()
		if instance3 == nil {
			instance3 = &singleton3{}
		}
	}
	return instance3
}

// sync.Once
type singleton4 struct{}

var instance4 *singleton4
var once sync.Once

func GetInstance4() *singleton4 {
	once.Do(func() {
		instance4 = &singleton4{}
	})
	return instance4
}
