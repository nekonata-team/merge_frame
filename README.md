# merge_frame

スクリーンショット画像をフレーム画像に合成し、指定したディレクトリに出力するツールです。

test_dataディレクトリにサンプルのスクリーンショット画像とフレーム画像が含まれています。

## 使い方

コマンドライン引数で各パスを指定します。

```bash
go run main.go --screenshots=test_data/screenshots/ --frame=test_data/internal/frame.png --output=test_data/output
```

- `--screenshots` : スクリーンショット画像が入ったフォルダ
- `--frame` : 合成するフレーム画像のパス
- `--output` : 出力先ディレクトリ

## ビルド手順

1. 依存パッケージをインストール

    ```bash
    go get github.com/fogleman/gg
    go get github.com/urfave/cli/v2
    ```

2. ビルド

    ```bash
    go build -o merge_frame main.go
    ```

3. 実行例

    ```bash
    ./merge_frame --screenshots=test_data/screenshots/ --frame=test_data/internal/frame.png --output=test_data/output
    ```

## 備考

- PNG画像のみ対応しています。
- 出力先ディレクトリは自動作成されます。
