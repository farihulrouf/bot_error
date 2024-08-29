package controllers

import (
	"fmt"
	"strings"
	"encoding/json"
	"wagobot.com/db"
	// "encoding/base64"
	// "wagobot.com/base"
	"wagobot.com/model"
	// "wagobot.com/response"
	// "go.mau.fi/whatsmeow"
	// "go.mau.fi/whatsmeow/store"
	// "go.mau.fi/whatsmeow/store/sqlstore"
	// "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func EventHandler(evt interface{}) {

	switch v := evt.(type) {
	case *events.Message:

		fmt.Println("------ new message ")
		// fmt.Println(evt)

		jsonResponse, err := json.MarshalIndent(v, "", "")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(jsonResponse))
		}

		chatId := v.Info.Chat.String()
		theType := "post"
		replyToPost := ""
		replyToUser := ""
		mediaType := v.Info.Type
		media := [] model.Media{}

		if !v.Info.IsGroup {
			chatId = model.GetPhoneNumber(chatId)
		}
		
		txtMessage := ""
		if v.Message.ExtendedTextMessage != nil {
			ext := v.Message.GetExtendedTextMessage()
			ci := ext.GetContextInfo()
			if ci != nil {
				txtMessage = ext.GetText()
				theType = "reply"
				replyToPost = ci.GetStanzaID()
				replyToUser = model.GetPhoneNumber(ci.GetParticipant())
			}
		} else {
			txtMessage = v.Message.GetConversation()
		}

		if v.Message.ImageMessage != nil {
			img := v.Message.GetImageMessage()
			imgCaption := img.GetCaption()

			// url := img.GetURL()

			media = append(media, model.Media {
				Type: mediaType,
				Caption: imgCaption,
				MimeType: img.GetMimetype(),
				Thumbnail: img.GetJPEGThumbnail(),
				FileLength: img.GetFileLength(),
			})
			if txtMessage != "" {
				txtMessage += ". "
			}
			txtMessage += imgCaption
		}

		if v.Message.DocumentMessage != nil {
			doc := v.Message.GetDocumentMessage()
			docCaption := doc.GetCaption()

			// url := doc.GetURL()

			media = append(media, model.Media {
				Type: mediaType,
				Caption: docCaption,
				FileName: doc.GetFileName(),
				FileLength: doc.GetFileLength(),
				MimeType: doc.GetMimetype(),
			})
			if txtMessage != "" {
				txtMessage += ". "
			}
			txtMessage += docCaption
		}

		if v.Message.AudioMessage != nil {
			aud := v.Message.GetAudioMessage()
			audCaption := "" //aud.GetCaption()

			// url := aud.GetURL()

			media = append(media, model.Media {
				Type: mediaType,
				Caption: audCaption,
				MimeType: aud.GetMimetype(),
				Seconds: aud.GetSeconds(),
				FileLength: aud.GetFileLength(),
			})
			if txtMessage != "" {
				txtMessage += ". "
			}
			txtMessage += audCaption
		}

		if v.Message.VideoMessage != nil {
			vid := v.Message.GetVideoMessage()
			vidCaption := vid.GetCaption()

			// url := vid.GetURL()

			media = append(media, model.Media {
				Type: mediaType,
				Caption: vidCaption,
				MimeType: vid.GetMimetype(),
				Thumbnail: vid.GetJPEGThumbnail(),
				FileLength: vid.GetFileLength(),
			})
			if txtMessage != "" {
				txtMessage += ". "
			}
			txtMessage += vidCaption
		}

		if v.Message.ContactMessage != nil {
			ctc := v.Message.GetContactMessage()
			media = append(media, model.Media {
				Name: ctc.GetDisplayName(),
				Contact: ctc.GetVcard(),
			})
			if txtMessage != "" {
				txtMessage += ". "
			}
			txtMessage += ctc.GetVcard()
		}

		if v.Message.PollCreationMessageV3 != nil {
			pol := v.Message.GetPollCreationMessageV3()
			polName := pol.GetName()
			jsonData, _ := json.Marshal(pol)
			media = append(media, model.Media {
				Poll: string(jsonData),
			})
			if txtMessage != "" {
				txtMessage += ". "
			}
			txtMessage += polName
		}

		if v.Message.LocationMessage != nil {
			loc := v.Message.GetLocationMessage()
			lat := loc.GetDegreesLatitude()
			lng := loc.GetDegreesLongitude()
			media = append(media, model.Media {
				Latitude: lat,
				Longitude: lng,
			})
			if txtMessage != "" {
				txtMessage += ". "
			}
			txtMessage += "Location: " + fmt.Sprintf("%f", lat) +", "+ fmt.Sprintf("%f", lng)
		}

		message := model.Event {
			ID: v.Info.ID,
			Chat: chatId, // group id or phone id
			SenderId: model.GetPhoneNumber(v.Info.Sender.String()), // phone id
			SenderName: v.Info.PushName,
			Time: v.Info.Timestamp.Unix(),
			IsGroup: v.Info.IsGroup,
			IsFromMe: v.Info.IsFromMe,
			Type: theType,
			MediaType: mediaType,
			Text: txtMessage,
			IsViewOnce: v.IsViewOnce,
			IsViewOnceV2: v.IsViewOnceV2,
			IsViewOnceV2Extension: v.IsViewOnceV2Extension,
			IsLottieSticker: v.IsLottieSticker,
			IsDocumentWithCaption: v.IsDocumentWithCaption,
			ReplyToPost: replyToPost,
			ReplyToUser: replyToUser,
			Media: media,
		}

		jsonResponse2, err := json.MarshalIndent(message, "", "")
		fmt.Println(string(jsonResponse2))

		// untuk group
		// if v.Info.IsGroup {
		// 	message = Message {
		// 	}
		// }

		// untuk personal
		// db.CheckDatabase()

		// for key := range clients {
		// 	//fmt.Println("whoami:", key)
		// 	fmt.Println("end client", clients[key])
		// }

		// if !v.Info.IsFromMe || v.Message.GetConversation() != "" ||
		// 	v.Message.GetImageMessage().GetCaption() != "" ||
		// 	v.Message.GetVideoMessage().GetCaption() != "" ||
		// 	v.Message.GetDocumentMessage().GetCaption() != "" {
		// 	id := v.Info.ID
		// 	chat := v.Info.Sender.String()
		// 	timestamp := v.Info.Timestamp
		// 	text := v.Message.GetConversation()
		// 	group := v.Info.IsGroup
		// 	isfrome := v.Info.IsFromMe
		// 	//doc := v.Message.GetDocumentMessage()
		// 	captionMessage := v.Message.GetImageMessage().GetCaption()
		// 	videoMessage := v.Message.GetVideoMessage().GetCaption()
		// 	docMessage := v.Message.GetDocumentMessage().GetCaption()
		// 	//docCaption := v.Message.GetDocumentMessage().GetTitle()
		// 	name := v.Info.PushName
		// 	to := v.Info.PushName

		// 	//to : = v.Info.na
		// 	thumbnail := v.Message.ImageMessage.GetJPEGThumbnail()
		// 	thumbnailvideo := v.Message.VideoMessage.GetJPEGThumbnail()
		// 	thumbnaildoc := v.Message.DocumentMessage.GetJPEGThumbnail()
		// 	url := v.Message.ImageMessage.GetURL()
		// 	mimeTipe := v.Message.ImageMessage.GetMimetype()

		// 	tipe := v.Info.Type
		// 	isdocument := v.IsDocumentWithCaption
		// 	//chatText := v.Info.Chat
		// 	mediatype := v.Info.MediaType
		// 	//smtext := v.Message.Conversation()
		// 	//fmt.Println("ID: %s, Chat: %s, Time: %d, Text: %s\n", to, mediatype, isdocument, chat, timestamp, text, group, isfrome, tipe)
		// 	//fmt.Println("info repley", reply, coba)

		// 	// Assuming replies are stored within a field named Replies
		// 	//fmt.Println("tipe messages", tipe, docCaption, isdocument, doc, mediatype, captionMessage, videoMessage, docMessage)
		// 	mu.Lock()
		// 	defer mu.Unlock() // Ensure mutex is always unlocked when the function returns
		// 	messages = append(messages, response.Message{
		// 		ID:             id,
		// 		Chat:           chat,
		// 		Time:           timestamp.Unix(),
		// 		Text:           text,
		// 		Group:          group,
		// 		Mediatipe:      mediatype,
		// 		IsDocument:     isdocument,
		// 		Tipe:           tipe,
		// 		IsFromMe:       isfrome,
		// 		Caption:        captionMessage,
		// 		VideoMessage:   videoMessage,
		// 		DocMessage:     docMessage,
		// 		Name:           name,
		// 		From:           chat,
		// 		To:             to,
		// 		Url:            url,
		// 		Thumbnail:      base64.StdEncoding.EncodeToString(thumbnail),
		// 		MimeTipe:       mimeTipe,
		// 		Thumbnaildoc:   base64.StdEncoding.EncodeToString(thumbnaildoc),
		// 		Thumbnailvideo: base64.StdEncoding.EncodeToString(thumbnailvideo),
		// 		//MimeType:     *mimesType,
		// 		//CommentMessage: comment,
		// 		//Replies: reply,
		// 		// Add replies to the message if available
		// 		// Replies: v.Message.Replies,
		// 	})
		// }

		/*
			payload := response.Message{
				ID:             v.Info.ID,
				Chat:           v.Info.Sender.String(),
				Time:           v.Info.Timestamp.Unix(),
				Text:           v.Message.GetConversation(),
				Group:          v.Info.IsGroup,
				IsFromMe:       v.Info.IsFromMe,
				Caption:        v.Message.GetImageMessage().GetCaption(),
				VideoMessage:   v.Message.GetVideoMessage().GetCaption(),
				DocMessage:     v.Message.GetDocumentMessage().GetCaption(),
				MimeTipe:       v.Message.GetImageMessage().GetMimetype(),
				Name:           v.Info.PushName,
				To:             v.Info.PushName,
				Url:            v.Message.GetImageMessage().GetUrl(),
				Thumbnail:      base64.StdEncoding.EncodeToString(v.Message.GetImageMessage().GetJpegThumbnail()),
				Thumbnailvideo: base64.StdEncoding.EncodeToString(v.Message.GetVideoMessage().GetJpegThumbnail()),
				Thumbnaildoc:   base64.StdEncoding.EncodeToString(v.Message.GetDocumentMessage().GetJpegThumbnail()),
				Tipe:           v.Info.Type,
				IsDocument:     v.IsDocumentWithCaption,
				Mediatipe:      v.Info.MediaType,
			}
			webhookURL := "https://webhook.site/aa9bbb63-611c-4d7a-97cd-f4eb6d4b775d"
			err := sendPayloadToWebhook(payload, webhookURL)
			if err != nil {
				fmt.Printf("Failed to send payload to webhook: %v\n", err)
			}
		*/

	case *events.PairSuccess:
		// fmt.Println("pari succeess", v.ID.User)
		fmt.Println("--- pairing success", v.ID.User)
		phonekey :=  model.GetPhoneNumber(v.ID.String())
		for key, client := range model.Clients {
			if !strings.HasPrefix(key, "DEV") {
				continue
			}

			iphone := model.GetPhoneNumber(client.Client.Store.ID.String())
			fmt.Println("comparing ", phonekey, iphone)
			if iphone == phonekey {
				db.InsertUserDevice(model.UserDevice{
					UserId:    client.User,
					DeviceJid: phonekey,
				})
				client.ExpiredTime = 0
				model.Clients[phonekey] = client
				delete(model.Clients, key)
				break;
			}
		}
		// initialClient()
		
	case *events.HistorySync:
		fmt.Println("Received a history sync")
		/*for _, conv := range v.Data.GetConversations() {
			for _, historymsg := range conv.GetMessages() {
				chatJID, _ := types.ParseJID(conv.GetId())
				evt, err := client.ParseWebMessage(chatJID, historymsg.GetMessage())
				if err != nil {
					log.Println(err)
				}
				eventHandler(evt)
			}
		}*/

	case *events.LoggedOut:
		fmt.Println("------ Logout from mobile device ----")
		for _, client := range model.Clients {
			cid := model.GetPhoneNumber(client.Client.Store.ID.String())
			if !client.Client.IsLoggedIn() {
				delete(model.Clients, cid)
				db.DeleteUserDevice(cid)
			}
		}
		//initialClient()

	case *events.Receipt:
		
		// fmt.Println("----- terima")
		if v.Type == events.ReceiptTypeRead || v.Type == events.ReceiptTypeReadSelf {
			fmt.Printf("%v was read by %s at %s\n", v.MessageIDs, v.SourceString(), v.Timestamp)
			// Membuat payload untuk webhook
			/*webhookPayload := model.ReadReceipt{
				MessageID: v.MessageIDs[0],
				ReadBy:    v.SourceString(),
				Time:      v.Timestamp.UnixMilli(),
			}*/
			// Mengirimkan payload ke webhook
			//webhookURL := "http://localhost:8080/webhook"
			/*err := sendPayloadToWebhook(string(v.Type), webhookURL)
			if err != nil {
				fmt.Printf("Failed to send read receipt to webhook: %v\n", err)
			}*/
		} else {
			// fmt.Println(v)
		}

		/*case *events.Receipt:
		if v.Type == types.ReceiptTypeRead || v.Type == types.ReceiptTypeReadSelf {
			fmt.Println("%v was read by %s at %s", v.MessageIDs, v.SourceString(), v.Timestamp)
		} else if v.Type == types.ReceiptTypeDelivered {
			fmt.Println("%s was delivered to %s at %s", v.MessageIDs[0], v.SourceString(), v.Timestamp)
		}
		*/

	}
}