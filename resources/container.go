package resource

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecr "github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	logs "github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	s3 "github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
)

type ContainerService struct {
	Stack  cdk.Stack
	Vpc    ec2.Vpc
	Config struct {
		RepositoryName string
		BucketName     string
		LogGroupName   string
	}
}

type ContainerResource struct {
	Fargate ecs.FargateService
}

func (s *ContainerService) New() ContainerResource {
	repository := ecr.Repository_FromRepositoryName(s.Stack, jsii.String("ecr"), jsii.String(s.Config.RepositoryName))
	logBucket := s3.Bucket_FromBucketName(s.Stack, jsii.String("bucket"), jsii.String(s.Config.BucketName))
	logGroup := logs.LogGroup_FromLogGroupName(s.Stack, jsii.String("CdkClusterLogAssociation"), jsii.String(s.Config.LogGroupName))

	cluster := ecs.NewCluster(s.Stack, jsii.String("cluster"),
		&ecs.ClusterProps{
			Vpc:         s.Vpc,
			ClusterName: jsii.String("CdkCluster"),
			ExecuteCommandConfiguration: &ecs.ExecuteCommandConfiguration{
				LogConfiguration: &ecs.ExecuteCommandLogConfiguration{
					CloudWatchLogGroup:          logGroup,
					CloudWatchEncryptionEnabled: jsii.Bool(true),
					S3Bucket:                    logBucket,
					S3EncryptionEnabled:         jsii.Bool(true),
					S3KeyPrefix:                 jsii.String("cdk-log"),
				},
				Logging: ecs.ExecuteCommandLogging_OVERRIDE,
			},
		},
	)

	taskDef := ecs.NewFargateTaskDefinition(s.Stack, jsii.String("task"),
		&ecs.FargateTaskDefinitionProps{
			Family:          jsii.String("CdkTaskFamily"),
			Cpu:             jsii.Number(256),
			RuntimePlatform: &ecs.RuntimePlatform{CpuArchitecture: ecs.CpuArchitecture_X86_64()},
		},
	)

	taskDef.AddContainer(jsii.String("addContainer"),
		&ecs.ContainerDefinitionOptions{
			Image:          ecs.ContainerImage_FromEcrRepository(repository, jsii.String("v0.01")),
			Cpu:            jsii.Number(256),
			MemoryLimitMiB: jsii.Number(512),
			ContainerName:  jsii.String("CdkContainer"),
			PortMappings: &[]*ecs.PortMapping{
				{
					Protocol:      ecs.Protocol_TCP,
					HostPort:      jsii.Number(3000),
					ContainerPort: jsii.Number(3000),
				},
			},
			Logging:     ecs.LogDriver_AwsLogs(&ecs.AwsLogDriverProps{StreamPrefix: jsii.String("cdklog"), LogGroup: logGroup}),
			Environment: &map[string]*string{},
		},
	)

	service := ecs.NewFargateService(s.Stack, jsii.String("fargateService"),
		&ecs.FargateServiceProps{
			Cluster:              cluster,
			DeploymentController: &ecs.DeploymentController{Type: "ECS"},
			DesiredCount:         jsii.Number(1),
			ServiceName:          jsii.String("CdkService"),
			TaskDefinition:       taskDef,
			VpcSubnets:           &ec2.SubnetSelection{SubnetType: ec2.SubnetType_PRIVATE_ISOLATED},
		},
	)

	return ContainerResource{Fargate: service}
}
