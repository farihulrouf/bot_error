package controllers

import (
	// "io"
	// "os"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	// "net/http"
	"encoding/json"
	"net/url"

	"wagobot.com/db"
	"wagobot.com/helpers"

	// "encoding/base64"
	// "wagobot.com/base"
	"wagobot.com/model"
	// "wagobot.com/response"
	"go.mau.fi/whatsmeow"
	// "go.mau.fi/whatsmeow/store"
	// "go.mau.fi/whatsmeow/store/sqlstore"
	wtypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func getExtension(filename string) string {
	ext := filepath.Ext(filename)
	return ext
}

func removeExtension(filename string) string {
	// Get the file extension
	ext := filepath.Ext(filename)
	// Remove the extension from the filename
	return strings.TrimSuffix(filename, ext)
}

func setLastMimetype(mimetype string) string {
	tmp := strings.Split(mimetype, "/")
	return tmp[1]
}

var processedMessages = struct {
	sync.RWMutex
	messages map[string]struct{}
}{
	messages: make(map[string]struct{}),
}

// Fungsi untuk memeriksa dan menyimpan ID pesan
func isProcessed(id string) bool {
	processedMessages.RLock()
	defer processedMessages.RUnlock()
	_, exists := processedMessages.messages[id]
	return exists
}

func markAsProcessed(id string) {
	processedMessages.Lock()
	defer processedMessages.Unlock()
	processedMessages.messages[id] = struct{}{}
}

func uploadToSpace(pathToFile string, fileData []byte, mimetype string) (string, error) {

	bucketName := "dragonfly"
	region := "sgp1"
	accessKey := model.SpaceConfig.AccessKey
	secretKey := model.SpaceConfig.SecretKey

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region), // Use the DigitalOcean region (e.g., nyc3)
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           model.SpaceConfig.Endpoint,
				SigningRegion: region,
			}, nil
		})),
	)
	if err != nil {
		return "", fmt.Errorf("failed to load config: %v", err)
	}

	// Create an S3 client
	client := s3.NewFromConfig(cfg)

	// Upload the file
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(pathToFile),
		Body:        bytes.NewReader(fileData),
		ACL:         types.ObjectCannedACLPublicRead, // Set to public or private based on your requirement
		ContentType: aws.String(mimetype),
	})
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	strurl := model.SpaceConfig.Endpoint + "/" + bucketName + "/" + pathToFile

	return strurl, nil
}

func saveMedia(
	client *whatsmeow.Client,
	mediaMessage whatsmeow.DownloadableMessage,
	chatId string, strdate string, filename string, targeturl string, mimetype string) string {

	if filename == "" {
		parsedURL, err := url.Parse(targeturl)
		if err != nil {
			fmt.Printf("Error parsing URL: %v\n", err)
			return ""
		}
		filename = path.Base(parsedURL.Path)
	}
	// filename = removeExtension(filename) +"-"+ setLastMimetype(mimetype)
	filename = removeExtension(filename) + "-" + getExtension(filename)

	byteData, err := client.Download(mediaMessage)
	if err != nil {
		fmt.Println("Error downloading encrypted image:", err)
		return ""
	}

	path := "media/wa/p/" + chatId + "/" + strdate + "/" + filename

	myurl, _ := uploadToSpace(path, byteData, mimetype)
	fmt.Println("FILEEEEE", myurl)

	return myurl
}

func saveProfilePicture(client *whatsmeow.Client, theJID wtypes.JID) string {

	params := &whatsmeow.GetProfilePictureParams{
		// JID:     theJID,
		// IsCommunity: false,
	}

	profilePictureURL, err := client.GetProfilePictureInfo(theJID, params) // false for high-res, true for low-res
	if err != nil {
		fmt.Println("Error getting profile picture:", err)
		return ""
	}

	response, err := http.Get(profilePictureURL.URL)
	if err != nil {
		return ""
	}
	defer response.Body.Close()

	byteData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}

	filename := model.GetPhoneNumber(theJID.String()) + "-jpg"

	path := "media/wa/a/" + filename

	myurl, _ := uploadToSpace(path, byteData, "image/jpg")
	// fmt.Println("FILEEEEE", myurl)

	return myurl
}

