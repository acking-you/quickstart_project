package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"quickstart/util"
	"testing"
)

func TestSnackCase2CamelCase(t *testing.T) {
	assert.Equal(t, "testSnackCase2CamelCase", util.SnackCase2CamelCase("test_snack_case2_camel_case"))
}

func TestSnackCase2PascalCase(t *testing.T) {
	assert.Equal(t, "TestSnackCase2CamelCase", util.SnackCase2PascalCase("test_snack_case2_camel_case"))
}

func TestCamelcase2SnackCase(t *testing.T) {
	assert.Equal(t, "test_base_you_at_your_time", util.CamelCase2SnackCase("testBaseYouAtYourTime"))

	assert.Equal(t, "do_parser_struct_name", util.CamelCase2SnackCase("doParserStructName"))
}

func TestPascalCase2SnackCase(t *testing.T) {
	assert.Equal(t, "test_camelcase2_snack_case", util.PascalCase2SnackCase("TestCamelcase2SnackCase"))
}

func TestSsanf(t *testing.T) {
	var s1, s2, s3 string
	//成功
	err := util.Sscanf("nihao(sfdsfsa safsd)", "$($)$", &s1, &s2, &s3)
	fmt.Println(s1, "  ", s2, " ", s3)
	if err != nil {
		panic(err)
	}
	//失败
	err = util.Sscanf("nihao(sfdsfsa safsd", "$($)$", &s1, &s2, &s3)
	if err != nil {
		panic(err)
	}
	fmt.Println(s1, "  ", s2, " ", s3)
}
