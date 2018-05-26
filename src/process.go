package main

import (
	cmc "github.com/bearbeard/go-coinmarketcap"
	"bytes"
	"fmt"
	"strings"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

const (
	_5   = "5"
	_10  = "10"
	_50  = "50"
	_100 = "100"
	CHECK = "check"
	START = "start"
	HELP = "help"
	START_TEXT = "Hi. I am BearBeardBot.\n"+
		"I can help you check cryptocurrency via CoinMarketCap.\n"+
		"Type /help to learn my commands."
	HELP_TEXT = "You can type a ticker name in input field to get info about that ticker " +
		"or you can use one of these commands:\n\n" +
		"/check - reveal options to check top currency"
	UNKNOWN_COMMAND_TEXT = "Sorry, I can't remember this command. Please, try another one or type /help."
)

var (
	buttons    = []string{_5, _10, _50, _100}
	errorTitle = "Error"
	errorText  = "Sorry, something goes wrong with me. You took an error."
)

func processUpdate(bot *tgbotapi.BotAPI) {
	var msg tgbotapi.MessageConfig
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10
	updates, err := bot.GetUpdatesChan(u)
	checkError(err)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				msg = invokeCommand(update)
			} else {
				msg = invokeTextCommand(update)
			}
		} else if update.InlineQuery != nil {
			inlineConfig := invokeInlineCommand(update)
			_, err := bot.AnswerInlineQuery(inlineConfig)
			checkError(err)
		} else if update.CallbackQuery != nil {
			msg = invokeCallbackCommand(update)
		}
		msg.ParseMode = "markdown"
		bot.Send(msg)
	}
}

func invokeCommand(update tgbotapi.Update) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	switch update.Message.Command() {
	case CHECK:
		{
			msg = check(update)
		}
	case START:
		{
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, START_TEXT)
		}
	case HELP:
		{
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, HELP_TEXT)
		}
	default:
		{
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, UNKNOWN_COMMAND_TEXT)
		}
	}
	return msg
}

func invokeTextCommand(update tgbotapi.Update) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	ticker := update.Message.Text
	reply, err := checkTicker(ticker)
	if err != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, errorText)
		log.Println(err)
	} else {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	}
	return msg
}

func invokeCallbackCommand(update tgbotapi.Update) tgbotapi.MessageConfig {
	var reply string
	var err error
	var msg tgbotapi.MessageConfig
	var data = update.CallbackQuery.Data
	switch data {
	case _5:
		{
			reply, err = checkTop(5)
		}
	case _10:
		{
			reply, err = checkTop(10)
		}
	case _50:
		{
			reply, err = checkTop(50)
		}
	case _100:
		{
			reply, err = checkTop(100)
		}
	}
	if err != nil {
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, errorText)
	} else {
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
	}
	return msg
}

func invokeInlineCommand(update tgbotapi.Update) tgbotapi.InlineConfig {
	var articles []interface{}
	var msg tgbotapi.InlineQueryResultArticle
	query := update.InlineQuery.Query
	reply, err := checkTicker(query)

	if err != nil {
		msg = tgbotapi.NewInlineQueryResultArticleMarkdown(update.InlineQuery.ID, errorTitle, errorText)
		log.Println(err)
	} else {
		msg = tgbotapi.NewInlineQueryResultArticleMarkdown(update.InlineQuery.ID, reply, reply)
	}

	articles = append(articles, msg)
	return tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       articles,
	}
}

func check(update tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Check top")
	keyboard := tgbotapi.InlineKeyboardMarkup{}

	for _, button := range buttons {
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(button, button)
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	msg.ReplyMarkup = keyboard
	return msg
}

func checkTop(lim int) (string, error) {
	tickers, err := cmc.Tickers(&cmc.TickersOptions{
		Start:   0,
		Limit:   lim,
		Convert: "USD",
	})

	if err != nil {
		return "", err
	}

	b := new(bytes.Buffer)
	fmt.Fprintf(b, "*Top %d:*\n", lim)

	for _, ticker := range tickers {
		symbol := ticker.Symbol
		price := ticker.Quotes["USD"].Price
		fmt.Fprintf(b, "%d) %s = %f\n", ticker.Rank, symbol, price)
	}
	return b.String(), nil
}

func checkTicker(name string) (string, error) {
	name = strings.ToUpper(name)
	usd := "USD"
	ticker, err := cmc.Ticker(&cmc.TickerOptions{
		Symbol:  name,
		Convert: usd,
	})

	if err != nil {
		return "", err
	}

	b := new(bytes.Buffer)
	fmt.Fprintf(b, "*%s (%s)*\n\n", ticker.Name, ticker.Symbol)
	fmt.Fprintf(b, "*Price*: %f\n", ticker.Quotes[usd].Price)
	fmt.Fprintf(b, "*Rank*: %d\n", ticker.Rank)
	fmt.Fprintf(b, "*Sirculating supply:* %f\n", ticker.CirculatingSupply)
	fmt.Fprintf(b, "*MarketCap:* %f\n", ticker.Quotes[usd].MarketCap)
	fmt.Fprintf(b, "*24H change (perc.):* %f\n", ticker.Quotes[usd].PercentChange24H)
	fmt.Fprintf(b, "*7D change (perc.):* %f\n", ticker.Quotes[usd].PercentChange7D)
	return b.String(), nil
}
