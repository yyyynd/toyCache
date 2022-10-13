package server

type Loader interface {
	Load(key string) (value interface{}, err error)
}

type LoadFun func(key string)(value interface{}, err error)

func (f LoadFun) Load(key string)(value interface{}, err error) {
	return f(key)
}

func SetLoadFunc(fn Loader) {
	LoadFunc = fn
}