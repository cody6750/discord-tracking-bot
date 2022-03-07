package functions

import (
	"bytes"
	"io"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func Log(loger *logrus.Logger, w io.Writer, f func(...interface{}), message string) {
	loger.SetOutput(w)
	f(message)
}
func LogToDiscord(loger *logrus.Logger, sess *discordgo.Session, c *discordgo.Channel, f func(...interface{}), message string) {
	var b bytes.Buffer
	Log(loger, &b, f, message)
	sess.ChannelMessageSend(c.ID, b.String())

}

func LogToStdout(loger *logrus.Logger, f func(...interface{}), message string) {
	Log(loger, os.Stdout, f, message)
}

func LogToDiscordAndStdOut(loger *logrus.Logger, sess *discordgo.Session, c *discordgo.Channel, f func(...interface{}), message string) {
	LogToStdout(loger, f, message)
	loger.SetFormatter(&logrus.TextFormatter{ForceColors: false, FullTimestamp: true})
	LogToDiscord(loger, sess, c, f, message)
	loger.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})

}
