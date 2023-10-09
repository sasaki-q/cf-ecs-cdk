package resource

import (
	"fmt"

	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	cf "github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	origin "github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	s3 "github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	s3_deployment "github.com/aws/aws-cdk-go/awscdk/v2/awss3deployment"
	"github.com/aws/jsii-runtime-go"
)

type CloudfrontService struct {
	Stack cdk.Stack
}

type CloudfrontResource struct{}

func (s *CloudfrontService) New() CloudfrontResource {
	oai := cf.NewOriginAccessIdentity(s.Stack, jsii.String("oai"), &cf.OriginAccessIdentityProps{
		Comment: jsii.String("my-website-oai"),
	})

	bucket := s3.NewBucket(s.Stack, jsii.String("bucket"), &s3.BucketProps{
		BucketName: jsii.String("my-app-bucket-2023-10-09"),
	})

	policy := iam.NewPolicyStatement(&iam.PolicyStatementProps{
		Actions: &[]*string{jsii.String("s3:GetObject")},
		Effect:  iam.Effect_ALLOW,
		Principals: &[]iam.IPrincipal{
			iam.NewCanonicalUserPrincipal(oai.CloudFrontOriginAccessIdentityS3CanonicalUserId()),
		},
		Resources: &[]*string{jsii.String(fmt.Sprintf("%s/*", *bucket.BucketArn()))},
	})

	bucket.AddToResourcePolicy(policy)

	distribution := cf.NewDistribution(s.Stack, jsii.String("destribution"), &cf.DistributionProps{
		DefaultRootObject: jsii.String("index.html"),
		DefaultBehavior: &cf.BehaviorOptions{
			AllowedMethods:       cf.AllowedMethods_ALLOW_GET_HEAD(),
			CachedMethods:        cf.CachedMethods_CACHE_GET_HEAD(),
			CachePolicy:          cf.CachePolicy_CACHING_OPTIMIZED(),
			ViewerProtocolPolicy: cf.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
			Origin: origin.NewS3Origin(bucket, &origin.S3OriginProps{
				OriginAccessIdentity: oai,
			}),
		},
		ErrorResponses: &[]*cf.ErrorResponse{
			{
				Ttl:                cdk.Duration_Seconds(jsii.Number(10)),
				HttpStatus:         jsii.Number(404),
				ResponseHttpStatus: jsii.Number(404),
				ResponsePagePath:   jsii.String("/error.html"),
			},
		},
	})

	s3_deployment.NewBucketDeployment(s.Stack, jsii.String("bucket_deployment"), &s3_deployment.BucketDeploymentProps{
		DestinationBucket: bucket,
		Distribution:      distribution,
		DistributionPaths: &[]*string{jsii.String("/*")},
		Sources: &[]s3_deployment.ISource{
			s3_deployment.Source_Data(
				jsii.String("/index.html"),
				jsii.String("<h1>Hello World!</h1>"),
			),
			s3_deployment.Source_Data(
				jsii.String("/error.html"),
				jsii.String("<h1>Not Found !</h1>"),
			),
		},
	})

	return CloudfrontResource{}
}
