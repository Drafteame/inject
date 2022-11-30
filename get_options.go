package inject

type getOptions struct {
	singleton bool
	shared    bool
}

type GetOption func(*getOptions)

func Singleton() GetOption {
	return func(o *getOptions) {
		o.singleton = true
	}
}

func FromShared() GetOption {
	return func(o *getOptions) {
		o.shared = true
	}
}
