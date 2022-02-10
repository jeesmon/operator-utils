module github.com/jeesmon/operator-utils

go 1.15

require (
	github.com/openshift/custom-resource-status v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	k8s.io/api v0.23.0
	k8s.io/apimachinery v0.23.0
	k8s.io/client-go v0.23.0
	sigs.k8s.io/controller-runtime v0.11.0
)

replace github.com/openshift/custom-resource-status => github.com/jeesmon/custom-resource-status v1.1.1
