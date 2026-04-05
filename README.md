# wasm-reversi

Go + WebAssemblyで実装したリバーシ（オセロ）ゲーム

## 概要

ゲームロジックをGoで実装し、WebAssemblyにコンパイルしてブラウザで動作させるリバーシゲームです。Go言語の`syscall/js`パッケージを使用して、JavaScriptとの相互運用を実現しています。

## 必要な環境

- Go 1.16以降
- Python 3（開発サーバー用）
- モダンなWebブラウザ（WebAssembly対応）

## セットアップ

```bash
# リポジトリをクローン
git clone <repository-url>
cd wasm-reversi

# WASMビルド
make build
```

## 使い方

### 開発サーバーの起動

```bash
make serve
```

ブラウザで `http://localhost:8080` にアクセスしてゲームをプレイできます。

### 利用可能なMakeコマンド

```bash
make              # ヘルプを表示
make build        # WASMファイルをビルド
make serve        # ビルドして開発サーバーを起動
make clean        # ビルド成果物を削除
make update-wasm-exec  # wasm_exec.jsを最新版に更新（Goバージョン変更時）
```

## 現在の状態

- ✅ 8x8の盤面表示
- ✅ クリックでコマを配置
- ✅ ターン切り替え
- ⚠️ **未実装**: 石をひっくり返すロジック

現在のバージョンではコマを置くことはできますが、リバーシのルールに従った石のひっくり返し処理は実装されていません。

## 技術スタック

- **Go**: ゲームロジック
- **WebAssembly**: ブラウザで実行可能なバイナリ形式
- **HTML/CSS/JavaScript**: UI
- **syscall/js**: Go ↔ JavaScript間の連携

## プロジェクト構成

```
.
├── main.go         # ゲームロジック（Go）
├── index.html      # UI
├── wasm_exec.js    # Go WASMランタイム
├── main.wasm       # コンパイル済みWASMバイナリ
├── Makefile        # ビルドタスク
└── README.md       # このファイル
```

## 開発

### コードを変更した場合

```bash
# 再ビルド
make build

# ブラウザをリロードして変更を確認
```

### Goバージョンを更新した場合

```bash
# wasm_exec.jsを最新版に更新
make update-wasm-exec

# 再ビルド
make serve
```

## ライセンス

MIT