func EventHandler(evt interface{}, cclient model.CustomClient) {

	switch v := evt.(type) {
	case *events.Message:

		fmt.Println("------ new message ")
		// fmt.Println(evt)
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v", r)
			}
		}()
		var media model.Media
		chatId := v.Info.Chat.String()
		theType := "post"
		replyToPost := ""
		replyToUser := ""
		mediaType := v.Info.Type
		strdate := v.Info.Timestamp.Format("20060102")
		// Jika ID pesan sudah diproses, hentikan proses lebih lanjut
		if isProcessed(v.Info.ID) {
			fmt.Println("Message already processed, skipping...")
			return
		}

		// Tandai ID pesan sebagai sudah diproses
		markAsProcessed(v.Info.ID)

		// if !v.Info.IsGroup {
		chatId = model.GetPhoneNumber(chatId)
		// }
		idChat := helpers.ConvertToLettersDetailed(cclient.Phone, chatId, v.Info.IsGroup)
		fmt.Println("---------- message  link ------")
		//fmt.Println(v.Message)

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
			turl := ext.GetCanonicalURL()
			if turl != "" {
				mediaType = "url"
				if txtMessage != "" {
					txtMessage += ". "
				}
				txtMessage += ext.GetText()
			}
		} else {
			txtMessage = v.Message.GetConversation()
		}

		if v.Message.ImageMessage != nil {
			img := v.Message.GetImageMessage()
			imgCaption := img.GetCaption()
			mediaType := "image"

			mediaUrl := saveMedia(cclient.Client, img, chatId, strdate, "", img.GetURL(), img.GetMimetype())

			media = model.Media{
				Url:        mediaUrl,
				Type:       mediaType,
				Caption:    imgCaption,
				MimeType:   img.GetMimetype(),
				FileLength: img.GetFileLength(),
			}
			if txtMessage != "" {
				txtMessage += ". "
			}
			txtMessage += imgCaption
		}

		if v.Message.DocumentMessage != nil {
			doc := v.Message.GetDocumentMessage()
			docCaption := doc.GetCaption()
			mediaType := "file"

			mediaUrl := saveMedia(cclient.Client, doc, chatId, strdate, doc.GetFileName(), doc.GetURL(), doc.GetMimetype())

			media = model.Media{
				Url:        mediaUrl,
				Type:       mediaType,
				Caption:    docCaption,
				FileName:   doc.GetFileName(),
				FileLength: doc.GetFileLength(),
				MimeType:   doc.GetMimetype(),
			}

			if txtMessage != "" {
				txtMessage += ". "
			}
			txtMessage += docCaption
		}

		if v.Message.AudioMessage != nil {
			aud := v.Message.GetAudioMessage()
			audCaption := ""
			mediaType = "audio"

			mediaUrl := saveMedia(cclient.Client, aud, chatId, strdate, "", aud.GetURL(), aud.GetMimetype())

			media = model.Media{
				Url:        mediaUrl,
				Type:       mediaType,
				Caption:    audCaption,
				MimeType:   aud.GetMimetype(),
				Seconds:    aud.GetSeconds(),
				FileLength: aud.GetFileLength(),
			}

			if txtMessage != "" {
				txtMessage += ". "
			}
			txtMessage += audCaption
		}

		if v.Message.VideoMessage != nil {
			vid := v.Message.GetVideoMessage()
			vidCaption := vid.GetCaption()
			mediaType = "video"

			mediaUrl := saveMedia(cclient.Client, vid, chatId, strdate, "", vid.GetURL(), vid.GetMimetype())

			media = model.Media{
				Url:      mediaUrl,
				Type:     mediaType,
				Caption:  vidCaption,
				MimeType: vid.GetMimetype(),
				// Thumbnail: vid.GetJPEGThumbnail(),
				FileLength: vid.GetFileLength(),
			}

			if txtMessage != "" {
				txtMessage += ". "
			}
			txtMessage += vidCaption
		}

		if v.Message.ContactMessage != nil {
			ctc := v.Message.GetContactMessage()
			mediaType = "contact"
			media = model.Media{
				Name:    ctc.GetDisplayName(),
				Contact: ctc.GetVcard(),
			}

			if txtMessage != "" {
				txtMessage += ". "
			}
			txtMessage += ctc.GetVcard()
		}

		if v.Message.PollCreationMessageV3 != nil {
			pol := v.Message.GetPollCreationMessageV3()
			polName := pol.GetName()
			mediaType = "polling"
			jsonData, _ := json.Marshal(pol)
			media = model.Media{
				Poll: string(jsonData),
			}
			if txtMessage != "" {
				txtMessage += ". "
			}
			txtMessage += polName
		}
		if v.Message.GetLiveLocationMessage() != nil {
			loc := v.Message.GetLiveLocationMessage()
			lat := loc.GetDegreesLatitude()
			lng := loc.GetDegreesLongitude()
			mediaType = "location"
			media = model.Media{
				Latitude:  lat,
				Longitude: lng,
			}
			if txtMessage != "" {
				txtMessage += ". "
			}
			txtMessage += "Location: " + fmt.Sprintf("%f", lat) + ", " + fmt.Sprintf("%f", lng)
		}
		if v.Message.LocationMessage != nil {
			loc := v.Message.GetLocationMessage()
			lat := loc.GetDegreesLatitude()
			lng := loc.GetDegreesLongitude()
			mediaType = "location"
			media = model.Media{
				Latitude:  lat,
				Longitude: lng,
			}
			if txtMessage != "" {
				txtMessage += ". "
			}
			txtMessage += "Location: " + fmt.Sprintf("%f", lat) + ", " + fmt.Sprintf("%f", lng)
		}

		avatar := saveProfilePicture(cclient.Client, v.Info.Sender)

		senderPhone := model.GetPhoneNumber(v.Info.Sender.String())
		senderName := v.Info.PushName

		message := model.Event{
			ID:                    idChat,
			Chat:                  chatId,      // group id or phone id
			SenderId:              senderPhone, // phone id
			SenderName:            senderName,
			SenderAvatar:          avatar,
			Time:                  v.Info.Timestamp.Unix(),
			IsGroup:               v.Info.IsGroup,
			IsFromMe:              v.Info.IsFromMe,
			Type:                  theType,
			MediaType:             mediaType,
			Text:                  txtMessage,
			IsViewOnce:            v.IsViewOnce,
			IsViewOnceV2:          v.IsViewOnceV2,
			IsViewOnceV2Extension: v.IsViewOnceV2Extension,
			IsLottieSticker:       v.IsLottieSticker,
			IsDocumentWithCaption: v.IsDocumentWithCaption,
			ReplyToPost:           replyToPost,
			ReplyToUser:           replyToUser,
			Media:                 media,
			IDChat:                v.Info.ID,
		}

		// payload, _ := json.MarshalIndent(message, "", "")
		// fmt.Println("HASILLL -----------")
		// fmt.Println(message)
		// fmt.Println("---- Active Webhook url", model.DefaultWebhook)

		payload := model.PayloadWebhook{
			Section: "single_message",
			Data:    message,
		}

		err := sendPayloadToWebhook(model.DefaultWebhook, payload)
		if err != nil {
			fmt.Printf("Failed to send payload to webhook: %v\n", err)
		}

		var groupid *string
		groupid = nil
		if v.Info.IsGroup {
			groupid = &chatId
		}

		// untuk sender
		members := []model.Member{}
		members = append(members, model.Member{
			ID:     senderPhone,
			Name:   senderName,
			Avatar: avatar,
			Phone:  senderPhone,
			// IsAdmin: member.IsAdmin,
			// IsSuperAdmin: member.IsSuperAdmin,
			GroupID: *groupid,
		})

		payload = model.PayloadWebhook{
			Section: "senders",
			Data:    members,
		}

		err = sendPayloadToWebhook(model.DefaultWebhook, payload)
		if err != nil {
			fmt.Printf("Failed to send payload to webhook: %v\n", err)
		}

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
		phoneref := ""
		ref := ""
		phonekey := model.GetPhoneNumber(v.ID.String())
		for key, client := range model.Clients {
			if !strings.HasPrefix(key, "DEV") {
				continue
			}

			params := strings.Split(key, "-")
			phoneref = params[1]
			ref = params[2]

			iphone := model.GetPhoneNumber(client.Client.Store.ID.String())
			fmt.Println("comparing ", phonekey, iphone, phoneref, ref)
			if phonekey == iphone {
				if phonekey == phoneref {
					db.InsertUserDevice(model.UserDevice{
						UserId:    client.User,
						DeviceJid: phonekey,
					})
					client.ExpiredTime = 0
					client.Phone = phonekey
					model.Clients[phonekey] = client
					delete(model.Clients, key)
				} else {
					time.Sleep(5 * time.Second)
					if client.Client.IsLoggedIn() {
						client.Client.Logout()
					}
					fmt.Printf("---> Mismatch number, deleting %s\n", key)
					delete(model.Clients, key)
				}
				break
			}
		}

		payload := model.PayloadWebhook{
			Section: "device_added",
			Data: model.PhoneVerifyParams{
				Phone:    phonekey,
				PhoneRef: phoneref,
				Ref:      ref,
			},
		}

		err := sendPayloadToWebhook(model.DefaultWebhook, payload)
		if err != nil {
			fmt.Printf("Failed to send payload to webhook: %v\n", err)
		}
		// initialClient()

	case *events.HistorySync:
		fmt.Println("---- sync history -----")
		for _, conv := range v.Data.GetConversations() {
			for _, historymsg := range conv.GetMessages() {
				fmt.Println(historymsg)
				chatJID, _ := wtypes.ParseJID(conv.GetId())
				evt, _ := cclient.Client.ParseWebMessage(chatJID, historymsg.GetMessage())
				// if err != nil {
				// 	log.Println(err)
				// }
				EventHandler(evt, cclient)
			}
		}
		fmt.Println("--- sync done -----")

	case *events.LoggedOut:
		fmt.Println("------ Logout from mobile device ----")
		for _, client := range model.Clients {
			cid := model.GetPhoneNumber(client.Client.Store.ID.String())
			if !client.Client.IsLoggedIn() {
				payload := model.PayloadWebhook{
					Section: "device_removed",
					Data: model.PhoneParams{
						Phone: cid,
					},
				}
				err := sendPayloadToWebhook(model.DefaultWebhook, payload)
				if err != nil {
					fmt.Printf("Failed to send payload to webhook: %v\n", err)
				}
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
