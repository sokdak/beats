package stat

import (
	"github.com/elastic/beats/v7/libbeat/common/cfgwarn"
	"github.com/elastic/beats/v7/metricbeat/helper"
	"github.com/elastic/beats/v7/metricbeat/mb"
	"github.com/elastic/beats/v7/metricbeat/mb/parse"
	"github.com/pkg/errors"
	"io/ioutil"
)

const (
	defaultScheme = "http"
	defaultPath   = "/metrics"
)

var (
	HostParser = parse.URLHostParserBuilder{
		PathConfigKey: "metrics_path",
		DefaultPath:   defaultPath,
		DefaultScheme: defaultScheme,
	}.Build()
)

// init registers the MetricSet with the central registry as soon as the program
// starts. The New function will be called later to instantiate an instance of
// the MetricSet for each host defined in the module's configuration. After the
// MetricSet has been created then Fetch will begin to be called periodically.
func init() {
	mb.Registry.MustAddMetricSet("mining", "stat", New,
		mb.WithHostParser(HostParser),
		mb.DefaultMetricSet())
}

// MetricSet holds any configuration or state information. It must implement
// the mb.MetricSet interface. And this is best achieved by embedding
// mb.BaseMetricSet because it implements all of the required mb.MetricSet
// interface methods except for Fetch.
type MetricSet struct {
	mb.BaseMetricSet
	http *helper.HTTP
}

// New creates a new instance of the MetricSet. New is responsible for unpacking
// any MetricSet specific configuration options if there are any.
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	cfgwarn.Beta("The mining stat metricset is beta.")

	config := struct{}{}
	if err := base.Module().UnpackConfig(&config); err != nil {
		return nil, err
	}

	http, err := helper.NewHTTP(base)
	if err != nil {
		return nil, err
	}

	return &MetricSet{
		BaseMetricSet: base,
		http: http,
	}, nil
}

// Fetch methods implements the data gathering and data conversion to the right
// format. It publishes the event which is then forwarded to the output. In case
// of an error set the Error field of mb.Event or simply call report.Error().
func (m *MetricSet) Fetch(reporter mb.ReporterV2) error {
	// fetch response
	response, err := m.http.FetchResponse()
	if err != nil {
		reporter.Error(errors.Wrapf(err, "unable to fetch data from metrics endpoint"))
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			m.Logger().Debug("error closing http body")
		}
	}()

	// parse events
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// report events
	events := parseResponse(body)
	for _, event := range events {
		reporter.Event(event)
	}
	return nil
}
