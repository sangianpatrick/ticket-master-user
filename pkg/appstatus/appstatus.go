package appstatus

type Status string

var (
	OK                  Status = Status("OK")
	NotFound            Status = Status("NotFound")
	InternalServerError Status = Status("InternalServerError")
)
