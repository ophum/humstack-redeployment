# humstack-redeployment

再展開用API&Agent

以下のようにJSONデータをAPIサーバーにPOSTするとgroup1グループ, test-nsネームスペースにあるtest-から始まるVMの再展開を行う
```
curl -X POST -H "Content-Type: application/json" http://localhost:8090/api/v0 -d @- <<EOF
{
    "meta": {
        "id": "unique-id"
    },
    "spec": {
        "group": "group1",
        "namespace": "test-ns",
        "vmIDPrefix: "test-",
        "restartTime": "2020-01-10T12:30:00+09:00"
    }
}
EOF
```

### Agentについて
- Agentは各Computeノード上で動作させる
    - group/namespace/vmIDPrefixで取得したVMのリストのうち最初の1つのAnnotationを見て自ノード上のVMかどうかを判別する
- humstackが保存しているBSのファイルを直接消しに行くので権限周り注意する
    - root
大まかに以下のような動作をする
1. 指定されたすべてのVMのActionStateをPowerOffにする
2. 指定されたすべてのVMに接続するBSを削除する
3. restartTimeの時間以降の場合に指定されたすべてのVMのActionStateをPowerOnにする

#### config.yaml
```
# redeployment apiserverのアドレスとポート
apiServerAddress: localhost
apiServerPort: 8090

# humstack apiserverのアドレスとポート
humstackAPIServerAddress: localhost
humstackAPIServerPort: 8080
# humstackのbsが保存されるベースパス
humstackBlockStorageDirPath: /var/lib/humstack/blockstorages

# bsを削除するときの並行数
bsDeleteParallelLimit: 10
```

### APIサーバーについて
以下のエンドポイントがある

| method| path | description |
| --- | --- | --- |
| GET | `/api/v0` | redeploymentリソースの一覧を取得 |
| GET | `/api/v0/:redeployment_id` | IDが`:redeployment_id`のredeploymentリソースを取得 |
| POST | `/api/v0` | redeploymentリソースの作成 |
| PUT | `/api/v0/:redeployment_id` | IDが`:redeployment_id`のredeploymentリソースを更新 |
| DELETE | `/api/v0/:redeployment_id` | IDが`:redeployment_id`のredeploymentリソースを削除