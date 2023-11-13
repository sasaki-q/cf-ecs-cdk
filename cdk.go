package main

import (
	resource "cdk/resources"
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const ARG = "ENV"

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

type GithubConfig struct {
	owner       string
	repository  string
	accessToken string
}

type Config struct {
	domainName   string
	githubConfig GithubConfig
	beConfig     ContainerConfig
}

type CdkStackProps struct {
	awscdk.StackProps
}

var SubnetOptions = map[string]string{
	"Public":     "Public",
	"Private":    "Private",
	"DB_Private": "DB_Private",
}

const (
	PublicSubnet    = "Public"
	PrivateSubnet   = "Private"
	DBPrivateSubnet = "DB_Private"
)

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps, config Config) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	var networkService IResource[resource.NetworkResource] = &resource.NetworkService{Stack: stack, SubnetOptions: SubnetOptions}
	networkService.New()
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
	ecsConfig := envVal["ecs_conf"].(map[string]interface{})
	ghConfig := envVal["github_conf"].(map[string]interface{})

	NewCdkStack(app, "MyCdkStack",
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
			domainName: envVal["domain_name"].(string),
			githubConfig: GithubConfig{
				owner:       ghConfig["owner"].(string),
				repository:  ghConfig["repository"].(string),
				accessToken: ghConfig["access_token"].(string),
			},
			beConfig: ContainerConfig{
				repository: ecsConfig["repository"].(string),
				log: struct {
					bucket   string
					logGroup string
				}{
					bucket:   ecsConfig["log"].(map[string]interface{})["bucket"].(string),
					logGroup: ecsConfig["log"].(map[string]interface{})["group"].(string),
				},
			},
		},
	)

	app.Synth(nil)
}

func serviceGenerator() {}
