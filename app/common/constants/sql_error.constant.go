package constants

type SqlQueryStatusType string

const (
	Success        SqlQueryStatusType = "success"
	NotUnique      SqlQueryStatusType = "not_unique"
	NotFound       SqlQueryStatusType = "not_found"
	Unknown        SqlQueryStatusType = "unknown"
	Error          SqlQueryStatusType = "error"
	NoDataToUpdate SqlQueryStatusType = "no_data_to_update"
)
