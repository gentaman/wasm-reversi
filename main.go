package main

import (
	"syscall/js"
)

var board [8][8]int // 0:空, 1:黒, 2:白
var turn = 1        // 1:黒から開始
var aiColor = 0     // AIの色（0:なし、1:黒、2:白）
var handicapMoves = 0 // ハンデ（AIが最初の何手をパスするか）
var moveCount = 0   // 手数カウント

func main() {
	c := make(chan struct{}, 0)

	// JSから呼び出せるように関数を登録
	js.Global().Set("pressCell", js.FuncOf(pressCell))
	js.Global().Set("getBoard", js.FuncOf(getBoard))
	js.Global().Set("getTurn", js.FuncOf(getTurn))
	js.Global().Set("getValidMoves", js.FuncOf(getValidMoves))
	js.Global().Set("isGameOver", js.FuncOf(isGameOverJS))
	js.Global().Set("getScore", js.FuncOf(getScoreJS))
	js.Global().Set("resetGame", js.FuncOf(resetGameJS))
	js.Global().Set("setGameConfig", js.FuncOf(setGameConfigJS))
	js.Global().Set("getAIMove", js.FuncOf(getAIMoveJS))
	js.Global().Set("isAITurn", js.FuncOf(isAITurnJS))

	initBoard()
	<-c // プログラムが終了しないように待機
}

func initBoard() {
	// 盤面をクリア
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			board[y][x] = 0
		}
	}
	// 初期配置
	board[3][3], board[4][4] = 2, 2
	board[3][4], board[4][3] = 1, 1
	// 黒から開始
	turn = 1
	moveCount = 0
}

func pressCell(this js.Value, args []js.Value) interface{} {
	x, y := args[0].Int(), args[1].Int()

	// 既に石がある、または置けない位置の場合は何もしない
	if board[y][x] != 0 || !canPlace(x, y, turn) {
		return nil
	}

	// 石を置く
	board[y][x] = turn

	// 石を反転
	flipStones(x, y, turn)

	// 手数カウント
	moveCount++

	// ターンを切り替え
	nextTurn := 3 - turn // 1→2、2→1

	// 次のプレイヤーが置ける場所がない場合はパス
	if !hasValidMove(nextTurn) {
		// 現在のプレイヤーも置けない場合はゲーム終了（ターンはそのまま）
		if !hasValidMove(turn) {
			return nil
		}
		// パス（ターンは変わらない）
	} else {
		turn = nextTurn
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

func getTurn(this js.Value, args []js.Value) interface{} {
	return turn
}

func getValidMoves(this js.Value, args []js.Value) interface{} {
	// 現在のターンで置ける場所を返す
	moves := []interface{}{}
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if canPlace(x, y, turn) {
				moves = append(moves, map[string]interface{}{
					"x": x,
					"y": y,
				})
			}
		}
	}
	return moves
}

// 8方向のベクトル
var directions = [][2]int{
	{-1, -1}, {0, -1}, {1, -1}, // 左上、上、右上
	{-1, 0}, {1, 0},            // 左、右
	{-1, 1}, {0, 1}, {1, 1},    // 左下、下、右下
}

// 指定位置に石を置けるかチェック
func canPlace(x, y, color int) bool {
	if board[y][x] != 0 {
		return false
	}

	opponent := 3 - color // 相手の色

	for _, dir := range directions {
		dx, dy := dir[0], dir[1]
		nx, ny := x+dx, y+dy
		foundOpponent := false

		// この方向に相手の石があるかチェック
		for nx >= 0 && nx < 8 && ny >= 0 && ny < 8 {
			if board[ny][nx] == opponent {
				foundOpponent = true
				nx += dx
				ny += dy
			} else if board[ny][nx] == color && foundOpponent {
				// 相手の石を挟んでいる
				return true
			} else {
				break
			}
		}
	}

	return false
}

