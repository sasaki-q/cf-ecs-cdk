package dns

import (
	cdk "github.com/aws/aws-cdk-go/awscdk/v2"
	acm "github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	lb "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	route53 "github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	tg "github.com/aws/aws-cdk-go/awscdk/v2/awsroute53targets"
	c "github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type PublicHostZoneService struct {
	Stack      cdk.Stack
	DomainName string
}

type PublicHostZoneResource struct {
	Cert acm.ICertificate
	Zone route53.HostedZone
}

func (s *PublicHostZoneService) New() PublicHostZoneResource {
	zone := route53.NewPublicHostedZone(s.Stack, jsii.String("route53"), &route53.PublicHostedZoneProps{
		ZoneName: jsii.String(s.DomainName),
	})
	myacm := acm.NewCertificate(s.Stack, jsii.String("acm"), &acm.CertificateProps{
		DomainName:      jsii.String(s.DomainName),
		CertificateName: jsii.String("CdkAcm"),
		Validation:      acm.CertificateValidation_FromDns(zone),
	})

	return PublicHostZoneResource{Cert: myacm, Zone: zone}
}

func NewARecord(
	stack c.Construct,
	zone route53.IHostedZone,
	alb lb.ILoadBalancerV2,
) {
	route53.NewARecord(stack, jsii.String("arecord"), &route53.ARecordProps{
		Zone:   zone,
		Target: route53.RecordTarget_FromAlias(tg.NewLoadBalancerTarget(alb)),
	})
}
