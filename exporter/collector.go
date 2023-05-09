package exporter

import (
	"reflect"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tris/weatherflow"
)

var (
	// https://weatherflow.github.io/Tempest/api/ws.html

	labelNames = []string{"device_id"}

	desc = &metricDescriptions{
		WindLull: prometheus.NewDesc(
			"weatherflow_wind_lull",
			"Wind lull in meters per second (minimum 3 second sample)",
			labelNames, nil,
		),
		WindAvg: prometheus.NewDesc(
			"weatherflow_wind_avg",
			"Wind speed in meters per second (average over report interval)",
			labelNames, nil,
		),
		WindGust: prometheus.NewDesc(
			"weatherflow_wind_gust",
			"Wind gust in meters per second (maximum 3 second sample)",
			labelNames, nil,
		),
		WindDirectionAvg: prometheus.NewDesc(
			"weatherflow_wind_direction_avg",
			"Wind direction in degrees (average over report interval)",
			labelNames, nil,
		),
		WindSampleInterval: prometheus.NewDesc(
			"weatherflow_wind_sample_interval_seconds",
			"Wind sample interval in seconds",
			labelNames, nil,
		),
		StationPressure: prometheus.NewDesc(
			"weatherflow_station_pressure",
			"Station pressure in millibars",
			labelNames, nil,
		),
		AirTemperature: prometheus.NewDesc(
			"weatherflow_air_temperature",
			"Air temperature in degrees Celsius",
			labelNames, nil,
		),
		RelativeHumidity: prometheus.NewDesc(
			"weatherflow_relative_humidity",
			"Relative humidity in percent",
			labelNames, nil,
		),
		Illuminance: prometheus.NewDesc(
			"weatherflow_illuminance",
			"Illuminance in lux",
			labelNames, nil,
		),
		UV: prometheus.NewDesc(
			"weatherflow_uv",
			"UV index",
			labelNames, nil,
		),
		SolarRadiation: prometheus.NewDesc(
			"weatherflow_solar_radiation",
			"Solar radiation in watts per square meter",
			labelNames, nil,
		),
		RainAccumulated: prometheus.NewDesc(
			"weatherflow_rain_accumulated",
			"Rain accumulated in millimeters",
			labelNames, nil,
		),
		PrecipitationType: prometheus.NewDesc(
			"weatherflow_precipitation_type",
			"Precipitation type (0: none, 1: rain, 2: hail)",
			labelNames, nil,
		),
		LightningStrikeAvgDistance: prometheus.NewDesc(
			"weatherflow_lightning_strike_avg_distance",
			"Lightning strike average distance in kilometers",
			labelNames, nil,
		),
		LightningStrikeCount: prometheus.NewDesc(
			"weatherflow_lightning_strike_total",
			"Lightning strike count",
			labelNames, nil,
		),
		Battery: prometheus.NewDesc(
			"weatherflow_battery_volts",
			"Battery in volts",
			labelNames, nil,
		),
		ReportInterval: prometheus.NewDesc(
			"weatherflow_report_interval_minutes",
			"Report interval in minutes",
			labelNames, nil,
		),
		LocalDailyRainAccumulation: prometheus.NewDesc(
			"weatherflow_local_daily_rain_total",
			"Local daily rain accumulation in millimeters",
			labelNames, nil,
		),
		RainAccumulatedFinal: prometheus.NewDesc(
			"weatherflow_rain_final_total",
			"Rain accumulated final (Rain Check) in millimeters",
			labelNames, nil,
		),
		LocalDailyRainAccumulationFinal: prometheus.NewDesc(
			"weatherflow_local_daily_rain_final_total",
			"Local daily rain accumulation final (Rain Check) in millimeters",
			labelNames, nil,
		),
		PrecipitationAnalysisType: prometheus.NewDesc(
			"weatherflow_precipitation_analysis_type",
			"Precipitation analysis type (0: none, 1: Rain Check with user display on, 2: Rain Check with user display off",
			labelNames, nil,
		),

		// Rapid wind
		WindSpeed: prometheus.NewDesc(
			"weatherflow_wind_speed",
			"Wind speed in meters per second (instant)",
			labelNames, nil,
		),
		WindDirection: prometheus.NewDesc(
			"weatherflow_wind_direction",
			"Wind direction in degrees (instant)",
			labelNames, nil,
		),
	}
)

