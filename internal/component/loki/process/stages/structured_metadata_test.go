package stages

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/require"

	"github.com/grafana/alloy/internal/featuregate"
	"github.com/grafana/loki/pkg/push"
	util_log "github.com/grafana/loki/v3/pkg/util/log"
)

var pipelineStagesStructuredMetadataFromLogfmt = `
stage.logfmt {
	mapping = { "app" = ""}
}

stage.structured_metadata { 
	values = {"app" = ""}
}
`

var pipelineStagesStructuredMetadataFromJSON = `
stage.json {
	expressions = {app = ""}
}

stage.structured_metadata { 
	values = {"app" = ""}
}
`

var pipelineStagesStructuredMetadataWithRegexParser = `
stage.regex {
	expression = "^(?s)(?P<time>\\S+?) (?P<stream>stdout|stderr) (?P<flags>\\S+?) (?P<content>.*)$"
}

stage.structured_metadata { 
	values = {"stream" = ""}
}
`

var pipelineStagesStructuredMetadataFromJSONWithTemplate = `
stage.json {
	expressions = {app = ""}
}

stage.template {
    source   = "app"
    template = "{{ ToUpper .Value }}"
}

stage.structured_metadata { 
	values = {"app" = ""}
}
`

var pipelineStagesStructuredMetadataAndRegularLabelsFromJSON = `
stage.json {
	expressions = {app = "", component = "" }
}

stage.structured_metadata { 
	values = {"app" = ""}
}

stage.labels { 
	values = {"component" = ""}
}
`

var pipelineStagesStructuredMetadataFromStaticLabels = `
stage.static_labels {
	values = {"component" = "querier", "pod" = "loki-querier-664f97db8d-qhnwg"}
}

stage.structured_metadata {
	values = {"pod" = ""}
}
`

var pipelineStagesStructuredMetadataFromStaticLabelsDifferentKey = `
stage.static_labels {
	values = {"component" = "querier", "pod" = "loki-querier-664f97db8d-qhnwg"}
}

stage.structured_metadata {
	values = {"pod_name" = "pod"}
}
`

func Test_StructuredMetadataStage(t *testing.T) {
	tests := map[string]struct {
		pipelineStagesYaml         string
		logLine                    string
		expectedStructuredMetadata push.LabelsAdapter
		expectedLabels             model.LabelSet
	}{
		"expected structured metadata to be extracted with logfmt parser and to be added to entry": {
			pipelineStagesYaml:         pipelineStagesStructuredMetadataFromLogfmt,
			logLine:                    "app=loki component=ingester",
			expectedStructuredMetadata: push.LabelsAdapter{push.LabelAdapter{Name: "app", Value: "loki"}},
		},
		"expected structured metadata to be extracted with json parser and to be added to entry": {
			pipelineStagesYaml:         pipelineStagesStructuredMetadataFromJSON,
			logLine:                    `{"app":"loki" ,"component":"ingester"}`,
			expectedStructuredMetadata: push.LabelsAdapter{push.LabelAdapter{Name: "app", Value: "loki"}},
		},
		"expected structured metadata to be extracted with regexp parser and to be added to entry": {
			pipelineStagesYaml:         pipelineStagesStructuredMetadataWithRegexParser,
			logLine:                    `2019-01-01T01:00:00.000000001Z stderr P i'm a log message!`,
			expectedStructuredMetadata: push.LabelsAdapter{push.LabelAdapter{Name: "stream", Value: "stderr"}},
		},
		"expected structured metadata to be extracted with json parser and to be added to entry after rendering the template": {
			pipelineStagesYaml:         pipelineStagesStructuredMetadataFromJSONWithTemplate,
			logLine:                    `{"app":"loki" ,"component":"ingester"}`,
			expectedStructuredMetadata: push.LabelsAdapter{push.LabelAdapter{Name: "app", Value: "LOKI"}},
		},
		"expected structured metadata and regular labels to be extracted with json parser and to be added to entry": {
			pipelineStagesYaml:         pipelineStagesStructuredMetadataAndRegularLabelsFromJSON,
			logLine:                    `{"app":"loki" ,"component":"ingester"}`,
			expectedStructuredMetadata: push.LabelsAdapter{push.LabelAdapter{Name: "app", Value: "loki"}},
			expectedLabels:             model.LabelSet{model.LabelName("component"): model.LabelValue("ingester")},
		},
		"expected structured metadata and regular labels to be extracted with static labels stage and to be added to entry": {
			pipelineStagesYaml:         pipelineStagesStructuredMetadataFromStaticLabels,
			logLine:                    `sample log line`,
			expectedStructuredMetadata: push.LabelsAdapter{push.LabelAdapter{Name: "pod", Value: "loki-querier-664f97db8d-qhnwg"}},
			expectedLabels:             model.LabelSet{model.LabelName("component"): model.LabelValue("querier")},
		},
		"expected structured metadata and regular labels to be extracted with static labels stage using different structured key": {
			pipelineStagesYaml:         pipelineStagesStructuredMetadataFromStaticLabelsDifferentKey,
			logLine:                    `sample log line`,
			expectedStructuredMetadata: push.LabelsAdapter{push.LabelAdapter{Name: "pod_name", Value: "loki-querier-664f97db8d-qhnwg"}},
			expectedLabels:             model.LabelSet{model.LabelName("component"): model.LabelValue("querier")},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			pl, err := NewPipeline(util_log.Logger, loadConfig(test.pipelineStagesYaml), nil, prometheus.DefaultRegisterer, featuregate.StabilityGenerallyAvailable)
			require.NoError(t, err)

			result := processEntries(pl, newEntry(nil, nil, test.logLine, time.Now()))[0]
			require.Equal(t, test.expectedStructuredMetadata, result.StructuredMetadata)
			if test.expectedLabels != nil {
				require.Equal(t, test.expectedLabels, result.Labels)
			} else {
				require.Empty(t, result.Labels)
			}
		})
	}
}
