# cdktf-golang-aws-iam

The Cloud Development Kit for Terraform (CDKTF) allows you to define your infrastructure in a familiar programming language such as TypeScript, Python, Go, C#, or Java.

In this tutorial, you will provision an EC2 instance on AWS using your preferred programming language.

## Prerequisites

* [Terraform](https://www.terraform.io/downloads) >= v1.0
* [CDK for Terraform](https://learn.hashicorp.com/tutorials/terraform/cdktf-install) >= v0.8
* A [Terraform Cloud](https://app.terraform.io/) account, with [CLI authentication](https://learn.hashicorp.com/tutorials/terraform/cloud-login) configured
* [an AWS account](https://portal.aws.amazon.com/billing/signup?nc2=h_ct&src=default&redirect_url=https%3A%2F%2Faws.amazon.com%2Fregistration-confirmation#/start)
* AWS Credentials [configured for use with Terraform](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#authentication)


Credentials can be provided by using the AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and optionally AWS_SESSION_TOKEN environment variables. The region can be set using the AWS_REGION or AWS_DEFAULT_REGION environment variables.

```shell
$ export AWS_ACCESS_KEY_ID="anaccesskey"
$ export AWS_SECRET_ACCESS_KEY="asecretkey"
$ export AWS_REGION="us-west-2"
```

## Install project dependencies

```shell
mkdir cdktf-go-aws-iam
cd cdktf-go-aws-iam
cdktf init --template="go"
```

## Install AWS provider
Open `cdktf.json` in your text editor, and add `aws` as one of the Terraform providers that you will use in the application.
```JSON
{
  "language": "go",
  "app": "go run main.go",
  "codeMakerOutput": "generated",
  "projectId": "02f2d864-a2f2-49e8-ab52-b472e233755e",
  "sendCrashReports": "false",
  "terraformProviders": [
	 "hashicorp/aws@~> 3.67.0"
  ],
  "terraformModules": [],
  "context": {
    "excludeStackIdFromLogicalIds": "true",
    "allowSepCharsInLogicalIds": "true"
  }
}
```
Run `cdktf get` to install the AWS provider you added to `cdktf.json`.
```SHELL
cdktf get
```

CDKTF uses a library called `jsii` to allow Go code to interact with CDK, 
which is written in TypeScript. 
Ensure that the jsii runtime is installed by running `go mod tidy`.

```SHELL
go mod tidy
```

## Define your CDK for Terraform Application

Replace the contents of main.py with the following code for a new Golang application

```golang
package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"

	"cdk.tf/go/stack/generated/hashicorp/aws"
	"cdk.tf/go/stack/generated/hashicorp/aws/iam"
)

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)
	aws.NewAwsProvider(stack, jsii.String("AWS"), &aws.AwsProviderConfig{
		Region: jsii.String("us-west-1"),
	})

	group := iam.NewIamGroup(stack, jsii.String("iam-group-demo"), &iam.IamGroupConfig{
		Name: jsii.String("CDKtf-Golang-Group-Demo"),
	})

	user := iam.NewIamUser(stack, jsii.String("iam-user-demo"), &iam.IamUserConfig{
		Name: jsii.String("CDKtf-Golang-User-Demo"),
		Tags: &map[string]*string{
			"Name":    jsii.String("CDKtf-Golang-User-Demo"),
			"Team":    jsii.String("Devops"),
			"Company": jsii.String("Your compnay"),
		},
	})

	role := iam.NewIamRole(stack, jsii.String("iam-role-demo"), &iam.IamRoleConfig{
		Name: jsii.String("CDKtf-Golang-role-Demo"),
		AssumeRolePolicy: jsii.String(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {
					"Service": "ec2.amazonaws.com"
				},
				"Action": "sts:AssumeRole"
			}]
		}`),

		Tags: &map[string]*string{
			"Name":    jsii.String("CDKtf-Golang-role-Demo"),
			"Team":    jsii.String("Devops"),
			"Company": jsii.String("Your compnay"),
		},
	})

	policy := iam.NewIamPolicy(stack, jsii.String("iam-policy-demo"), &iam.IamPolicyConfig{
		Name: jsii.String("CDKtf-Golang-policy-Demo"),
		Policy: jsii.String(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Action": "*",
				"Resource": ["arn:aws:ec2:*:*:client-vpn-endpoint/*"],
				"Effect": "Allow"
			}]
		}`),
		Description: jsii.String("This policy is for Golang demo"),
	})

	iam.NewIamGroupMembership(stack, jsii.String("iam-group-membership-demo"), &iam.IamGroupMembershipConfig{
		Name:  jsii.String("group-membership"),
		Group: group.Name(),
		Users: jsii.Strings(*user.Name()),
	})

	attachment := iam.NewIamPolicyAttachment(stack, jsii.String("iam-application-managed-policy-demo"), &iam.IamPolicyAttachmentConfig{
		Name:      jsii.String("CDKtf-Golang-iam-attachment-Demo"),
		Groups:    jsii.Strings(*group.Name()),
		Roles:     jsii.Strings(*role.Name()),
		Users:     jsii.Strings(*user.Name()),
		PolicyArn: jsii.String(*policy.Arn()),
	})

	cdktf.NewTerraformOutput(stack, jsii.String("iam-group"), &cdktf.TerraformOutputConfig{
		Value: group.Name(),
	})

	cdktf.NewTerraformOutput(stack, jsii.String("iam-user"), &cdktf.TerraformOutputConfig{
		Value: group.Name(),
	})

	cdktf.NewTerraformOutput(stack, jsii.String("iam-role"), &cdktf.TerraformOutputConfig{
		Value: role.Arn(),
	})

	cdktf.NewTerraformOutput(stack, jsii.String("iam-policy"), &cdktf.TerraformOutputConfig{
		Value: policy.Arn(),
	})

	cdktf.NewTerraformOutput(stack, jsii.String("iam-attachment"), &cdktf.TerraformOutputConfig{
		Value: attachment.Name(),
	})

	return stack
}

func main() {
	app := cdktf.NewApp(nil)
	stack := NewMyStack(app, "aws_instance")
	cdktf.NewRemoteBackend(stack, &cdktf.RemoteBackendProps{
		Hostname:     jsii.String("app.terraform.io"),
		Organization: jsii.String("jigsaw373"),
		Workspaces:   cdktf.NewNamedRemoteWorkspace(jsii.String("cdktf-go-aws-iam")),
	})

	app.Synth()
}

```
## Provision infrastructure
```shell
cdktf deploy
```
After the instance is created, visit the AWS EC2 Dashboard.

## Clean up your infrastructure
```shell
cdktf destroy
```
