package envlib

import (
	"sync"

	"github.com/dogmatiq/ferrite"
)

// AWSSDKv2 is a [ferrite.Registry] containing the environment variables used by
// the AWS Go SDK (v2).
//
// See https://github.com/aws/aws-sdk-go-v2/blob/main/config/env_config.go
func AWSSDKv2() ferrite.Registry {
	return build(
		"aws.sdk-v2",
		"AWS SDK for Go (v2)",
		"https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/configure-gosdk.html",
		func(reg ferrite.Registry) {
			accessKeyID := ferrite.
				String("AWS_ACCESS_KEY_ID", "AWS access key ID associated with an IAM account").
				WithExample("AKIAIOSFODNN7EXAMPLE", "general format of an AWS access key ID").
				Optional(ferrite.WithRegistry(reg))

			ferrite.
				String("AWS_SECRET_ACCESS_KEY", "AWS IAM secret access key").
				WithSensitiveContent().
				Optional(
					ferrite.WithRegistry(reg),
					ferrite.RelevantIf(accessKeyID),
				)
		},
	)
}

var registries sync.Map

func build(
	name string,
	description string,
	url string,
	setup func(ferrite.Registry),
) ferrite.Registry {
	if v, ok := registries.Load(name); ok {
		return v.(ferrite.Registry)
	}

	reg := ferrite.NewRegistry(
		name,
		description,
		ferrite.WithDocumentationURL(url),
	)

	setup(reg)

	v, _ := registries.LoadOrStore(name, reg)
	return v.(ferrite.Registry)
}
