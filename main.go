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