// 石を反転
func flipStones(x, y, color int) {
	opponent := 3 - color

	for _, dir := range directions {
		dx, dy := dir[0], dir[1]
		nx, ny := x+dx, y+dy
		toFlip := [][]int{}

		// この方向に反転できる石を収集
		for nx >= 0 && nx < 8 && ny >= 0 && ny < 8 {
			if board[ny][nx] == opponent {
				toFlip = append(toFlip, []int{nx, ny})
				nx += dx
				ny += dy
			} else if board[ny][nx] == color && len(toFlip) > 0 {
				// 挟んでいるので反転
				for _, pos := range toFlip {
					board[pos[1]][pos[0]] = color
				}
				break
			} else {
				break
			}
		}
	}
}

// 指定した色が置ける場所があるかチェック
func hasValidMove(color int) bool {
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if canPlace(x, y, color) {
				return true
			}
		}
	}
	return false
}

// ゲームが終了したかチェック
func isGameOver() bool {
	// 両者とも置ける場所がない
	if !hasValidMove(1) && !hasValidMove(2) {
		return true
	}

	// どちらかの石がなくなった
	blackCount, whiteCount := 0, 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if board[y][x] == 1 {
				blackCount++
			} else if board[y][x] == 2 {
				whiteCount++
			}
		}
	}
	if blackCount == 0 || whiteCount == 0 {
		return true
	}

	// 盤面が全て埋まった
	if blackCount+whiteCount == 64 {
		return true
	}

	return false
}

// スコアを取得
func getScore() (int, int) {
	blackCount, whiteCount := 0, 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if board[y][x] == 1 {
				blackCount++
			} else if board[y][x] == 2 {
				whiteCount++
			}
		}
	}
	return blackCount, whiteCount
}

// JS用のラッパー関数
func isGameOverJS(this js.Value, args []js.Value) interface{} {
	return isGameOver()
}

func getScoreJS(this js.Value, args []js.Value) interface{} {
	black, white := getScore()
	return map[string]interface{}{
		"black": black,
		"white": white,
	}
}

func resetGameJS(this js.Value, args []js.Value) interface{} {
	initBoard()
	return nil
}

// ゲーム設定
func setGameConfigJS(this js.Value, args []js.Value) interface{} {
	// args[0]: aiColor (0: なし, 1: 黒, 2: 白)
	// args[1]: handicap (0〜5程度)
	aiColor = args[0].Int()
	handicapMoves = args[1].Int()
	return nil
}

// AIのターンかどうか
func isAITurnJS(this js.Value, args []js.Value) interface{} {
	return aiColor != 0 && turn == aiColor
}

// AIの思考（最も多く石を取れる手を選ぶ）
func getAIMoveJS(this js.Value, args []js.Value) interface{} {
	// ハンデ中はパス（手を打たない）
	if moveCount < handicapMoves && aiColor == turn {
		// ターン切り替え
		nextTurn := 3 - turn
		if hasValidMove(nextTurn) {
			turn = nextTurn
			moveCount++
		}
		return map[string]interface{}{
			"pass": true,
		}
	}

	bestX, bestY := -1, -1
	bestScore := -1

	// 角の位置に重みをつける
	cornerWeight := 100
	edgeWeight := 10

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if !canPlace(x, y, turn) {
				continue
			}

			// この手で取れる石の数をシミュレーション
			score := countFlips(x, y, turn)

			// 角の場合は優先
			if (x == 0 || x == 7) && (y == 0 || y == 7) {
				score += cornerWeight
			} else if x == 0 || x == 7 || y == 0 || y == 7 {
				// 辺の場合も少し優先
				score += edgeWeight
			}

			if score > bestScore {
				bestScore = score
				bestX = x
				bestY = y
			}
		}
	}

	if bestX == -1 {
		return map[string]interface{}{
			"pass": true,
		}
	}

	return map[string]interface{}{
		"x":    bestX,
		"y":    bestY,
		"pass": false,
	}
}

// 指定位置に置いた時に反転する石の数を数える
func countFlips(x, y, color int) int {
	opponent := 3 - color
	totalFlips := 0

	for _, dir := range directions {
		dx, dy := dir[0], dir[1]
		nx, ny := x+dx, y+dy
		flips := 0

		for nx >= 0 && nx < 8 && ny >= 0 && ny < 8 {
			if board[ny][nx] == opponent {
				flips++
				nx += dx
				ny += dy
			} else if board[ny][nx] == color && flips > 0 {
				totalFlips += flips
				break
			} else {
				break
			}
		}
	}

	return totalFlips
}
