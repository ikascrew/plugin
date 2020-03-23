
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

