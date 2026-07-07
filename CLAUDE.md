# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

ikascrew VJ システムの映像プラグイン集。コンテンツの「型」と「JSON パラメータ」の**本家リポジトリ**であり、型語彙・param 形式のルールはすべてここで決まる。server / ikasbox は `video` レジストリ経由でこのリポジトリを使う(client は import しない — gocv 依存を持ち込まないため)。

ビルドには OpenCV(gocv)が必要。テストは `video/terminal` に1ファイルのみ。

## video レジストリ(video/video.go)— 型解決の唯一の入口

- `video.Get(t, param string) (core.Video, error)` — 型名から実装を生成する**唯一の対応表**。server の Effect、ikasbox の生成型コンテンツ登録の両方がこれを通る。型ごとの分岐を利用側に書かないこと。
- 正語彙は **`file` / `img` / `cd` / `terminal`**(`video.Types()`)。`video.Normalize` が旧語彙を吸収する(`image`→`img`、`countdown`→`cd`、空→`file`)。正規化はこの関数だけに置く。
- `video.IsGenerative(t)` — 実体ファイル不要の生成型(`cd`, `terminal`)かの判定。ikasbox の登録方法の分岐用。

## param(JSON)の規約 — 全体ルール

1. param は **JSON 文字列**。ikasbox の DB(`contents.params`)→ `/project/content/list` API → server の work file → gRPC 経由の各層は**不透明な文字列として素通し**する(中間層でパースしない)。
2. **解釈は各プラグインの `New(param string)` だけ**が行う。各プラグインは自分の `Params` 構造体 + `parseParams` を持つ(構造体の所有権はプラグイン側)。
3. **非 JSON(`{` で始まらない)文字列は旧来の生文字列**として解釈する(file/img はパス、terminal/cd はテキスト)。server の旧 work file や起動スプラッシュ(`Get("terminal", 生テキスト)`)がこれに依存しているため、このフォールバックは削らないこと。
4. 検証は登録時(ikasbox `content register`)に `video.Get` を実際に呼んで行う。フィールド追加は encoding/json が未知フィールドを無視するため後方互換。

## 各プラグインの param スキーマ

| 型 | パッケージ | JSON param | 非 JSON 時の解釈 |
|---|---|---|---|
| `file` | `video/file` | `{"path": "動画ファイルパス"}` | 全体をパスとして扱う |
| `img` | `video/image` | `{"path": "画像ファイルパス"}` | 全体をパスとして扱う |
| `cd` | `video/countdown` | `{"target": "RFC3339 または 2006-01-02 15:04:05(JST)", "text": "終了後に表示する文字列"}`。target 省略時は過去日時(即終了表示) | 全体を終了後テキストとして扱う |
| `terminal` | `video/terminal` | `{"text": "表示テキスト(\n 区切り)"}` | 全体を表示テキストとして扱う |

`video/telop` は agouti(ブラウザ自動化)依存の実験実装で、レジストリには**登録していない**(server/ikasbox に agouti を持ち込まないため)。

## 新しい video プラグインの追加手順

1. `video/<name>/` に `New(param string) (core.Video, error)` を実装(`Params` 構造体 + `parseParams` は既存プラグインの形を踏襲)。`core.Video` インターフェースは `Next/Wait/Set/Current/Source/Release`。
2. `video/video.go` の `Get` の switch と `Types()` に型名を追加(必要なら `Normalize` に別名、`IsGenerative` に生成型判定も)。
3. これだけで ikasbox の `content register` と server の再生が両方使えるようになる(中間層の変更は不要 — それがこの設計の眼目)。

## WIP(ビルド対象外のスケッチ)

`effect/`(light, speed)・`transition/switch`・`cmd/`(plugins.json → plugins_gen.go のレジストリ生成、plugin.so 動的ロード)・`cap.go` は設計スケッチで、`//go:build ignore` によりビルド対象外(`go run` での個別実行は可)。動的ロード(root の `plugin.go`, Go plugin 機構)は Linux 限定なので、Windows 開発中は静的レジストリ(video/video.go)が本線。effect/transition の設計は core リポジトリの `core.Effect` / `core.Transition` インターフェース案と対になっている。

## Conventions

- エラーは `golang.org/x/xerrors` でラップ。コミット件名は `fix:` / `feat:` / `wip:` プレフィックス(日本語)。
- 関連リポジトリ(core / server / ikasbox)からは開発中 `replace` でローカル参照されている。ここを変更したら push + 各リポジトリの pseudo-version 更新(または replace 維持)が必要。
