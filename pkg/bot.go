package discordbot

// IF broken import, use go get github.com/bwmarrin/discordgo
// DOCCUMENTATION REFFERENCE: https://github.com/bwmarrin/discordgo/blob/master/examples/pingpong/main.go
import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/bwmarrin/discordgo"
	"github.com/cody6750/discordbot/pkg/functions"
	bothandler "github.com/cody6750/discordbot/pkg/handlers"
	services "github.com/cody6750/discordbot/pkg/services/aws"
)

const (
	defaultDiscordTokenID = "discord/token"
	defaultDiscordToken   = ""
	defaultMaxRetries     = 5
	defaultRegion         = "us-east-1"
	defaultLocalRun       = false
)

var (
	// Need to store token in secrets manager or something
	// discordToken string = "ODE5NDgyNjU1NjczODEwOTU0.YEnQsg.Pmu9ppRwPpvJl6ebPS2ryQCyrnc"
	localRun       bool   = defaultLocalRun
	maxRetries     int    = defaultMaxRetries
	discordTokenID string = defaultDiscordTokenID
	discordToken   string = defaultDiscordToken
	region         string = defaultRegion
	err            error
	secretsSvc     *secretsmanager.SecretsManager
)

//Init initializes the discord bot
func Init() {

	fmt.Println("Discord bot initializing......")
	getEnvVariables()

	if !localRun {
		initAWS()
		discordToken, err = getDiscordTokenFromAWS(secretsSvc, discordTokenID)
		if err != nil {
			log.Fatalf("getDiscordTokenFRrom AWS failed with ERROR: %v", err)
		}
	}
	err := checkEnvVariables()
	if err != nil {
		log.Fatalf("checkEnvVariables failed with ERROR: %v", err)
	}
	// Create a new Discord session using the provided bot token.
	discord, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("Error has occured")
		return
	}
	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(bothandler.Ready)
	discord.AddHandler(bothandler.MessageCreate)
	discord.AddHandler(bothandler.Disconnect)
	discord.AddHandler(bothandler.SlashCommands)

	// In this example, we only care about receiving message events.
	discord.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	go functions.StartTracking(discord, functions.GetChannels(discord, "tracking"))
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func initAWS() {
	configs := aws.Config{
		Region:     aws.String(region),
		MaxRetries: aws.Int(maxRetries),
	}
	session := session.Must(session.NewSession(&configs))
	secretsSvc = secretsmanager.New(session)
}

func getDiscordTokenFromAWS(secretsSvc *secretsmanager.SecretsManager, discordTokenPath string) (string, error) {
	token, err := services.GetSecret(secretsSvc, discordTokenPath)
	if err != nil {
		return "", err
	}
	return token, nil
}

func getEnvVariables() {
	if os.Getenv("AWS_MAX_RETRIES") != "" {
		maxRetries, err = strconv.Atoi(os.Getenv("AWS_MAX_RETRIES"))
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
	if os.Getenv("AWS_REIGON") != "" {
		region = os.Getenv("AWS_REGION")
	}

	if os.Getenv("LOCAL_RUN") != "" {
		switch r := os.Getenv("LOCAL_RUN"); strings.ToLower(r) {
		case "true":
			localRun = true
		case "false":
			localRun = false
		default:
			log.Fatalf("LOCAL_RUN environment variables must be `true` | `false` ")
		}
	}
	if os.Getenv("DISCORD_TOKEN_ID") != "" {
		discordTokenID = os.Getenv("DISCORD_TOKEN_ID")
	}
	if os.Getenv("DISCORD_TOKEN") != "" {
		discordToken = os.Getenv("DISCORD_TOKEN")
	}
}

func checkEnvVariables() error {
	if discordToken == "" {
		return fmt.Errorf("Discord token is not set, please set DISCORD_TOKEN environment variable or check AWS")
	}
	return nil
}
