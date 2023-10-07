package resource

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	acm "github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	lb "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	"github.com/aws/jsii-runtime-go"
)

type LoadBalancerService struct {
	Stack          cdk.Stack
	Cert           acm.ICertificate
	Vpc            ec2.Vpc
	FargateService ecs.FargateService
}

type LoadBalancerResource struct {
	Alb lb.ApplicationLoadBalancer
}

func (s *LoadBalancerService) New() LoadBalancerResource {
	alb := lb.NewApplicationLoadBalancer(s.Stack, jsii.String("alb"), &lb.ApplicationLoadBalancerProps{
		Vpc:              s.Vpc,
		InternetFacing:   jsii.Bool(true),
		LoadBalancerName: jsii.String("CdkAlb"),
		VpcSubnets: &ec2.SubnetSelection{
			SubnetType: ec2.SubnetType_PUBLIC,
		},
	})

	tg := lb.NewApplicationTargetGroup(s.Stack, jsii.String("tg"), &lb.ApplicationTargetGroupProps{
		TargetGroupName: jsii.String("CdkTargetGroup"),
		TargetType:      lb.TargetType_IP,
		Port:            jsii.Number(80),
		Vpc:             s.Vpc,
		Targets:         &[]lb.IApplicationLoadBalancerTarget{s.FargateService},
	})

	alb.AddListener(jsii.String("listener"), &lb.BaseApplicationListenerProps{
		Certificates:        &[]lb.IListenerCertificate{lb.ListenerCertificate_FromCertificateManager(s.Cert)},
		DefaultTargetGroups: &[]lb.IApplicationTargetGroup{tg},
		Protocol:            lb.ApplicationProtocol_HTTPS,
	})

	return LoadBalancerResource{Alb: alb}
}
