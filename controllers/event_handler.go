package controllers

import (
	// "io"
	// "os"

	"fmt"

	//"log"

	"path/filepath"
	"strings"
	"sync"
	"time"

	// "net/http"
	//"encoding/json"

	"wagobot.com/db"
	//"wagobot.com/helpers"

	// "encoding/base64"
	// "wagobot.com/base"
	"wagobot.com/model"
	// "wagobot.com/response"

	// "go.mau.fi/whatsmeow/store"
	// "go.mau.fi/whatsmeow/store/sqlstore"

	"go.mau.fi/whatsmeow/types/events"
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

func EventHandler(evt interface{}, cclient model.CustomClient) {

	switch v := evt.(type) {
	case *events.Message:

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
		fmt.Println(payload)
		// initialClient()

	case *events.HistorySync:

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
				fmt.Println(payload)

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
