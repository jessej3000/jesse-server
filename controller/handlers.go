package controller

import (
	"context"
	"errors"
	"fmt"
	"server/router"
	"time"

	"github.com/google/uuid"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type category struct {
	ID   int
	Name string
}

type verifyConn struct {
	Status        int
	Paid          int
	ShowMessaging bool
	Categories    []category
	Msg           addMsg
	Session       string
}

type addMsg struct {
	MsgID     string
	MsgFromID string
	MsgToID   string
	MsgBody   string
	MsgTime   time.Time
	MsgAuthor string
}

func MessageIn(Client *router.Client, data interface{}, manager *router.Manager) {
	var msgReceived map[string]interface{}
	// Add message

	msgReceived = data.(map[string]interface{})

	msgTime := time.Now()

	detail := manager.Clients[Client]
	// msgfrom := detail.ID
	msgto := "242" // msgReceived["recipientid"].(string)

	resp, error := DetectIntentText(
		"virtualself-gseo",
		msgReceived["session"].(string),
		msgReceived["msg"].(string),
		"us-en",
	)

	if error != nil {
		return
	}

	msgUid := uuid.New()
	msg := router.Message{
		Name: "msg_in",
		Data: addMsg{
			MsgID:     msgUid.String(),
			MsgFromID: msgto,
			MsgToID:   detail.ID,
			MsgBody:   resp,
			MsgTime:   msgTime,
			MsgAuthor: "Jesse",
		},
	}
	Client.Send <- msg

	fmt.Println("::::::::::::::::>>>>>>>>>>>" + msgReceived["msg"].(string) + " SESSION : " + msgReceived["session"].(string))
	// for cli, _ := range manager.Clients {
	// 	// if (msgto == det.ID) || (msgfrom == det.ID) {
	// 	msg := router.Message{Name: "msg_in", Data: addMsg{MsgID: "2", MsgFromID: detail.ID, MsgToID: msgto, MsgBody: msgReceived["msg"].(string), MsgTime: msgTime, MsgAuthor: detail.Name}}
	// 	cli.Send <- msg
	// 	// }
	// }

	// resp, _ := DetectIntentText(
	// 	"irtualself-gseo",
	// 	"12345678",
	// 	msgReceived["msg"].(string),
	// 	"us-en",
	// )
	// fmt.Printf(resp)
}

func DetectIntentText(projectID, sessionID, text, languageCode string) (string, error) {
	ctx := context.Background()

	sessionClient, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		return "", err
	}
	defer sessionClient.Close()

	if projectID == "" || sessionID == "" {
		return "", errors.New(fmt.Sprintf("Received empty project (%s) or session (%s)", projectID, sessionID))
	}

	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)
	textInput := dialogflowpb.TextInput{Text: text, LanguageCode: languageCode}
	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}
	request := dialogflowpb.DetectIntentRequest{Session: sessionPath, QueryInput: &queryInput}

	response, err := sessionClient.DetectIntent(ctx, &request)
	if err != nil {
		return "", err
	}

	queryResult := response.GetQueryResult()
	fulfillmentText := queryResult.GetFulfillmentText()
	return fulfillmentText, nil
}

func Identify(Client *router.Client, data interface{}, chatManager *router.Manager) {
	var dataReceived map[string]interface{}
	var categories []category

	dataReceived = data.(map[string]interface{})
	chatManager.Clients[Client] = router.ClientDetail{
		ID:      dataReceived["ID"].(string),
		Name:    dataReceived["Name"].(string),
		IP:      chatManager.Clients[Client].IP,
		Session: chatManager.Clients[Client].Session,
	} // 0 Means client is connected

	resp, error := DetectIntentText(
		"virtualself-gseo",
		chatManager.Clients[Client].Session,
		"hi",
		"us-en",
	)

	if error != nil {
		return
	}
	msgTime := time.Now()

	msgUid := uuid.New()
	dataMsg := addMsg{
		MsgID:     msgUid.String(),
		MsgFromID: "242",
		MsgToID:   chatManager.Clients[Client].ID,
		MsgBody:   resp,
		MsgTime:   msgTime,
		MsgAuthor: "Jesse",
	}

	msg := router.Message{
		Name: "identified",
		Data: verifyConn{
			Status:        1,
			Paid:          0,
			ShowMessaging: false,
			Categories:    categories,
			Msg:           dataMsg,
			Session:       chatManager.Clients[Client].Session,
		},
	}
	fmt.Println(chatManager.Clients[Client].Name + " Identified... With ID: " + chatManager.Clients[Client].ID + " & Session ID: " + chatManager.Clients[Client].Session)

	Client.Send <- msg
}
