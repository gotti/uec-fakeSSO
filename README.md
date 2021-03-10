# これは？
いろいろなサービスにUECメールだけでパスワードレスログインできるようにするための認証サービスです．
ただし今のところサービスは1つしか登録できません．

# TODO
現在以下の実装ができていません．
- UTF-8のメール送信をきちんとやる(UTF-8を解釈できないと文字化け)
- ポート番号8084がハードコーディングされているのでなんとかしたい
- 複数のサービスで同じAPIを使えるようにしたいな

# Creating API Server
- cp config.toml.sample config.toml
- config.tomlを環境に応じて編集してください．

APIServerTokenは正当な利用者以外からアクセスできないようにするためのものです(SlackとかDiscordのWeb Hookのようなものです)．十分に長い英数字を書き込んでこの情報は公開しないでください．悪意あるユーザーが任意のUECメールに迷惑メールを送信できる脆弱性になります．英数字(大文字可)のみを保証し，記号は正常な動作を保証しません．

利用する際にこのAPIServerTokenをWeb HookのappTokenとして指定してください．

OTTExpireにはワンタイムトークンの有効期間を書き込んでください．単位は分です．

UserListDatabaseにはsqlite3データベースの場所を書き込んでください．特に問題がなければ初期設定のままで構いません．

UserListFileにはユーザー名を改行区切りで記入したファイルを指定してください．特に問題がなければ初期設定の./usersを利用してください．

[Mail]以下について

利用しているメールサーバのSMTPのログイン情報を書き込んでください．
  - smtpAddress: SMTPサーバのアドレス
  - port: ポート番号
  - from: 送信元として表示されるメールアドレス
  - username: ユーザー名
  - password: パスワード
  - sub: メールの題名
  - msg: メール本文

# How to use this API
- このAPIは共通して正常であれば200を返し異常であれば400,500番台を返します．internal/handler.goを雰囲気で読むとどのレスポンスコードが返ってくるかわかります．

- このAPIはregisterとverifyという2つの処理を提供します．
  - registerはオプションとしてappToken,usernameをとり，指定されたUECアカウントに生成したワンタイムトークンを送信します．以下の3つの条件が揃えばメールを送信して200を返答します．そうでなければ401を返答します．
    - appTokenがAPIServerTokenと一致していること
    - UECアカウントが存在すること
    - そのサービスに既に登録されていないこと
    - 有効なトークンが残っていないこと
  - verifyはオプションとしてappToken,username,ottをとり，指定されたUECアカウントに正当なワンタイムトークン(ott)が発行されていることを確認します．以下の条件が揃えば200を返答し登録済みとして記録しregisterを無効化し，そのワンタイムトークンを無効化します．そうでなければ401を返答します．
    - appTokenがAPIServerTokenと一致していること
    - UECアカウントが存在すること
    - ワンタイムトークンが正当であること

- HTTPのGETが送信できてレスポンスのステータスコードが見れる言語であり，かつappTokenをユーザに漏洩させない言語であれば任意の言語が使えます．htmlとjavascriptだけではソースからappTokenが見えるので，必ずサーバサイドで実装してください(PHPに一度POSTするなど)． 以下にGolangでの参考実装を示します．
```go:register.go
var appToken, username string //適切に初期化されているとします．appTokenは自分で指定するので信用できますが，usernameはユーザ入力なので信用できません．
reg := regexp.MustCompile(`[a-z]\d{7}`) //ここまで検証しなくてもいいのですが，記号は送信しないようにしてください．(URLがこわれるかも)
if !reg.MatchString(username) {
    fmt.Println("不適切なユーザ名です")
    return
}
url := fmt.Sprintf("https://example.com/api/register?appToken=%s&username=%s", appToken, username)
client := &http.Client{}
responce, _ := client.Get(url)
if responce.Code != 200{
    //異常終了，ユーザにエラーを表示
} else {
    //正常終了
}
```

```go:verify.go
var appToken, username, ott string //適切に初期化されているとします．appTokenは自分で指定するので信用できますが，usernameとottはユーザ入力なので信用できません．
regUser := regexp.MustCompile(`[a-z]\d{7}`) //ここまで検証しなくてもいいのですが，記号は送信しないようにしてください．(URLがこわれるかも)
if !reg.MatchString(username) {
    fmt.Println("不適切なユーザ名です")
    return
}
regUser := regexp.MustCompile(`[0-9a-f]{6}`) //ここまで検証しなくてもいいのですが，記号は送信しないようにしてください．(URLがこわれるかも)
if !reg.MatchString(ott) {
    fmt.Println("不適切なワンタイムトークンです")
    return
}
url := fmt.Sprintf("https://example.com/api/verify?appToken=%s&username=%s&ott=%s", appToken, username, ott)
client := &http.Client{}
responce, _ := client.Get(url)
if responce.Code != 200{
    //異常，ユーザ登録をしないでください
    return
} else {
    //正常，ユーザ登録をしてください
}
```
