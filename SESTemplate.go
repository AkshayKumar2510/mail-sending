package mail_sending

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

func (mh *MailHandler) GetHTMLTemplate(ctx context.Context, org string, details []TokenLimit) string {
	var templateBuffer bytes.Buffer
	data := HtmlDetails{
		Org:  org,
		Data: details,
	}
	//key := os.Getenv("KEY")
	//bucket := os.Getenv("BUCKET")
	htmlData, err := os.ReadFile("./sampleFiles/email-notifier-template.html")
	if err != nil {
		log.Println("Failed to read template file.", " Error:", err)
		return ""
	}
	//templateFromS3, err := mh.S3Instance.GetObjectWithContext(ctx, &s3.GetObjectInput{
	//	Bucket: aws.String(bucket),
	//	Key:    aws.String(key),
	//})
	//if err != nil {
	//	log.Fatalf("unable to fetch template name:%s from bucket:%s, error:%v", key, bucket, err)
	//}
	//htmlData, err := io.ReadAll(templateFromS3.Body)
	//if err != nil {
	//	log.Fatalf("unable to read template from bucket, %v", err)
	//}
	htmlTemplate := template.Must(template.New("template.html").Parse(string(htmlData)))

	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "template.html", data)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return templateBuffer.String()
}

func (mailHandle *MailHandler) GenerateSESTemplate(ctx context.Context, org string, details []TokenLimit) (template *ses.SendEmailInput, err error) {
	from := os.Getenv("FROM")
	// pathParam := os.Getenv("SSM_PATH")
	// res, err := mailHandle.SSMInstance.GetParameterWithContext(ctx, &ssm.GetParameterInput{
	// 	Name:           &pathParam,
	// 	WithDecryption: aws.Bool(true),
	// })
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to fetch email from ssm for org: %s, err:%w", org, err)
	// }

	file, err := os.ReadFile(os.Getenv("CONFIG_FILE"))
	if err != nil {
		log.Fatalf("couldn't read config file: %s. Error: %v", os.Getenv("CONFIG_FILE"), err)
	}

	var email map[string]EMailDetails
	err = json.Unmarshal(file, &email)
	if err != nil {
		log.Fatalf("couldn't marshal config data %s. Error:%v", string(file), err)
	}

	// var email map[string]EMailDetails
	// value := []byte(*res.Parameter.Value)
	// if err = json.Unmarshal(value, &email); err != nil {
	// 	return nil, fmt.Errorf("failed to parse stored email details for org: %s got value:%s, err:%w", org, string(value), err)
	// }

	toDetails := convertToStringPointers(strings.Split(email[org].To, ","))
	ccDetails := convertToStringPointers(strings.Split(email[org].Cc, ","))
	bccDetails := convertToStringPointers(strings.Split(email[org].Bcc, ","))

	html := mailHandle.GetHTMLTemplate(ctx, org, details)

	title := os.Getenv("SUBJECT")

	template = &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses:  ccDetails,
			BccAddresses: bccDetails,
			ToAddresses:  toDetails,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("utf-8"),
					Data:    aws.String(html),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    aws.String(title),
			},
		},
		Source: aws.String(from),
	}
	return
}
func convertToStringPointers(slice []string) []*string {
	ptrs := make([]*string, len(slice))
	for i, s := range slice {
		ptrs[i] = aws.String(s)
	}
	return ptrs
}
