// +build ignore

package main

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

const filename = `list_tags_gen.go`

var serviceNames = []string{
	"accessanalyzer",
	"acm",
	"acmpca",
	"amplify",
	"appmesh",
	"appstream",
	"appsync",
	"athena",
	"backup",
	"cloudhsmv2",
	"cloudtrail",
	"cloudwatch",
	"cloudwatchevents",
	"cloudwatchlogs",
	"codecommit",
	"codedeploy",
	"codepipeline",
	"cognitoidentity",
	"cognitoidentityprovider",
	"configservice",
	"databasemigrationservice",
	"dataexchange",
	"datasync",
	"dax",
	"devicefarm",
	"directconnect",
	"directoryservice",
	"dlm",
	"docdb",
	"dynamodb",
	"ecr",
	"ecs",
	"efs",
	"eks",
	"elasticache",
	"elasticbeanstalk",
	"elasticsearchservice",
	"elbv2",
	"firehose",
	"fsx",
	"glue",
	"guardduty",
	"greengrass",
	"imagebuilder",
	"inspector",
	"iot",
	"iotanalytics",
	"iotevents",
	"kafka",
	"kinesis",
	"kinesisanalytics",
	"kinesisanalyticsv2",
	"kms",
	"lambda",
	"licensemanager",
	"mediaconnect",
	"mediaconvert",
	"medialive",
	"mediapackage",
	"mediastore",
	"mq",
	"neptune",
	"opsworks",
	"organizations",
	"qldb",
	"rds",
	"resourcegroups",
	"route53",
	"route53resolver",
	"sagemaker",
	"securityhub",
	"sfn",
	"sns",
	"sqs",
	"ssm",
	"storagegateway",
	"swf",
	"transfer",
	"waf",
	"wafregional",
	"wafv2",
	"workspaces",
}

type TemplateData struct {
	ServiceNames []string
}

