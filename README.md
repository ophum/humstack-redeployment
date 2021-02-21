# humstack-redeployment

再展開用API&Agent

以下のようにJSONデータをAPIサーバーにPOSTするとgroup1グループ, test-nsネームスペースにあるtest-から始まるVMの再展開を行う
```
curl -X POST -H "Content-Type: application/json" http://localhost:8090/api/v0/redeployments -d @- <<EOF
{
    "meta": {
        "id": "unique-id"
    },
    "spec": {
        "group": "group1",
        "namespace": "test-ns",
        "vmIDPrefix": "test-",
        "restartTime": "2020-01-10T12:30:00+09:00"
    }
}
EOF
```

### Agentについて
- Agentは1ノードで動作させる
大まかに以下のような動作をする
1. 指定されたすべてのVMのActionStateをPowerOffにする
2. 指定されたすべてのVMに接続されているBSのステータスを`Error`に変更する
    - bsAgentは一度イメージを削除し作り直す
3. restartTimeの時間以降の場合に指定されたすべてのVMのActionStateをPowerOnにする

#### config.yaml
```
# redeployment apiserverのアドレスとポート
apiServerAddress: localhost
apiServerPort: 8090

# humstack apiserverのアドレスとポート
humstackAPIServerAddress: localhost
humstackAPIServerPort: 8080
```

### APIサーバーについて
以下のエンドポイントがある

| method| path | description |
| --- | --- | --- |
| GET | `/api/v0/redeployments` | redeploymentリソースの一覧を取得 |
| GET | `/api/v0/redeployments/:redeployment_id` | IDが`:redeployment_id`のredeploymentリソースを取得 |
| POST | `/api/v0/redeployments` | redeploymentリソースの作成 |
| PUT | `/api/v0/redeployments/:redeployment_id` | IDが`:redeployment_id`のredeploymentリソースを更新 |
| DELETE | `/api/v0/redeployments/:redeployment_id` | IDが`:redeployment_id`のredeploymentリソースを削除