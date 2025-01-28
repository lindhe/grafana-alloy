package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/alloy/internal/component"
	http_service "github.com/grafana/alloy/internal/service/http"
	"github.com/grafana/alloy/syntax"
)

func Test_getCollectors(t *testing.T) {
	t.Run("nothing specified (default behavior)", func(t *testing.T) {
		var exampleDBO11yAlloyConfig = `
		data_source_name = ""
		forward_to = []
	`

		var args Arguments
		err := syntax.Unmarshal([]byte(exampleDBO11yAlloyConfig), &args)
		require.NoError(t, err)

		actualCollectors := getCollectors(args)

		assert.Equal(t, map[string]bool{"QuerySample": false, "SchemaTable": false}, actualCollectors)
	})

	t.Run("enableCollectors", func(t *testing.T) {
		var exampleDBO11yAlloyConfig = `
		data_source_name = ""
		forward_to = []
		enable_collectors = ["QuerySample", "SchemaTable"]
	`

		var args Arguments
		err := syntax.Unmarshal([]byte(exampleDBO11yAlloyConfig), &args)
		require.NoError(t, err)

		actualCollectors := getCollectors(args)

		assert.Equal(t, map[string]bool{"QuerySample": true, "SchemaTable": true}, actualCollectors)
	})
}

func TestName(t *testing.T) {
	var exampleDBO11yAlloyConfig = `
		data_source_name = "sqlmockdb0/"
		forward_to = []
		enable_collectors = ["QuerySample", "SchemaTable"]
	`

	var args Arguments
	err := syntax.Unmarshal([]byte(exampleDBO11yAlloyConfig), &args)
	require.NoError(t, err)

	comp, err := New(component.Options{
		GetServiceData: func(name string) (interface{}, error) { return http_service.Data{}, nil },
		OnStateChange:  func(e component.Exports) {},
	}, args)
	require.NoError(t, err)

	assert.Equal(t, []Collector{}, comp.collectors)
}