func main() {
	// Always sort to reduce any potential generation churn
	sort.Strings(serviceNames)

	templateData := TemplateData{
		ServiceNames: serviceNames,
	}
	templateFuncMap := template.FuncMap{
		"ClientType":                           keyvaluetags.ServiceClientType,
		"ListTagsFunction":                     ServiceListTagsFunction,
		"ListTagsInputIdentifierField":         ServiceListTagsInputIdentifierField,
		"ListTagsInputIdentifierRequiresSlice": ServiceListTagsInputIdentifierRequiresSlice,
		"ListTagsInputResourceTypeField":       ServiceListTagsInputResourceTypeField,
		"ListTagsOutputTagsField":              ServiceListTagsOutputTagsField,
		"TagPackage":                           keyvaluetags.ServiceTagPackage,
		"Title":                                strings.Title,
	}

	tmpl, err := template.New("listtags").Funcs(templateFuncMap).Parse(templateBody)

	if err != nil {
		log.Fatalf("error parsing template: %s", err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, templateData)

	if err != nil {
		log.Fatalf("error executing template: %s", err)
	}

	generatedFileContents, err := format.Source(buffer.Bytes())

	if err != nil {
		log.Fatalf("error formatting generated file: %s", err)
	}

	f, err := os.Create(filename)

	if err != nil {
		log.Fatalf("error creating file (%s): %s", filename, err)
	}

	defer f.Close()

	_, err = f.Write(generatedFileContents)

	if err != nil {
		log.Fatalf("error writing to file (%s): %s", filename, err)
	}
}

var templateBody = `
// Code generated by generators/listtags/main.go; DO NOT EDIT.

package keyvaluetags

import (
	"github.com/aws/aws-sdk-go/aws"
{{- range .ServiceNames }}
	"github.com/aws/aws-sdk-go/service/{{ . }}"
{{- end }}
)
{{ range .ServiceNames }}

// {{ . | Title }}ListTags lists {{ . }} service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func {{ . | Title }}ListTags(conn {{ . | ClientType }}, identifier string{{ if . | ListTagsInputResourceTypeField }}, resourceType string{{ end }}) (KeyValueTags, error) {
	input := &{{ . | TagPackage  }}.{{ . | ListTagsFunction }}Input{
		{{- if . | ListTagsInputIdentifierRequiresSlice }}
		{{ . | ListTagsInputIdentifierField }}:   aws.StringSlice([]string{identifier}),
		{{- else }}
		{{ . | ListTagsInputIdentifierField }}:   aws.String(identifier),
		{{- end }}
		{{- if . | ListTagsInputResourceTypeField }}
		{{ . | ListTagsInputResourceTypeField }}: aws.String(resourceType),
		{{- end }}
	}

	output, err := conn.{{ . | ListTagsFunction }}(input)

	if err != nil {
		return New(nil), err
	}

	return {{ . | Title }}KeyValueTags(output.{{ . | ListTagsOutputTagsField }}), nil
}
{{- end }}
`

// ServiceListTagsFunction determines the service tagging function.
func ServiceListTagsFunction(serviceName string) string {
	switch serviceName {
	case "acm":
		return "ListTagsForCertificate"
	case "acmpca":
		return "ListTags"
	case "backup":
		return "ListTags"
	case "cloudhsmv2":
		return "ListTags"
	case "cloudtrail":
		return "ListTags"
	case "cloudwatchlogs":
		return "ListTagsLogGroup"
	case "dax":
		return "ListTags"
	case "directconnect":
		return "DescribeTags"
	case "dynamodb":
		return "ListTagsOfResource"
	case "efs":
		return "DescribeTags"
	case "elasticsearchservice":
		return "ListTags"
	case "elbv2":
		return "DescribeTags"
	case "firehose":
		return "ListTagsForDeliveryStream"
	case "glue":
		return "GetTags"
	case "kinesis":
		return "ListTagsForStream"
	case "kms":
		return "ListResourceTags"
	case "lambda":
		return "ListTags"
	case "mq":
		return "ListTags"
	case "opsworks":
		return "ListTags"
	case "redshift":
		return "DescribeTags"
	case "resourcegroups":
		return "GetTags"
	case "sagemaker":
		return "ListTags"
	case "sqs":
		return "ListQueueTags"
	case "workspaces":
		return "DescribeTags"
	default:
		return "ListTagsForResource"
	}
}

// ServiceListTagsInputIdentifierField determines the service tag identifier field.
func ServiceListTagsInputIdentifierField(serviceName string) string {
	switch serviceName {
	case "acm":
		return "CertificateArn"
	case "acmpca":
		return "CertificateAuthorityArn"
	case "athena":
		return "ResourceARN"
	case "cloudhsmv2":
		return "ResourceId"
	case "cloudtrail":
		return "ResourceIdList"
	case "cloudwatch":
		return "ResourceARN"
	case "cloudwatchevents":
		return "ResourceARN"
	case "cloudwatchlogs":
		return "LogGroupName"
	case "dax":
		return "ResourceName"
	case "devicefarm":
		return "ResourceARN"
	case "directconnect":
		return "ResourceArns"
	case "directoryservice":
		return "ResourceId"
	case "docdb":
		return "ResourceName"
	case "efs":
		return "FileSystemId"
	case "elasticache":
		return "ResourceName"
	case "elasticsearchservice":
		return "ARN"
	case "elbv2":
		return "ResourceArns"
	case "firehose":
		return "DeliveryStreamName"
	case "fsx":
		return "ResourceARN"
	case "kinesis":
		return "StreamName"
	case "kinesisanalytics":
		return "ResourceARN"
	case "kinesisanalyticsv2":
		return "ResourceARN"
	case "kms":
		return "KeyId"
	case "lambda":
		return "Resource"
	case "mediaconvert":
		return "Arn"
	case "mediastore":
		return "Resource"
	case "neptune":
		return "ResourceName"
	case "organizations":
		return "ResourceId"
	case "rds":
		return "ResourceName"
	case "redshift":
		return "ResourceName"
	case "resourcegroups":
		return "Arn"
	case "route53":
		return "ResourceId"
	case "sqs":
		return "QueueUrl"
	case "ssm":
		return "ResourceId"
	case "storagegateway":
		return "ResourceARN"
	case "transfer":
		return "Arn"
	case "workspaces":
		return "ResourceId"
	case "waf":
		return "ResourceARN"
	case "wafregional":
		return "ResourceARN"
	case "wafv2":
		return "ResourceARN"
	default:
		return "ResourceArn"
	}
}

// ServiceTagInputIdentifierRequiresSlice determines if the service tagging resource field requires a slice.
func ServiceListTagsInputIdentifierRequiresSlice(serviceName string) string {
	switch serviceName {
	case "cloudtrail":
		return "yes"
	case "directconnect":
		return "yes"
	case "elbv2":
		return "yes"
	default:
		return ""
	}
}

// ServiceListTagsInputResourceTypeField determines the service tagging resource type field.
func ServiceListTagsInputResourceTypeField(serviceName string) string {
	switch serviceName {
	case "route53":
		return "ResourceType"
	case "ssm":
		return "ResourceType"
	default:
		return ""
	}
}

// ServiceListTagsOutputTagsField determines the service tag field.
func ServiceListTagsOutputTagsField(serviceName string) string {
	switch serviceName {
	case "cloudhsmv2":
		return "TagList"
	case "cloudtrail":
		return "ResourceTagList[0].TagsList"
	case "databasemigrationservice":
		return "TagList"
	case "directconnect":
		return "ResourceTags[0].Tags"
	case "docdb":
		return "TagList"
	case "elasticache":
		return "TagList"
	case "elasticbeanstalk":
		return "ResourceTags"
	case "elasticsearchservice":
		return "TagList"
	case "elbv2":
		return "TagDescriptions[0].Tags"
	case "mediaconvert":
		return "ResourceTags.Tags"
	case "neptune":
		return "TagList"
	case "rds":
		return "TagList"
	case "route53":
		return "ResourceTagSet.Tags"
	case "ssm":
		return "TagList"
	case "waf":
		return "TagInfoForResource.TagList"
	case "wafregional":
		return "TagInfoForResource.TagList"
	case "wafv2":
		return "TagInfoForResource.TagList"
	case "workspaces":
		return "TagList"
	default:
		return "Tags"
	}
}