type metricDescriptions struct {
	WindLull                        *prometheus.Desc
	WindAvg                         *prometheus.Desc
	WindGust                        *prometheus.Desc
	WindDirectionAvg                *prometheus.Desc
	WindSampleInterval              *prometheus.Desc
	StationPressure                 *prometheus.Desc
	AirTemperature                  *prometheus.Desc
	RelativeHumidity                *prometheus.Desc
	Illuminance                     *prometheus.Desc
	UV                              *prometheus.Desc
	SolarRadiation                  *prometheus.Desc
	RainAccumulated                 *prometheus.Desc
	PrecipitationType               *prometheus.Desc
	LightningStrikeAvgDistance      *prometheus.Desc
	LightningStrikeCount            *prometheus.Desc
	Battery                         *prometheus.Desc
	ReportInterval                  *prometheus.Desc
	LocalDailyRainAccumulation      *prometheus.Desc
	RainAccumulatedFinal            *prometheus.Desc
	LocalDailyRainAccumulationFinal *prometheus.Desc
	PrecipitationAnalysisType       *prometheus.Desc

	// Rapid wind
	WindSpeed     *prometheus.Desc
	WindDirection *prometheus.Desc
}

type WeatherCollector struct {
	deviceID int
	metric   struct {
		WindLull                        prometheus.Metric
		WindAvg                         prometheus.Metric
		WindGust                        prometheus.Metric
		WindDirectionAvg                prometheus.Metric
		WindSampleInterval              prometheus.Metric
		StationPressure                 prometheus.Metric
		AirTemperature                  prometheus.Metric
		RelativeHumidity                prometheus.Metric
		Illuminance                     prometheus.Metric
		UV                              prometheus.Metric
		SolarRadiation                  prometheus.Metric
		RainAccumulated                 prometheus.Metric
		PrecipitationType               prometheus.Metric
		LightningStrikeAvgDistance      prometheus.Metric
		LightningStrikeCount            prometheus.Metric
		Battery                         prometheus.Metric
		ReportInterval                  prometheus.Metric
		LocalDailyRainAccumulation      prometheus.Metric
		RainAccumulatedFinal            prometheus.Metric
		LocalDailyRainAccumulationFinal prometheus.Metric
		PrecipitationAnalysisType       prometheus.Metric

		// Rapid wind
		WindSpeed     prometheus.Metric
		WindDirection prometheus.Metric
	}
	timer *time.Timer
}

func NewWeatherCollector(deviceID int) *WeatherCollector {
	return &WeatherCollector{
		deviceID: deviceID,
	}
}

func (wc *WeatherCollector) Describe(ch chan<- *prometheus.Desc) {
	d := reflect.ValueOf(*desc)
	for i := 0; i < d.NumField(); i++ {
		value := d.Field(i).Interface().(*prometheus.Desc)
		if value == nil {
			// TODO move this to init()?
			panic("missing desc for field " + d.Type().Field(i).Name)
		} else {
			ch <- value
		}
	}
}

func (wc *WeatherCollector) Collect(ch chan<- prometheus.Metric) {
	m := reflect.ValueOf(wc.metric)
	for i := 0; i < m.NumField(); i++ {
		value := m.Field(i).Interface()
		if value != nil {
			ch <- value.(prometheus.Metric)
		}
	}
}

