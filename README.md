SSM to read client and email in dev/stage/prod
Verify email identity in SES

Need variables:
FROM : from email id "root@imfo.se"
SSM_PATH : containing json 
`{
"org1": {
"to": "abc@gmail.com, def@gmail.com",
"cc": "z@mail.com, y@mail.com, x@mail.com",
"bcc": "o@hmail.com, p@hmail.com"
},
"org2": {
"to": "xyz@gmail.com.com",
"cc": "",
"bcc": ""
},
"org3":{
"to": "mno@gmail.com.com",
"cc": "",
"bcc": ""
}
}`

MONGO_URI : containing mongoDB string uri
MONGO_DATABASE : database name ("demo")
MONGO_COLLECTION : collection name ("tokenLimit")
SUBJECT : email subject
BUCKET : where html template is placed on s3
KEY : directory and name of file in S3 bucket