# k8s-observability-app

## プロジェクト概要
本プロジェクトは、Cloud Native技術とObservabilityの実践的な習得を目的とした個人開発ポートフォリオです。以下のスキルセットを実践的な内容で向上できるように設計、構築しました。
- Observability：Prometheus Clientを用いたメトリクス収集、Lokiによるログ収集、Tempoによるトレース収集、Grafanaによる可視化、Alertmanagerによるアラート通知までの一貫した可観測性の実装
- Go：gRPCサーバの開発とPrometheus Exporterによるメトリクス出力
- OpenTelemetry：Logs(OpenTelemetry collectorのk8slogs receiver経由でLoki送信)とTraces(OTLP Exporter→Tempo)を利用
- Kubernetes：Helmチャートを用いたNamespace分離、Deployment、Service、ConfigMapのリソース管理
- Cloud Native：EKSを利用した本番環境想定のデプロイ
- CI/CD：GitHub ActionsによるCIとArgoCDによるEKSへの自動デプロイ、GitOps運用の実践
- IaC：HelmとKubernetesマニフェスト、ArgoCDによるIaC管理

## 使用技術
| 分類                         | 使用技術                                                                                                           |
| -------------------------- | -------------------------------------------------------------------------------------------------------------- |
| **言語・フレームワーク**             | Go (gRPCサーバ、Prometheus Exporter、zap logger)                                                                    |
| **可観測性**                   | Prometheus (Metrics)、Loki (Logs)、Tempo (Traces)、Grafana (ダッシュボード可視化)、Alertmanager (通知)、OpenTelemetry Collector |
| **CI/CD・GitOps**           | GitHub Actions (CI)、ArgoCD (CD、GitOpsフロー)                                                                      |
| **Kubernetesリソース**         | Helmチャート、Deployment、Service、ConfigMap、Namespace分離、HorizontalPodAutoscaler (HPA)                                |
| **クラウドプラットフォーム**           | Amazon EKS (Elastic Kubernetes Service)                                                                        |
| **Infrastructure as Code** | Helm、ArgoCD ApplicationマニフェストによるIaC管理                                                                          |

## ディレクトリ構成
```
.
├── cmd                  # Goアプリケーション(gRPCサーバ、Prometheus Exporter)
│   └── grpc-server
├── docs                 # Grafanaダッシュボード定義
│   └── dashboards
├── gen                  # Protobufコンパイル生成ファイル
│   └── proto
├── go.mod               # Go Module定義ファイル
├── go.sum               # Go Module依存情報
├── helm                 # Helmチャート(Kubernetes用のマニフェスト管理)
│   ├── grpc-app         # gRPCアプリ用チャート
│   └── monitoring-stack # Monitoring Stack用チャート(Prometheus,Grafana,Tempo,Loki,Collector)
├── manifests            # ArgoCD Application,Namespace,HPAなどKubernetesマニフェスト
├── proto                # gRPC通信定義
│   └── echo.proto
├── README.md            # プロジェクト概要
```

## EKSデプロイ手順

EKSへのデプロイは以下の手順で構成

### 1_前提条件
- AWSアカウント、IAMユーザ作成済み
- eksctl, kubectl, helm, argocd CLIインストール済み
- GitHubリポジトリのclone済み

### 2_EKSクラスタ作成
```bash
eksctl create cluster \
  --name k8s-observability-app \
  --region ap-northeast-1 \
  --nodes 2 \
  --node-type t3.medium
```
### 3_ArgoCDのインストール
```bash
kubectl create namespace argocd
kubectl apply -n argocd -f manifests/argocd-install.yaml
```

ArgoCD UIへのアクセス例（port-forward経由）：
```bash
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

### 4_ArgoCD Application登録
本リポジトリの manifests/argocd-app-grpc-eks.yaml と argocd-app-monitoring-eks.yaml を適用します：
```bash
kubectl apply -f manifests/argocd-app-grpc-eks.yaml
kubectl apply -f manifests/argocd-app-monitoring-eks.yaml
```
GitHub mainブランチ更新でEKSへ自動デプロイされます。

### 5_Grafanaアクセス確認
Grafana（LoadBalancer経由）へのアクセス例：
```bash
kubectl get svc -n monitoring
```
Grafana初期ログイン：admin / admin（helm/monitoring-stack/values.yaml にて管理）

### 6_Slack通知確認
AlertmanagerからSlack Webhookへの通知も確認できます（helm/monitoring-stack/values.yaml に設定済み）。

