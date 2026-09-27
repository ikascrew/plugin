
# plugins.json

ikasbox上でjsonを編集すると

go build --buildmode=plugin -o plugins.so plugins_gen.go

でplugins.soを作成して、サーバ上のプロセスを更新する

- サーバがlinuxでないといけなくなる。
- Windows上での動作確認ができなくなる？

GOOS=linuxで行うと、gocvがエラーになる。

https://github.com/hybridgroup/gocv/issues/615

なので同一OSでのビルドになる

linuxでビルドした後にWindows上でビルドできるかを試す

jsonを元にWindows用のDLLやPluginを作成して読み込めるようにする


プラグインは現状Videoを取得することしかできてないが

DB上での設定画面
コマの取得

# TODO

## countdown / terminal の解像度追従描画

2026-07-26 の作業ツリー消失で失われ、まだ作り直していない。

- `video/countdown/countdown.go` と `video/terminal/terminal.go` の `gocv.PutText` は、座標(`image.Pt(60, 400)` など)と文字サイズが固定値のまま。出力解像度を変えると表示がはみ出したり寄ったりする
- 消失前は中央寄せの `putCenterText` と、フレームサイズから求める `unit` / `thickness` で文字サイズ・線幅をスケーリングしていた

## 再構成したコードの見直し

消失後に作り直したため、消失前の実装と細部が異なる可能性がある。

- `video/output` — 呼び出し側から API(`Set` / `Size` / `DefaultWidth` / `DefaultHeight`)を逆算して書き直した
- countdown / terminal の `*core.Frame` 対応 — 描画ロジックは 2026-07-07 版のまま機械的に適合させた

