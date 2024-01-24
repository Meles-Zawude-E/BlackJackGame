package handler

import (
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/meles-zawude-e/shuffleGame/database"
	"github.com/meles-zawude-e/shuffleGame/model"
	"gorm.io/gorm"
)

func CreateDeck(c echo.Context) error {
	var cardsParam struct {
		Shuffled bool   `json:"shuffled"`
		Cards    string `json:"cards"`
	}

	if err := c.Bind(&cardsParam); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	var cards []model.Card

	if cardsParam.Cards != "" {
		cards = parseCards(cardsParam.Cards)
	} else {
		cards = createStandardDeck()
	}

	if cardsParam.Shuffled {
		shuffleDeck(cards)
	}

	deckID := uuid.New()
	deck := model.Deck{
		ID:        deckID,
		Shuffled:  cardsParam.Shuffled,
		Remaining: len(cards),
	}

	db := database.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&deck).Error; err != nil {
			return err
		}

		for _, card := range cards {
			card.DeckID = deckID
			if err := tx.Create(&card).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save deck to database"})
	}

	return c.JSON(http.StatusOK, deck)
}

func OpenDeck(c echo.Context) error {
	deckID, err := uuid.Parse(c.Param("deckID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid deck ID"})
	}

	db := database.GetDB()
	var deck model.Deck
	if err := db.Preload("Cards").First(&deck, deckID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Deck not found"})
	}

	return c.JSON(http.StatusOK, deck)
}

func DrawCard(c echo.Context) error {
	deckID, err := uuid.Parse(c.Param("deckID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid deck ID"})
	}

	db := database.GetDB()
	var deck model.Deck
	if err := db.Preload("Cards").First(&deck, deckID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Deck not found"})
	}

	countParam := strings.TrimSpace(c.QueryParam("count"))
	count, err := strconv.Atoi(countParam)

	if err != nil || count <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid count parameter"})
	}

	if count > deck.Remaining {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Not enough cards remaining in the deck"})
	}

	drawnCards := deck.Cards[len(deck.Cards)-count:]
	deck.Cards = deck.Cards[:len(deck.Cards)-count]
	deck.Remaining -= count

	// Update the deck in the database
	if err := db.Save(&deck).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update deck in the database"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"cards": drawnCards})
}

func parseCards(cardsParam string) []model.Card {
	var cards []model.Card
	cardCodes := strings.Split(cardsParam, ",")

	for _, code := range cardCodes {
		value, suit := parseCardCode(code)
		card := model.Card{
			Value: value,
			Suit:  suit,
			Code:  code,
		}
		cards = append(cards, card)
	}

	return cards
}

func createStandardDeck() []model.Card {
	var cards []model.Card
	suits := []string{"SPADES", "DIAMONDS", "CLUBS", "HEARTS"}
	values := []string{"ACE", "2", "3", "4", "5", "6", "7", "8", "9", "10", "JACK", "QUEEN", "KING"}

	for _, value := range values {
		for _, suit := range suits {
			code := value[:1] + strings.ToUpper(suit[:1])
			card := model.Card{
				Value: value,
				Suit:  suit,
				Code:  code,
			}
			cards = append(cards, card)
		}
	}

	return cards
}

func shuffleDeck(cards []model.Card) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
}

func parseCardCode(code string) (value, suit string) {
	if len(code) < 2 {
		return "", ""
	}
	return code[:len(code)-1], code[len(code)-1:]
}
