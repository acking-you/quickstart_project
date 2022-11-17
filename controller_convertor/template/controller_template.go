package template

var KControllerDefine = `

type ${class}Controller struct {
	*gin.Context
}
func New${class}Controller(context *gin.Context) *${class}Controller {
	return &${class}Controller{Context: context}
}

`

var KControllerAction = `

func (s *${class}Controller) Do${action}${class}() error {

	return nil
}

`

var KHandlerBind = `

func bind${class}VO(c *gin.Context) error {
	v := &vo.${class}VO{}
	if err := c.ShouldBind(v); err != nil {
		return err
	}
	return nil
}

`

var KNoResponseHandlerByAction = `

func ${action}${class}Handler(c *gin.Context) {
	_ = New${class}Controller(c).Do${action}${class}() 
}

`

var KResponseHandlerByAction = `

func ${action}${class}Handler(c *gin.Context) {
	if err := New${class}Controller(c).Do${action}${class}(); err != nil {
		r.BadResponse(c, err.Error())
	}
}

`
