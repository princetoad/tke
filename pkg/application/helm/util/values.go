/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package util

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/strvals"
	"sigs.k8s.io/yaml"
)

// MergeValues merges values from values and rawValues
// values: can specify multiple or separate values: key1=val1,key2=val2
// rawValues: json format or yaml format
func MergeValues(values []string, rawValues string, vType string) (map[string]interface{}, error) {
	base := map[string]interface{}{}
	var err error
	if len(rawValues) > 0 {
		switch strings.ToLower(vType) {
		case "yaml", "":
			if base, err = parseYAMLValue(rawValues); err != nil {
				return nil, err
			}
			break
		case "json":
			if base, err = parseJSONValue(rawValues); err != nil {
				return nil, err
			}
			break
		default:
			return nil, errors.New(fmt.Sprintf("unsupport value type: %s", vType))
		}
	}
	if len(values) > 0 {
		for _, value := range values {
			if err := strvals.ParseInto(value, base); err != nil {
				return nil, errors.Wrap(err, "failed parsing value data")
			}
		}
	}
	return base, err
}

func parseValue(values []string) (map[string]interface{}, error) {
	base := map[string]interface{}{}
	for _, value := range values {
		if err := strvals.ParseInto(value, base); err != nil {
			return nil, errors.Wrap(err, "failed parsing value data")
		}
	}
	return base, nil
}

func parseYAMLValue(values string) (map[string]interface{}, error) {
	// rv := fmt.Sprintf("%q", values)
	base := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(values), &base); err != nil {
		return nil, errors.Wrap(err, "failed to parse yaml value")
	}
	return base, nil
}

func parseJSONValue(values string) (map[string]interface{}, error) {
	// rv := fmt.Sprintf("%q", values)
	base := map[string]interface{}{}
	if err := json.Unmarshal([]byte(values), &base); err != nil {
		return nil, errors.Wrap(err, "failed to parse json value")
	}
	return base, nil
}
