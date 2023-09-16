package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Sign string

const XSign Sign = "X"
const OSign Sign = "O"

type GameApp struct {
	mainApp           fyne.App
	mainWindow        fyne.Window
	signLbl           *widget.Label
	buttons           [3][3]*widget.Button
	currentPlayerSign Sign
	opponentSign      Sign
	gameOver          bool
	conn              net.Conn
	isMyTurn          bool
}

func (a *GameApp) Initialize() {
	log.WithField("func", "Game.Initialize").Debug("Initializing...")
	a.mainApp = app.New()
	a.mainWindow = a.mainApp.NewWindow("Tic-Tac-Toe")

	// Set the initial window size to 400x400 pixels
	a.mainWindow.Resize(fyne.NewSize(400, 300))

	// Create a 3x3 grid of buttons for the game board
	var buttons [3][3]*widget.Button

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			buttons[i][j] = widget.NewButton("", a.handleButtonClick(false, i, j))
		}
	}
	// Create a grid layout to arrange the buttons
	grid := container.NewGridWithColumns(3,
		buttons[0][0], buttons[0][1], buttons[0][2],
		buttons[1][0], buttons[1][1], buttons[1][2],
		buttons[2][0], buttons[2][1], buttons[2][2],
	)

	signLbl := widget.NewLabel("")

	// Create a vertical box to hold the game board, status label, and buttons
	content := container.NewVBox(signLbl, grid)

	a.mainWindow.SetContent(content)

	a.signLbl = signLbl
	a.buttons = buttons
}

func (a *GameApp) Show() {
	log.WithField("func", "Game.Show").Debug("show...")
	a.mainWindow.ShowAndRun()
}

func (a *GameApp) SetTurn(turn Sign) {
	a.signLbl.SetText(string(turn))
}

func (a *GameApp) handleButtonClick(opponentMove bool, row, col int) func() {
	return func() {
		// Check if the button is empty and the game is not over
		if a.buttons[row][col].Text == "" && !a.gameOver {
			if !opponentMove && a.isMyTurn {
				move := Move{row, col, false, ""}
				if err := a.SendMoveData(move); err != nil {
					log.WithField("func", "Game.handleButtonClick").Errorln("error sending my move:", err)
					return
				}

				a.buttons[row][col].SetText(string(a.currentPlayerSign))
				a.isMyTurn = false
			} else if opponentMove && !a.isMyTurn {
				a.isMyTurn = true
				a.buttons[row][col].SetText(string(a.opponentSign))
			}

			// Check if the current player has won
			if a.IsCurrentPlayerWon() {
				move := Move{row, col, true, a.currentPlayerSign}
				if err := a.SendMoveData(move); err != nil {
					log.WithField("func", "Game.handleButtonClick").Errorln("error sending winner data:", err)
					return
				}
				a.gameOver = true
				a.showWinnerPopup(a.currentPlayerSign)
			}
		}
	}
}

// IsCurrentPlayerWon to check if the current player has won
func (a *GameApp) IsCurrentPlayerWon() bool {
	btns := a.buttons
	sign := string(a.currentPlayerSign)
	// Check rows
	for i := 0; i < 3; i++ {
		if btns[i][0].Text == sign && btns[i][1].Text == sign && btns[i][2].Text == sign {
			return true
		}
	}

	// Check columns
	for i := 0; i < 3; i++ {
		if btns[0][i].Text == sign && btns[1][i].Text == sign && btns[2][i].Text == sign {
			return true
		}
	}

	// Check diagonals
	if btns[0][0].Text == sign && btns[1][1].Text == sign && btns[2][2].Text == sign {
		return true
	}
	if btns[0][2].Text == sign && btns[1][1].Text == sign && btns[2][0].Text == sign {
		return true
	}

	return false
}

func (a *GameApp) showWinnerPopup(winner Sign) {
	popUpWindow := a.mainApp.NewWindow(fmt.Sprintf("Player %s Wins!", winner))

	content := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Player %s wins!", winner)),
		widget.NewButton("OK", func() {
			a.resetGame()
			popUpWindow.Close()
		}),
	)

	popUpWindow.SetContent(content)
	popUpWindow.Resize(fyne.NewSize(200, 100))
	popUpWindow.CenterOnScreen()
	popUpWindow.Show()
	a.resetGame()
}

