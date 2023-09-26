package main

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bitrise-io/go-utils/log"
	"github.com/stretchr/testify/require"
)

// Example file:
//
// # Generated by pub on 2019-08-06 13:38:51.234081.
// archive:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/archive-2.0.8/lib/
// args:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/args-1.5.0/lib/
// async:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/async-2.2.0/lib/
// boolean_selector:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/boolean_selector-1.0.4/lib/
// charcode:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/charcode-1.1.2/lib/
// collection:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/collection-1.14.11/lib/
// convert:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/convert-2.1.1/lib/
// crypto:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/crypto-2.0.6/lib/
// cupertino_icons:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/cupertino_icons-0.1.2/lib/
// dart_config:file:///Users/vagrant/.pub-cache/git/dart_config-a7ed88a4793e094a4d5d5c2d88a89e55510accde/lib/
// flutter:file:///Users/vagrant/flutter-sdk/flutter/packages/flutter/lib/
// flutter_launcher_icons:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/flutter_launcher_icons-0.7.0/lib/
// flutter_test:file:///Users/vagrant/flutter-sdk/flutter/packages/flutter_test/lib/
// font_awesome_flutter:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/font_awesome_flutter-8.4.0/lib/
// image:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/image-2.0.7/lib/
// intl:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/intl-0.15.7/lib/
// matcher:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/matcher-0.12.5/lib/
// meta:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/meta-1.1.6/lib/
// path:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/path-1.6.2/lib/
// pedantic:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/pedantic-1.7.0/lib/
// petitparser:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/petitparser-2.1.1/lib/
// quiver:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/quiver-2.0.3/lib/
// scoped_model:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/scoped_model-1.0.1/lib/
// shared_preferences:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/shared_preferences-0.4.3/lib/
// sky_engine:file:///Users/vagrant/flutter-sdk/flutter/bin/cache/pkg/sky_engine/lib/
// source_span:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/source_span-1.5.5/lib/
// stack_trace:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/stack_trace-1.9.3/lib/
// stream_channel:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/stream_channel-2.0.0/lib/
// string_scanner:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/string_scanner-1.0.4/lib/
// term_glyph:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/term_glyph-1.1.0/lib/
// test_api:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/test_api-0.2.5/lib/
// typed_data:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/typed_data-1.1.6/lib/
// vector_math:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/vector_math-2.0.8/lib/
// xml:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/xml-3.3.1/lib/
// yaml:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/yaml-2.1.15/lib/
// veggieseasons:lib/

func Test_parsePackageResolutionFile(t *testing.T) {
	var analyzerPath url.URL
	analyzerPath.Path = "/Users/vagrant/.pub-cache/hosted/pub.dartlang.org/analyzer-0.36.4/lib/"
	analyzerPath.Scheme = "file"

	var relPath url.URL
	relPath.Path = "../../.pub-cache/hosted/pub.dartlang.org/analyzer-0.36.4/lib/"

	tests := []struct {
		name     string
		contents string
		want     map[string]url.URL
		wantErr  bool
	}{
		{
			name: "empty file",
			contents: `# Generated by pub on 2019-08-05 14:50:08.261783.

# Other comment`,
			want:    map[string]url.URL{},
			wantErr: false,
		},
		{
			name: "package with file scheme",
			contents: `# Generated by pub on 2019-08-05 14:50:08.261783.
analyzer:file:///Users/vagrant/.pub-cache/hosted/pub.dartlang.org/analyzer-0.36.4/lib/`,
			want: map[string]url.URL{
				"analyzer": analyzerPath,
			},
			wantErr: false,
		},
		{
			name: "relative path dependency",
			contents: `# Generated by pub on 2019-08-05 14:50:08.261783.
analyzer:../../.pub-cache/hosted/pub.dartlang.org/analyzer-0.36.4/lib/`,
			want: map[string]url.URL{
				"analyzer": relPath,
			},
			wantErr: false,
		},
		{
			name: "invalid URI",
			contents: `# Generated by pub on 2019-08-05 14:50:08.261783.
analyzer::invalid/ss`,
			want:    map[string]url.URL{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePackageResolutionFile(tt.contents)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePackageResolutionFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePackageResolutionFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cacheableFlutterDepPaths(t *testing.T) {
	log.SetEnableDebugLog(true)

	var yamlURL url.URL
	yamlURL.Path = "/Users/vagrant/.pub-cache/hosted/pub.dartlang.org/yaml-2.1.15/lib/"

	var gitURL url.URL
	gitURL.Path = "/Users/vagrant/.pub-cache/git/sample-apps-flutter-sample-pub-package-afc598ac6dc1a8e39ff7fd505463fd1df9f7c600/mypath/lib/"

	var gitURL2 url.URL
	gitURL2.Path = "/Users/vagrant/.pub-cache/git/sample-apps-flutter-ios-android-package-f44f5a21cd47f45db70faa1a6aed8c8035483d73/lib"

	tests := []struct {
		name              string
		packageToLocation map[string]url.URL
		want              []string
		wantErr           bool
	}{
		{
			name: "valid dependency from system cache",
			packageToLocation: map[string]url.URL{
				"yaml": yamlURL,
			},
			want: []string{
				"/Users/vagrant/.pub-cache/hosted/pub.dartlang.org/yaml-2.1.15",
			},
			wantErr: false,
		},
		{
			name: "package from git",
			packageToLocation: map[string]url.URL{
				"sample_package": gitURL,
			},
			want: []string{
				"/Users/vagrant/.pub-cache/git",
			},
			wantErr: false,
		},
		{
			name: "multiple packages from git",
			packageToLocation: map[string]url.URL{
				"sample_package":  gitURL,
				"sample_package2": gitURL2,
			},
			want: []string{
				"/Users/vagrant/.pub-cache/git",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cacheableFlutterDepPaths(tt.packageToLocation)
			if (err != nil) != tt.wantErr {
				t.Errorf("cacheableFlutterDepPaths() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cacheableFlutterDepPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONFormatParsing(t *testing.T) {
	testData := `
{
  "configVersion": 2,
  "packages": [
    {
      "name": "insert-cool-name-here",
      "rootUri": "file:///path/to/file",
      "packageUri": "lib/",
      "languageVersion": "2.12"
    }
  ],
  "generated": "2023-04-26T10:11:25.639598Z",
  "generator": "pub",
  "generatorVersion": "2.15.1"
}
`
	result, err := parseJSON(testData)
	require.NoError(t, err)

	fileURL, err := url.Parse(filepath.Join("file:///path/to/file", "lib/"))
	require.NoError(t, err)

	expected := map[string]url.URL{
		"insert-cool-name-here": *fileURL,
	}
	assert.Equal(t, expected, result)
}
