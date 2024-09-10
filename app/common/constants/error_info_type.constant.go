package constants

type errorInfoTypeType struct {
	NotUnique  string
	NotFound   string
	NotAllowed string
}

var ErrorInfoTypeConstant = errorInfoTypeType{
	NotUnique:  "not_unique",
	NotFound:   "not_found",
	NotAllowed: "not_allowed",
}
