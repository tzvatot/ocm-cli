/*
Copyright (c) 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file contains functions that add common flags to the command line.

package flags

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/openshift-online/ocm-sdk-go"
	"github.com/spf13/pflag"

	"github.com/openshift-online/ocm-cli/pkg/debug"
)

// AddDebugFlag adds the '--debug' flag to the given set of command line flags.
func AddDebugFlag(fs *pflag.FlagSet) {
	debug.AddFlag(fs)
}

// AddParameterFlag adds the '--parameter' flag to the given set of command line flags.
func AddParameterFlag(fs *pflag.FlagSet, values *[]string) {
	fs.StringArrayVar(
		values,
		"parameter",
		nil,
		"Query parameters to add to the request. The value must be the name of the "+
			"parameter, followed by an optional equals sign and then the value "+
			"of the parameter. Can be used multiple times to specify multiple "+
			"parameters or multiple values for the same parameter. Example: "+
			"--parameter search=\"username like 'myname%'\"",
	)
}

// AddHeaderFlag adds the '--header' flag to the given set of command line flags.
func AddHeaderFlag(fs *pflag.FlagSet, values *[]string) {
	fs.StringArrayVar(
		values,
		"header",
		nil,
		"Headers to add to the request. The value must be the name of the header "+
			"followed by an optional equals sign and then the value of the "+
			"header. Can be used multiple times to specify multiple headers "+
			"or multiple values for the same header.",
	)
}

// AddBodyFlag adds the '--body' flag to the given set of command line flags.
func AddBodyFlag(fs *pflag.FlagSet, value *string) {
	fs.StringVar(
		value,
		"body",
		"",
		"Name fo the file containing the request body. If this isn't given then "+
			"the body will be taken from the standard input.",
	)
}

// ApplyParameterFlag applies the value of the '--parameter' command line flag to the given
// request.
func ApplyParameterFlag(request interface{}, values []string) {
	applyNVFlag(request, "Parameter", values)
}

// ApplyHeaderFlag applies the value of the '--header' command line flag to the given request.
func ApplyHeaderFlag(request interface{}, values []string) {
	applyNVFlag(request, "Header", values)
}

// applyNVFlag finds the method with the given name in a request and calls it to set a collection of
// name value pairs.
func applyNVFlag(request interface{}, method string, values []string) {
	// Find the method:
	callable := reflect.ValueOf(request).MethodByName(method)
	if !callable.IsValid() {
		return
	}

	// Split the values into name value pairs and call the method for each one:
	for _, value := range values {
		var name string
		position := strings.Index(value, "=")
		if position != -1 {
			name = value[:position]
			value = value[position+1:]
		} else {
			name = value
			value = ""
		}
		args := []reflect.Value{
			reflect.ValueOf(name),
			reflect.ValueOf(value),
		}
		callable.Call(args)
	}
}

// ApplyBodyFlag applies the value of the '--body' command line flag to the given request.
func ApplyBodyFlag(request *sdk.Request, value string) error {
	var body []byte
	var err error
	if value != "" {
		// #nosec G304
		body, err = ioutil.ReadFile(value)
	} else {
		body, err = ioutil.ReadAll(os.Stdin)
	}
	if err != nil {
		return err
	}
	request.Bytes(body)
	return nil
}
