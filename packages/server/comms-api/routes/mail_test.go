package routes

import (
	"github.com/gin-gonic/gin"
)

var mailRouter *gin.Engine

const mailApiKey = "f76f637a-0863-4509-8956-591fdec6a73e"

func init() {
	//mailRouter = gin.Default()
	//route := mailRouter.Group("/")
	//
	//conf := &config.Config{}
	//dft := test_utils.MakeDialFactoryTest(conf)
	//conf.Mail.ApiKey = mailApiKey
	//addMailRoutes(conf, dft, route)
}

const rawEmail = "To: agent@agent.secretcorp.com\r\n" +
	"From: gabi@example.org\r\n" +
	"Subject: Help Please\r\n" +
	"\r\n" +
	"Hello Gabi\r\n"

	//
	//func TestGetMail(t *testing.T) {
	//	messageId := int64(7)
	//	contactId := "2"
	//	sentMessageEvent := false
	//
	//	mpr := &MailPostRequest{
	//		ApiKey:     mailApiKey,
	//		Sender:     "john.doe@example.org",
	//		Subject:    "Help Please",
	//		RawMessage: rawEmail,
	//	}
	//	test_utils.SetMessageStoreCallbacks(&test_utils.MockMessageServiceCallbacks{SaveMessage: func(ctx context.Context, message *msProto.MessageDeprecate) (*msProto.MessageDeprecate, error) {
	//		log.Printf("Inside SaveMessage")
	//		var tm *time.Time = nil
	//		if message.GetTime() != nil {
	//			var timeref = message.GetTime().AsTime()
	//			tm = &timeref
	//		}
	//
	//		if tm == nil {
	//			var timeRef = time.Now()
	//			tm = &timeRef
	//		}
	//
	//		message.Time = timestamppb.New(*tm)
	//		message.Id = &messageId
	//
	//		if message.Username == nil {
	//			if !assert.Equal(t, "john.doe@example.org", message.Username) {
	//				return nil, status.Error(400, "Unexpected username")
	//			}
	//		} else {
	//			if !assert.Equal(t, contactId, *message.ContactId) {
	//				return nil, status.Error(400, "Unexpected contact ID")
	//			}
	//		}
	//
	//		message = &msProto.MessageDeprecate{
	//			Id:        &messageId,
	//			ContactId: &contactId,
	//		}
	//		return message, nil
	//	}})
	//
	//	test_utils.SetChannelApiCallbacks(&test_utils.MockChannelApiCallbacks{NewMessageEvent: func(ctx context.Context, id *proto.NewMessage) (*proto.OasisEmpty, error) {
	//		if !assert.Equal(t, id, messageId) {
	//			return nil, status.Error(400, "Unexpected message id!")
	//		}
	//		sentMessageEvent = true
	//		return &proto.OasisEmpty{}, nil
	//	}})
	//
	//	w := httptest.NewRecorder()
	//	reqBody, _ := json.Marshal(mpr)
	//	req, _ := http.NewRequest("POST", "/mail/", bytes.NewReader(reqBody))
	//	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	//	log.Printf("Going to post body %s", reqBody)
	//	mailRouter.ServeHTTP(w, req)
	//	log.Printf("Got Body %s", w.Body)
	//	if !assert.Equal(t, w.Code, 200) {
	//		return
	//	}
	//	assert.True(t, sentMessageEvent, "NewMessageEvent not called!")
	//}

	//func TestGetMailInvalidKey(t *testing.T) {
	//	id1 := int64(7)
	//	gabyId := int64(2)
	//	sentMessageEvent := false
	//	sentSaveMessage := false
	//
	//	mpr := &MailPostRequest{
	//		ApiKey:     "Invalid Key",
	//		Sender:     "gabi@example.org",
	//		Subject:    "Help Please",
	//		RawMessage: rawEmail,
	//	}
	//	test_utils.SetMessageStoreCallbacks(&test_utils.MockMessageServiceCallbacks{SaveMessage: func(ctx context.Context, message *msProto.MessageDeprecate) (*msProto.MessageDeprecate, error) {
	//		log.Printf("Inside SaveMessage")
	//		sentSaveMessage = true
	//
	//		var tm *time.Time = nil
	//		if message.GetTime() != nil {
	//			var timeref = message.GetTime().AsTime()
	//			tm = &timeref
	//		}
	//
	//		if tm == nil {
	//			var timeRef = time.Now()
	//			tm = &timeRef
	//		}
	//
	//		message.Time = timestamppb.New(*tm)
	//		message.Id = &id1
	//
	//		if message.Contact == nil {
	//			if !assert.Equal(t, "gabi@example.org", message.Username) {
	//				return nil, status.Error(400, "Unexpected username")
	//			}
	//		} else {
	//			if !assert.Equal(t, gabyId, *message.Contact.Id) {
	//				return nil, status.Error(400, "Unexpected contact ID")
	//			}
	//		}
	//
	//		message.Contact = &msProto.Contact{
	//			ContactId: "77775553",
	//			FirstName: "Gabriel",
	//			LastName:  "Gontariu",
	//			Id:        &gabyId,
	//		}
	//		return message, nil
	//	}})
	//
	//	test_utils.SetChannelApiCallbacks(&test_utils.MockChannelApiCallbacks{NewMessageEvent: func(ctx context.Context, id *proto.OasisMessageId) (*proto.OasisEmpty, error) {
	//		if !assert.Equal(t, id.MessageId, id1) {
	//			return nil, status.Error(400, "Unexpected message id!")
	//		}
	//		sentMessageEvent = true
	//		return &proto.OasisEmpty{}, nil
	//	}})
	//
	//	w := httptest.NewRecorder()
	//	reqBody, _ := json.Marshal(mpr)
	//	req, _ := http.NewRequest("POST", "/mail/", bytes.NewReader(reqBody))
	//	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	//	log.Printf("Going to post body %s", reqBody)
	//	mailRouter.ServeHTTP(w, req)
	//	log.Printf("Got Body %s", w.Body)
	//	assert.Equal(t, w.Code, 403)
	//	assert.False(t, sentMessageEvent, "NewMessageEvent  called when it shouldn't!")
	//	assert.False(t, sentSaveMessage, "SaveMessage Called when it shouldn't!")
	//
	//}
