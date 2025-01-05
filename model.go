package mail_sending

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type MailHandler struct {
	S3Instance      *s3.S3
	SessionInstance *session.Session
	SSMInstance     *ssm.SSM
}

type TokenLimit struct {
	AssignedTokenLimit      int    `bson:"assignedTokenLimit"`
	Org                     string `bson:"org"`
	Provider                string `bson:"provider"`
	RemainingTokenLimit     int    `bson:"remainingTokenLimit"`
	TokenConsumptionTrigger int    `bson:"tokenConsumptionTrigger,omitempty"`
	UsedTokenLimit          int    `bson:"usedTokenLimit"`
}

type HtmlDetails struct {
	Org  string
	Data []TokenLimit
}

type EMailDetails struct {
	To  string `json:"to"`
	Cc  string `json:"cc"`
	Bcc string `json:"bcc"`
}
