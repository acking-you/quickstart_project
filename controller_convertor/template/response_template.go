package template

var KResponseTemplate = `
const (
	KStatusOk = iota
	KStatusErr
)

type ResponseMsg[T any] struct {
	StatusCode int    ` + "`json:\"status_code\"`" + `
	StatusMsg  string ` + "`json:\"status_msg\"`" + `
	Msg        *T     ` + "`json:\"msg\"`" + `
}

func BadResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"status_code": KStatusErr,
		"status_msg":  msg,
	})
}

func GoodResponse[T any](c *gin.Context, msg *T) {
	c.JSON(http.StatusOK, &ResponseMsg[T]{
		StatusCode: KStatusOk,
		StatusMsg:  "response OK!",
		Msg:        msg,
	})
}
`
