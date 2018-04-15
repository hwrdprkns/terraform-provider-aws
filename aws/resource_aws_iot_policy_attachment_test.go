package aws

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAWSIoTPolicyAttachment_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAWSIoTPolicyAttachmentConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("aws_iot_policy_attachment.foo_attachment", "policy_name"),
					resource.TestCheckResourceAttrSet("aws_iot_policy_attachment.foo_attachment", "target"),
				),
			},
		},
	})
}

const testAccAWSIoTPolicyAttachmentConfig_basic = `
resource "aws_iot_certificate" "foo_cert" {
  csr = "${file("test-fixtures/iot-csr.pem")}"
  active = true
}

resource "aws_iot_policy" "foo_policy" {
  name = "PubSubToAnyTopic"
  policy = <<EOF
{
	"Version": "2012-10-17",
	"Statement": [{
		"Effect": "Allow",
		"Action": ["iot:*"],
		"Resource": ["*"]
	}]
}
EOF
}

resource "aws_iot_policy_attachment" "foo_attachment" {
	policy_name = "${aws_iot_policy.foo_policy.name}"
	target = "${aws_iot_certificate.foo_cert.arn}"
}
`
