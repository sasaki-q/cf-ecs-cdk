package main

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"

	resource "cdk/resources"
)

type IResource[T any] interface {
	New() T
}

type ContainerConfig struct {
	repository string
	log        struct {
		bucket   string
		logGroup string
	}
}

type Config struct {
	domainname string
	beConfig   ContainerConfig
}

type CdkStackProps struct {
	awscdk.StackProps
}

const ARG = "ENV"

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps, config Config) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	var networkService IResource[resource.NetworkResource] = &resource.NetworkService{Stack: stack}
	var networkResource = networkService.New()
	var dnsService IResource[resource.DnsResource] = &resource.DnsService{Stack: stack, DomainName: config.domainname}
	var dnsResource = dnsService.New()
	var containerService IResource[resource.ContainerResource] = &resource.ContainerService{
		Stack: stack,
		Vpc:   networkResource.Vpc,
		Config: struct {
			RepositoryName string
			BucketName     string
			LogGroupName   string
		}{
			RepositoryName: config.beConfig.repository,
			BucketName:     config.beConfig.log.bucket,
			LogGroupName:   config.beConfig.log.logGroup,
		},
	}
	var containerResource = containerService.New()
	var albService IResource[resource.LoadBalancerResource] = &resource.LoadBalancerService{
		Stack:          stack,
		Cert:           dnsResource.Cert,
		Vpc:            networkResource.Vpc,
		FargateService: containerResource.Fargate,
	}
	var albResource = albService.New()
	resource.NewARecord(stack, dnsResource.Zone, albResource.Alb)
	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	env := app.Node().TryGetContext(jsii.String(fmt.Sprintf("%s", ARG)))
	if env == nil {
		panic("please pass context")
	}
	envVal, ok := app.Node().TryGetContext(jsii.String(fmt.Sprintf("%s", env))).(map[string]interface{})
	if !ok {
		panic("please set context")
	}

	awscdk.Tags_Of(app).Add(jsii.String("Project"), jsii.String("CDK-GO"), nil)
	awscdk.Tags_Of(app).Add(jsii.String("Env"), jsii.String(fmt.Sprintf("%s", env)), nil)
	beConfig := envVal["beconf"].(map[string]interface{})

	NewCdkStack(app, "CdkStack",
		&CdkStackProps{
			awscdk.StackProps{
				Env: &awscdk.Environment{
					Account: jsii.String(envVal["account"].(string)),
					Region:  jsii.String(envVal["region"].(string)),
				},
				Synthesizer: awscdk.NewDefaultStackSynthesizer(
					&awscdk.DefaultStackSynthesizerProps{
						FileAssetsBucketName: jsii.String(envVal["bucketname"].(string)),
					},
				),
			},
		},
		Config{
			domainname: envVal["domainname"].(string),
			beConfig: ContainerConfig{
				repository: beConfig["repository"].(string),
				log: struct {
					bucket   string
					logGroup string
				}{
					bucket:   beConfig["log"].(map[string]interface{})["bucket"].(string),
					logGroup: beConfig["log"].(map[string]interface{})["group"].(string),
				},
			},
		},
	)

	app.Synth(nil)
}
