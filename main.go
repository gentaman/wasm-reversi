package main

import (
	"syscall/js"
)

var board [8][8]int // 0:空, 1:黒, 2:白
var turn = 1        // 1:黒から開始

func main() {
	c := make(chan struct{}, 0)

	// JSから呼び出せるように関数を登録
	js.Global().Set("pressCell", js.FuncOf(pressCell))
	js.Global().Set("getBoard", js.FuncOf(getBoard))

	initBoard()
	<-c // プログラムが終了しないように待機
}

func initBoard() {
	board[3][3], board[4][4] = 2, 2
	board[3][4], board[4][3] = 1, 1
}

func pressCell(this js.Value, args []js.Value) interface{} {
	x, y := args[0].Int(), args[1].Int()
	
	if board[y][x] == 0 {
		board[y][x] = turn
		// 本来はここで「ひっくり返すロジック」を実装
		if turn == 1 { turn = 2 } else { turn = 1 }
	}
	return nil
}

func getBoard(this js.Value, args []js.Value) interface{} {
	// JSに現在の盤面を二次元配列として返す
	res := make([]interface{}, 8)
	for y := 0; y < 8; y++ {
		row := make([]interface{}, 8)
		for x := 0; x < 8; x++ {
			row[x] = board[y][x]
		}
		res[y] = row
	}
	return res
}
