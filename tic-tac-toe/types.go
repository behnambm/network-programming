package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

// InitialData will be received by the peer that connected to a server
type InitialData struct {
	OpponentSign Sign
	GameStarter  Sign
}

func (i *InitialData) Marshal() ([]byte, error) {
	data, err := json.Marshal(i)
	if err != nil {
		log.WithField("func", "InitialData.Marshal").Debugln("error marshaling json:", err)
		return nil, err
	}
	return data, nil
}

func (i *InitialData) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, i)
	log.WithField("func", "InitialData.Unmarshal").Debugln("error unmarshalling json:", err)
	return err
}

func (i *InitialData) ConfigureGame(game *GameApp) {
	if i.OpponentSign == OSign {
		game.SetSign(XSign)
	} else if i.OpponentSign == XSign {
		game.SetSign(OSign)
	}
	game.SetOpSign(i.OpponentSign)

	if i.GameStarter == game.currentPlayerSign {
		game.isMyTurn = true
	}
}

// Move represents all user actions and their selection of position in board
type Move struct {
	Row      int
	Col      int
	GameOver bool
	Winner   Sign
}

func (m *Move) Marshal() ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		log.WithField("func", "Move.Marshal").Debugln("error marshalling json:", err)
		return nil, err
	}
	return data, nil
}

func (m *Move) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, m)
	log.WithField("func", "Move.Unmarshal").Debugln("error unmarshalling json:", err)
	return err
}
