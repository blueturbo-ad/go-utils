package environment

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment_ConfPath(t *testing.T) {
	type args struct {
		environment string
		workName    string
		region      string
	}

	var (
		err    error
		inputs = []args{
			{
				environment: "Pro",
				workName:    "go-utils",
				region:      "us-east-1",
			},
		}
	)

	for _, arg := range inputs {
		err = os.Setenv("ENVIRONMENT", arg.environment)
		assert.NoError(t, err)
		err = os.Setenv("REGION", arg.region)
		assert.NoError(t, err)
		err = os.Setenv("WORK_NAME", arg.workName)
		assert.NoError(t, err)

		Init()

		assert.Equal(t, GetConfPath(), filepath.Join(GetWorkPath(), "conf", "us_east_1_conf_dir"))
	}
}

func TestEnv(t *testing.T) {
	t.Run("input=valid", func(t *testing.T) {
		os.Setenv("WORK_NAME", "go-utils")
		Init()
		workPath := GetWorkPath()
		configPath := path.Join(workPath, "config")
		fmt.Println(configPath)
	})
}
