package flags

import (
	"log/slog"
	"reflect"
	"testing"

	"github.com/spf13/pflag"
)

func TestStruct(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	var cfg Config
	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Struct(&cfg, NamePrefix("test"), EnvPrefix("TEST"))(flagSet)
	flagSet.Parse([]string{"--help"})
}

func TestSet(t *testing.T) {
	var cfg Config
	rv := reflect.ValueOf(&cfg).Elem()
	rt := rv.Type()

	for i := range rt.NumField() {
		field := rt.Field(i)
		ft, fv := field.Type, rv.Field(i)

		t.Logf("%-18s type: %-18s nil: %t canset: %t addr: %t value: %v", field.Name, ft, IsNil(fv), fv.CanSet(), fv.CanAddr(), fv.Interface())

		if IsNil(fv) {
			t.Logf(" -- %s is nil, then create instance", field.Name)
			switch ft.Kind() {
			case reflect.Pointer:
				fv.Set(reflect.New(ft.Elem()))
			case reflect.Slice:
				fv.Set(reflect.MakeSlice(ft, 0, 0))
			}
			t.Logf(" -- %s created, value: %v", field.Name, fv.Interface())
		}
	}
}

type Config struct {
	Version         string   `json:"version,omitempty" yaml:"version,omitempty" flag:"version" usage:"sing-box version prefix" short:"v"`
	Workdir         string   `json:"workdir,omitempty" yaml:"workdir,omitempty" flag:"workdir" usage:"workdir" short:"d"`
	Template        string   `json:"template,omitempty" yaml:"template,omitempty" flag:"template" usage:"template url or path" short:"t"`
	Outbound        string   `json:"outbound,omitempty" yaml:"outbound,omitempty" flag:"outbound" usage:"outbound url or path" short:"o"`
	ConfigUpdate    string   `json:"config_update,omitempty" yaml:"config_update,omitempty" flag:"config_update" usage:"config auto update cron expr" default:"@every 24h"`
	BinUpdate       string   `json:"bin_update,omitempty" yaml:"bin_update,omitempty" flag:"bin_update" usage:"in auto update cron expr" default:"@every 48h"`
	DashboardUpdate string   `json:"dashboard_update,omitempty" yaml:"dashboard_update,omitempty" flag:"dashboard_update" usage:"dashboard auto update cron expr" default:"@every 48h"`
	AutoRestart     string   `json:"auto_restart,omitempty" yaml:"auto_restart,omitempty" flag:"auto_restart" usage:"auto restart cron expr" default:"@every 1h"`
	DownloadProxy   string   `json:"download_proxy,omitempty" yaml:"download_proxy,omitempty" flag:"download_proxy" usage:"download proxy"`
	ApiProxy        string   `json:"api_proxy,omitempty" yaml:"api_proxy,omitempty" flag:"api_proxy" usage:"api proxy"`
	Inbounds        []string `json:"inbounds,omitempty" yaml:"inbounds,omitempty" flag:"in" usage:"inbounds"`

	Logger *Logger
}

type Logger struct {
	*RotateOptions `json:",inline" yaml:",inline"`
	Stderr         *RotateOptions `json:"stderr,omitempty" yaml:"stderr,omitempty"`
	Stdout         *RotateOptions `json:"stdout,omitempty" yaml:"stdout,omitempty"`
}

type RotateOptions struct {
	Path       string `json:"path,omitempty" yaml:"path,omitempty"`               // 文件路径
	Std        string `json:"std,omitempty" yaml:"std,omitempty"`                 // 标准输出
	MaxSize    int64  `json:"max_size,omitempty" yaml:"max_size,omitempty"`       // 单文件最大大小
	MaxBackups int    `json:"max_backups,omitempty" yaml:"max_backups,omitempty"` // 最大备份文件数量
}

func IsNil(v reflect.Value) bool {
	k := v.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	}
	return false
}