func (wc *WeatherCollector) update(msg weatherflow.Message, apiToken string) {
	deviceID := msg.GetDeviceID()

	switch m := msg.(type) {
	case *weatherflow.MessageObsSt:
		timestamp := time.Unix(int64(m.Obs[0].TimeEpoch), 0)
		deviceIDStr := strconv.Itoa(deviceID)
		wc.metric.WindLull = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.WindLull,
			prometheus.GaugeValue,
			m.Obs[0].WindLull,
			deviceIDStr,
		))
		wc.metric.WindAvg = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.WindAvg,
			prometheus.GaugeValue,
			m.Obs[0].WindAvg,
			deviceIDStr,
		))
		wc.metric.WindGust = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.WindGust,
			prometheus.GaugeValue,
			m.Obs[0].WindGust,
			deviceIDStr,
		))
		wc.metric.WindDirectionAvg = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.WindDirectionAvg,
			prometheus.GaugeValue,
			float64(m.Obs[0].WindDirection),
			deviceIDStr,
		))
		wc.metric.WindSampleInterval = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.WindSampleInterval,
			prometheus.GaugeValue,
			float64(m.Obs[0].WindSampleInterval),
			deviceIDStr,
		))
		wc.metric.StationPressure = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.StationPressure,
			prometheus.GaugeValue,
			m.Obs[0].StationPressure,
			deviceIDStr,
		))
		airTemperature := m.Obs[0].AirTemperature
		if airTemperature != nil {
			wc.metric.AirTemperature = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
				desc.AirTemperature,
				prometheus.GaugeValue,
				*airTemperature,
				deviceIDStr,
			))
		}
		relativeHumidity := m.Obs[0].RelativeHumidity
		if relativeHumidity != nil {
			wc.metric.RelativeHumidity = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
				desc.RelativeHumidity,
				prometheus.GaugeValue,
				*relativeHumidity,
				deviceIDStr,
			))
		}
		wc.metric.Illuminance = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.Illuminance,
			prometheus.GaugeValue,
			float64(m.Obs[0].Illuminance),
			deviceIDStr,
		))
		wc.metric.UV = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.UV,
			prometheus.GaugeValue,
			float64(m.Obs[0].UV),
			deviceIDStr,
		))
		wc.metric.SolarRadiation = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.SolarRadiation,
			prometheus.GaugeValue,
			float64(m.Obs[0].SolarRadiation),
			deviceIDStr,
		))
		wc.metric.RainAccumulated = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.RainAccumulated,
			prometheus.CounterValue,
			m.Obs[0].RainAccumulated,
			deviceIDStr,
		))
		wc.metric.PrecipitationType = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.PrecipitationType,
			prometheus.GaugeValue,
			float64(m.Obs[0].PrecipitationType),
			deviceIDStr,
		))
		wc.metric.LightningStrikeAvgDistance = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.LightningStrikeAvgDistance,
			prometheus.GaugeValue,
			float64(m.Obs[0].LightningStrikeAvgDistance),
			deviceIDStr,
		))
		wc.metric.LightningStrikeCount = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.LightningStrikeCount,
			prometheus.CounterValue,
			float64(m.Obs[0].LightningStrikeCount),
			deviceIDStr,
		))
		wc.metric.Battery = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.Battery,
			prometheus.GaugeValue,
			m.Obs[0].Battery,
			deviceIDStr,
		))
		wc.metric.ReportInterval = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.ReportInterval,
			prometheus.GaugeValue,
			float64(m.Obs[0].ReportInterval),
			deviceIDStr,
		))
		wc.metric.LocalDailyRainAccumulation = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.LocalDailyRainAccumulation,
			prometheus.CounterValue,
			m.Obs[0].LocalDailyRainAccumulation,
			deviceIDStr,
		))
		wc.metric.RainAccumulatedFinal = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.RainAccumulatedFinal,
			prometheus.CounterValue,
			m.Obs[0].RainAccumulatedFinal,
			deviceIDStr,
		))
		wc.metric.LocalDailyRainAccumulationFinal = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.LocalDailyRainAccumulationFinal,
			prometheus.CounterValue,
			m.Obs[0].LocalDailyRainAccumulationFinal,
			deviceIDStr,
		))
		wc.metric.PrecipitationAnalysisType = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.PrecipitationAnalysisType,
			prometheus.GaugeValue,
			float64(m.Obs[0].PrecipitationAnalysisType),
			deviceIDStr,
		))

	case *weatherflow.MessageRapidWind:
		timestamp := time.Unix(int64(m.Ob.TimeEpoch), 0)
		deviceIDStr := strconv.Itoa(deviceID)
		wc.metric.WindSpeed = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.WindSpeed,
			prometheus.GaugeValue,
			m.Ob.WindSpeed,
			deviceIDStr,
		))
		wc.metric.WindDirection = prometheus.NewMetricWithTimestamp(timestamp, prometheus.MustNewConstMetric(
			desc.WindDirection,
			prometheus.GaugeValue,
			float64(m.Ob.WindDirection),
			deviceIDStr,
		))
	}
}
