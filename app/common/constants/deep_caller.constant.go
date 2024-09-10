package constants

type deepCallerType struct {
	Modules    int
	Common     int
	AppRoot    int
	CommonRoot int
}

var DeepCallerConstant = deepCallerType{
	Modules:    1,
	Common:     1,
	AppRoot:    0,
	CommonRoot: 1,
}