func (a *GameApp) ShowPopUp(msg string) {
	myApp := fyne.CurrentApp()
	popUpWindow := myApp.NewWindow("Message")

	content := container.NewVBox(
		widget.NewLabel(msg),
		widget.NewButton("OK", func() {
			popUpWindow.Close()
		}),
	)

	popUpWindow.SetContent(content)
	popUpWindow.Resize(fyne.NewSize(200, 150))
	popUpWindow.CenterOnScreen()
	popUpWindow.Show()
}

// resetGame to reset the game by clearing the button labels
func (a *GameApp) resetGame() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			a.buttons[i][j].SetText("")
		}
	}
	a.gameOver = false
	if a.currentPlayerSign == XSign {
		a.isMyTurn = true
	} else {
		a.isMyTurn = false
	}
}

func (a *GameApp) SetSign(s Sign) {
	a.currentPlayerSign = s
	a.signLbl.SetText(fmt.Sprintf("Current Player: %s", string(a.currentPlayerSign)))
}

func (a *GameApp) SetOpSign(s Sign) {
	a.opponentSign = s
}

func (a *GameApp) StartServer(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s
		a.conn.Close()
		l.Close()
	}()

	log.WithField("func", "Game.StartServer").Infoln("Listening...")

	conn, err := l.Accept()
	a.conn = conn
	if err != nil {
		log.WithField("func", "Game.StartServer").Errorln("error accepting connecting:", err)
		return
	}

	if a.currentPlayerSign == XSign {
		a.isMyTurn = true
	}

	data, err := (&InitialData{OpponentSign: a.currentPlayerSign, GameStarter: XSign}).Marshal()
	if err != nil {
		log.WithField("func", "Game.StartServer").Errorln("initial data marshal error:", err)
		return
	}

	err = binary.Write(conn, binary.BigEndian, int64(len(data)))
	if err != nil {
		log.WithField("func", "Game.StartServer").Errorln("error writing data length to connection:", err)
	}
	conn.Write(data)

	a.Process()
}

func (a *GameApp) ConnectToServer(url string) {
	log.WithField("func", "Game.ConnectToServer").Debugln("connecting...")
	conn, err := net.Dial("tcp", url)
	a.conn = conn
	if err != nil {
		log.WithField("func", "Game.ConnectToServer").Errorln("error in dialing:", err)
		return
	}

	var size int64
	err = binary.Read(conn, binary.BigEndian, &size)
	if err != nil {
		log.WithField("func", "Game.ConnectToServer").Errorln("error reading initial data length:", err)
		return
	}

	buf := new(bytes.Buffer)
	_, err = io.CopyN(buf, conn, size)
	if err != nil {
		log.WithField("func", "Game.ConnectToServer").Errorln("error copying data to buffer:", err)
		return
	}

	initData := InitialData{}
	err = initData.Unmarshal(buf.Bytes())
	if err != nil {
		log.WithField("func", "Game.ConnectToServer").Errorln("error unmarshal init data: ", err)
		return
	}

	initData.ConfigureGame(a)

	a.Process()
}

func (a *GameApp) Process() {
	for {
		var size int64
		err := binary.Read(a.conn, binary.BigEndian, &size)
		if err != nil {
			if err == io.EOF {
				a.resetGame()
				a.ShowPopUp("Connection lost")
				a.gameOver = true
				a.conn.Close()
				return
			}
			log.WithField("func", "Game.Process").Errorln("error reading data length:", err)
			return
		}

		buf := new(bytes.Buffer)
		_, err = io.CopyN(buf, a.conn, size)
		if err != nil {
			log.WithField("func", "Game.Process").Errorln("error copying data to buffer:", err)
			return
		}

		move := Move{}
		err = move.Unmarshal(buf.Bytes())
		if err != nil {
			log.WithField("func", "Game.Process").Errorln("error marshalling move:", err)
			return
		}

		a.handleButtonClick(true, move.Row, move.Col)()
		if move.GameOver {
			a.gameOver = true
			a.showWinnerPopup(move.Winner)
		}
	}
}

func (a *GameApp) SendMoveData(move Move) error {
	data, err := move.Marshal()
	if err != nil {
		log.WithField("func", "Game.SendMoveData").Errorln("error marshalling move:", data)
		return err
	}

	err = binary.Write(a.conn, binary.BigEndian, int64(len(data)))
	if err != nil {
		log.WithField("func", "Game.SendMoveData").Errorln("error writing data length to connection:", data)
		return err
	}

	_, err = a.conn.Write(data)
	if err != nil {
		log.WithField("func", "Game.SendMoveData").Errorln("error writing data to connection:", data)
		return err
	}
	return nil
}
