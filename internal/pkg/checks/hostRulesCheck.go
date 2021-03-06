/*
Copyright 2019 The Kubernetes Authors.

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

package checks

import (
	"fmt"

	"sigs.k8s.io/ingress-controller-conformance/internal/pkg/k8s"
)

func init() {
	Checks.AddCheck(hostRulesCheck)
}

var hostRulesCheck = &Check{
	Name:        "host-rules",
	Description: "Ingress with host rule should send traffic to the correct backend service",
	Run: func(check *Check, config Config) (success bool, err error) {
		host, err := k8s.GetIngressHost("default", "host-rules")
		if err != nil {
			return
		}

		resp, err := captureRequest(fmt.Sprintf("http://%s", host), "foo.bar.com")
		if err != nil {
			return
		}

		a := new(assertionSet)
		a.equals(assert{resp.StatusCode, 200, "Expected StatusCode to be %s but was %s"})
		a.equals(assert{resp.TestId, "host-rules", "Expected the responding service would be '%s' but was '%s'"})
		a.equals(assert{resp.Host, "foo.bar.com", "Expected the request host would be '%s' but was '%s'"})

		// TODO: Implement more assertions on request Headers for example

		if a.Error() == "" {
			success = true
		} else {
			fmt.Print(a)
		}
		return
	},
}
