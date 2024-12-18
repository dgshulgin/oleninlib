package main

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(os.Stdout).With().Str("component", "OLENIN-BOT").Timestamp().Logger()

	err := run(log)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

func run(log zerolog.Logger) error {
	log.Info().Msg("запущен")

	//err :=
	godotenv.Load()
	// if err != nil {
	// 	return err
	// }
	apiToken, exist := os.LookupEnv("TG_TOKEN")
	if !exist || len(apiToken) == 0 {
		return errors.New("TG_TOKEN не определен")
	}

	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		return err
	}

	_, exist = os.LookupEnv("DEBUG_MODE") //важно наличие, но не значение
	if exist {
		bot.Debug = true
	}

	updateConfig := tgbotapi.UpdateConfig{
		Timeout: 60,
	}
	updates := bot.GetUpdatesChan(updateConfig)

	var wg sync.WaitGroup

	sigCh := make(chan os.Signal, 1) //для обработки CTRL+C
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	quitCh := make(chan bool, 1) // для завершения главного цикла обработки

	// обработка CTRL+C
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-sigCh
		quitCh <- true
		log.Info().Msg("Перехвачен CTRL+C. Завершение работы")
	}()

	// главный цикл обработки сообщений
	wg.Add(1)
	go Update(bot, log, &wg, updates, quitCh)

	wg.Wait()
	log.Info().Msg("завершен успешно")

	return nil
}

func Update(
	bot *tgbotapi.BotAPI,
	log zerolog.Logger,
	wg *sync.WaitGroup,
	updates tgbotapi.UpdatesChannel,
	quit <-chan bool) {

	defer wg.Done()

	for {
		select {
		case quitFlag := <-quit:
			if quitFlag {
				log.Info().Msg("завершаю главный цикл")
				return
			}

		case update := <-updates:
			if update.Message == nil {
				continue
			}

			switch update.Message.Command() {
			case "help":
				processHelpCommand(bot, log, update.Message)
			default:
				processEchoText(bot, log, update.Message)
			}
		}
	}
}

func processHelpCommand(
	bot *tgbotapi.BotAPI,
	log zerolog.Logger,
	inputMessage *tgbotapi.Message) {

	out := tgbotapi.NewMessage(
		inputMessage.Chat.ID,
		"/help - help")
	if _, err := bot.Send(out); err != nil {
		log.Error().Err(err).Send()
	}
}
func processEchoText(
	bot *tgbotapi.BotAPI,
	log zerolog.Logger,
	inputMessage *tgbotapi.Message) {

	out := tgbotapi.NewMessage(inputMessage.Chat.ID, "OLB says: "+inputMessage.Text)
	out.ReplyToMessageID = inputMessage.MessageID

	if _, err := bot.Send(out); err != nil {
		log.Error().Err(err).Send()
	}
}
