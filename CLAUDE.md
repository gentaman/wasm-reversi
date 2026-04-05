# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

Go + WebAssemblyで実装されたリバーシ（オセロ）ゲーム。ゲームロジックはGoで実装し、WebAssemblyにコンパイルしてブラウザで実行する。

## アーキテクチャ

### Go ↔ JavaScript 連携

- `main.go`でGoの関数を`js.FuncOf`を使ってJavaScriptに公開
- `js.Global().Set()`でグローバルなJavaScript関数として登録
- JavaScript側から直接Go関数を呼び出せる

**公開されている関数**:
- `pressCell(x, y)`: セルがクリックされた時の処理。盤面を更新してターンを切り替え
- `getBoard()`: 現在の盤面状態（8x8の二次元配列）を取得

### ゲームロジック

- `board [8][8]int`: 盤面の状態（0:空、1:黒、2:白）
- `turn`: 現在のプレイヤー（1:黒、2:白）
- `aiColor`: AIの色（0:なし、1:黒、2:白）
- `handicapMoves`: ハンデ（AIが最初の何手をパスするか）
- 初期配置は`initBoard()`で中央4マスにセット
- 石の反転ロジック、終局判定、AI思考が実装済み

### AIアルゴリズム

- 評価基準：反転できる石の数 + 位置ボーナス
  - 角：+100点（最重要）
  - 辺：+10点
- ハンデ機能：序盤の指定手数をパス

## 開発コマンド

### Makefile経由（推奨）

```bash
make serve            # ビルド + 開発サーバー起動（http://localhost:8080）
make build            # WASMファイルのみビルド
make clean            # ビルド成果物を削除
make update-wasm-exec # wasm_exec.jsを最新版に更新
make help             # ヘルプを表示
```

### 直接実行

```bash
# ビルド
GOOS=js GOARCH=wasm go build -o main.wasm main.go

# 開発サーバー起動
python3 -m http.server 8080

# wasm_exec.jsの更新（Goバージョン変更時）
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

## ファイル構成

- `main.go`: ゲームロジック（Go + syscall/js）
- `index.html`: UI（HTML + JavaScript）
- `wasm_exec.js`: Go WASMランタイム（Goツールチェーンから提供）
- `main.wasm`: コンパイル済みWASMバイナリ
