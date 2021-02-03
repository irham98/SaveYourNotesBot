package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/irham/agung/notes-bot/gdrive"
	"github.com/irham/agung/notes-bot/handler"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

const (
	configFile = "config.yaml"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}

func init() {

}

func main() {
	// ======================= Configuration ======================= //
	config, err := loadConfig(configFile)
	exitOnErr(err)

	// ================= Google Drive Integration ================== //
	gConfig, err := gdrive.Setup(config.Credentials)
	exitOnErr(err)

	// =================== Bot Initialization ======================= //
	bot, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	exitOnErr(err)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(config.Telegram.WebhookUrl))
	exitOnErr(err)

	updates := bot.ListenForWebhook("/")

	h := handler.New(bot, log)

	http.HandleFunc("/auth", func(writer http.ResponseWriter, request *http.Request) {
		queries := request.URL.Query()
		code := queries.Get("code")
		log.Info(code)

		tok, err := gConfig.Exchange(context.TODO(), code)
		if err != nil {
			log.Fatalf("Unable to retrieve token from web %v", err)
		}

		gdrive.SaveToken("token.json", tok)

		writer.WriteHeader(200)
	})

	go func() {
		for update := range updates {
			if update.Message.Photo != nil {
			}
			switch update.Message.Text {
			case "/start":
				h.MessageStart(&update)
			case "/auth":
				h.MessageAuth(gConfig, &update)
			default:
				h.MessageUnknown(&update)
			}
		}
	}()

	// ===================== HTTP Listen ============================ //

	err = http.ListenAndServe(fmt.Sprintf(":%v", config.Port), nil)
	exitOnErr(err)

	//router := gin.Default()
	//router.POST("/start", func(c *gin.Context) {
	//	fmt.Println(c.Request.Body)
	//})
	//
	//router.POST("/", func(c *gin.Context) {
	//	response := &Response{}
	//	if err := c.BindJSON(response); err != nil {
	//		c.AbortWithStatus(http.StatusBadRequest)
	//	}
	//	fmt.Println(response.Message.Text)
	//	if response.Message.Photo != nil {
	//		err := SaveImage(response.Message.Photo[0].FileID)
	//		fmt.Println(err)
	//		c.AbortWithStatus(http.StatusInternalServerError)
	//	}
	//	c.JSON(200, gin.H{"message": "received"})
	//	return
	//})
	//
	//if err := router.Run(":80"); err != nil {
	//	panic(err)
	//}
}

//
//func main() {
//	// ======================= Configuration ======================= //
//	config, err := loadConfig(configFile)
//	exitOnErr(err)
//
//	b, err := ioutil.ReadFile(config.Credentials)
//	if err != nil {
//		log.Fatalf("Unable to read client secret file: %v", err)
//	}
//
//	// If modifying these scopes, delete your previously saved token.json.
//	gConfig, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
//	if err != nil {
//		log.Fatalf("Unable to parse client secret file to config: %v", err)
//	}
//	client := gdrive.GetClient(gConfig)
//
//	srv, err := drive.New(client)
//	if err != nil {
//		log.Fatalf("Unable to retrieve Drive client: %v", err)
//	}
//
//	r, err := srv.Files.List().PageSize(10).
//		Fields("nextPageToken, files(id, name)").Do()
//	if err != nil {
//		log.Fatalf("Unable to retrieve files: %v", err)
//	}
//	fmt.Println("Files:")
//	if len(r.Files) == 0 {
//		fmt.Println("No files found.")
//	} else {
//		for _, i := range r.Files {
//			fmt.Printf("%s (%s)\n", i.Name, i.Id)
//		}
//	}
//}
//
// for fatal error on server initialization
func exitOnErr(err error) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
