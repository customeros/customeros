const http = require('http');
var AWS = require('aws-sdk');
var S3 = new AWS.S3({ signatureVersion: 'v4' });

exports.handler = function(event, context, callback) {
    console.log('Spam filter');

    var sesNotification = event.Records[0].ses;
    console.log("SES Notification:\n", JSON.stringify(event.Records[0], null, 2));

    // Check if any spam check failed
    if (sesNotification.receipt.spfVerdict.status === 'FAIL' ||
        sesNotification.receipt.dkimVerdict.status === 'FAIL' ||
        sesNotification.receipt.spamVerdict.status === 'FAIL' ||
        sesNotification.receipt.virusVerdict.status === 'FAIL') {
        console.log('Dropping spam');
        // Stop processing rule set, dropping message
        callback(null, { 'disposition': 'STOP_RULE_SET' });
    }
    else {
        const bucket = process.env.OL_MAIL_S3_BUCKET;
        const key = sesNotification.mail.messageId;

        var params = { Bucket: bucket, Key: key };
        console.log('Checking file existence');
        console.log(params);
        console.log('Calling s3.getObject');
        S3.getObject(params, function(err, data) {
            console.log('S3.getObject called');
            console.log('err = ' + err);
            if (err) {
                console.log(err, err.stack); // an error occurred
                callback(err);
            }
            else {
                console.log(data); // successful response
                var body = data.Body.toString();
                console.log(body); // successful response
                const mailData = JSON.stringify({
                    sender: sesNotification.mail.commonHeaders['from'][0],
                    rawMessage: body,
                    subject: sesNotification.mail.commonHeaders['subject'],
                    'api-key': process.env.OL_API_KEY,
                    'X-Openline-TENANT': process.env.OL_TENANT_NAME
                });
                console.log(mailData);
                const options = {
                    hostname: process.env.OL_MAIL_CB_HOST,
                    port: process.env.OL_MAIL_CB_PORT,
                    path: process.env.OL_MAIL_CB_PATH,
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Content-Length': mailData.length,
                    },
                };



                const req = http.request(options, res => {
                    console.log(`statusCode: ${res.statusCode}`);

                    res.on('data', d => {
                        process.stdout.write(d);
                    });
                });

                req.on('error', error => {
                    console.error(error);
                });

                req.write(mailData);
                req.end();
                S3.deleteObject(params, function(err, data) {
                    console.log('S3.getObject called');
                    console.log('err = ' + err);
                    if (err) {
                        console.log(err, err.stack); // an error occurred
                        callback(err);
                    }
                    else {
                        console.log('Leaving s3.deleteObject');
                        callback(null, null);
                    }
                });

            }
            console.log('Leaving s3.getObject');
        });
    }
};