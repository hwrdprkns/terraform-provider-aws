package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsIotPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsIotPolicyAttachmentCreate,
		Read:   resourceAwsIotPolicyAttachmentRead,
		Delete: resourceAwsIotPolicyAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"policy_name": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"target": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAwsIotPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).iotconn
	name := d.Get("policy_name").(string)
	target := d.Get("target").(string)

	_, err := conn.AttachPolicy(&iot.AttachPolicyInput{
		PolicyName: aws.String(name),
		Target:     aws.String(target),
	})

	if err != nil {
		return fmt.Errorf("Error attaching IoT Policy to Target: %s", err)
	}

	d.SetId(fmt.Sprintf("%s-%s", name, target))
	return resourceAwsIotPolicyAttachmentRead(d, meta)
}

func resourceAwsIotPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).iotconn
	name := d.Get("policy_name").(string)
	target := d.Get("target").(string)

	resp, err := conn.ListTargetsForPolicy(&iot.ListTargetsForPolicyInput{
		PolicyName: aws.String(name),
	})
	if err != nil {
		if isAWSErr(err, iot.ErrCodeResourceNotFoundException, "") {
			log.Printf("[WARN] IoT Policy Attachment %q not found, removing from state", name)
			d.SetId("")
		}
		return err
	}

	_, found := sliceContainsString(flattenStringList(resp.Targets), target)
	if !found {
		d.SetId("")
	}

	return nil
}

func resourceAwsIotPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).iotconn

	params := &iot.DetachPolicyInput{
		PolicyName: aws.String(d.Get("policy_name").(string)),
		Target:     aws.String(d.Get("target").(string)),
	}
	log.Printf("[DEBUG] Detaching IoT policy %s", params)

	_, err := conn.DetachPolicy(params)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return err
	}

	return nil
}
