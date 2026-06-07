# DevOps Monorepo 架構設計

本文件收斂整個架構討論的最終決策,作為新專案的設計基準,同時作為 **Claude Code 的實作 context** 與**團隊 onboarding 文件**。內容為規範性(normative):列出的原則與鐵則視為實作約束,不是建議。

---

## 一、核心設計原則

這幾條原則驅動下面所有結構決策,順序就是優先順序。

1. **依「部署單元」組織,不依技術類型**。用最穩定的軸切 repo;技術層(backend/frontend)是最易變的軸,不拿來當目錄結構。
2. **巢狀只代表「版本生命週期綁定」**。要不要巢狀,只看「是否必須一起切版、一起上線」,跟打包成不成同一顆 image 無關。
3. **發版與部署是兩件事**。release 負責造出帶版本號的不可變 artifact;deploy 負責選哪個版本上哪個環境。
4. **共用 code 用 live at head**。monorepo 內部 lib 直接引用原始碼,不發版,改了立即生效。
5. **避免分散式單體**。只共用真正穩定、跨切面的東西;不同關注點不硬綁。
6. **避開單行道工具**。建構工具起步從簡,被現實逼到再升級。
7. **image 不存在前,manifest 絕不指向它**。任何更新 image tag 的動作都必須在 image 已 build + push 到 registry 之後,避免部署競態(ImagePullBackOff)。

---

## 二、目錄結構

```
monorepo/
├── apps/                          # 可獨立部署單元(有 image、會部署)
│   ├── service-a/                 # Go 後端
│   │   ├── README.md             # app 自己的說明(本地語境)
│   │   └── ...
│   ├── service-b/                 # Python 後端
│   ├── web-console/               # 前端
│   ├── webhook-agent/             # agent 服務也是 app;v1 單一進程,需要時再拆
│   └── booking/                   # 前後端綁定上線才巢狀
│       ├── api/
│       └── web/
├── libs/                          # 被引用但不部署(live at head)
│   ├── go/common/
│   ├── python/common/
│   └── proto/                     # 跨語言合約 → 各語言產生 client
├── deploy/                        # 部署描述(與 app 解耦)
│   ├── charts/                    # Helm chart 本體(怎麼部署)
│   │   ├── service-a/
│   │   └── web-console/
│   └── envs/                      # 部署哪個版本到哪個環境(同一顆 image 跨環境)
│       ├── dev/                   # 可選:開發自動部署
│       ├── stage/                 # 驗證
│       └── prd/                   # 正式
├── infra/                         # 純 IaC(provision 基礎設施)
│   ├── terraform/
│   ├── ansible/
│   └── modules/
├── ci/                            # 非 GHA 的 CI 資產(若決定納入,見第九章)
│   └── jenkins/
│       ├── shared-library/        # Jenkins Shared Library(vars/, src/)
│       └── jobs/                  # Jenkinsfile / job 定義
├── tools/                         # 內部 script / codegen / CLI
├── docs/                          # 跨切面文件
│   ├── adr/                       # 架構決策紀錄(本架構文件放這)
│   ├── guides/                    # 教學 / 操作指南
│   └── runbooks/                  # 維運手冊
├── .github/
│   ├── workflows/                 # GHA:ci / release / build-image / deploy
│   ├── CODEOWNERS                 # 擁有權邊界
│   ├── dependabot.yml
│   ├── pull_request_template.md
│   └── ISSUE_TEMPLATE/
├── README.md                      # repo 入口
├── CONTRIBUTING.md
├── SECURITY.md
├── .editorconfig
├── .gitignore
├── .gitattributes
├── Taskfile.yml
├── commitlint.config.js
├── release-please-config.json
└── .release-please-manifest.json
```

---

## 三、apps / libs / 巢狀 規則

| 概念 | 定義 | 判準 |
|------|------|------|
| apps | 會被部署的單元(有 image、有部署目標) | 「它以獨立單元部署嗎?」是 → apps |
| libs | 被 import 但不部署的程式碼 | 沒有 image、沒有部署目標 → libs |
| 平輩 | 各自切版、各自決定上線 | 獨立部署 |
| 巢狀 | 共用版本號、一起切版上線 | 綁定上線(image 仍可各自一顆) |

agent 服務屬於 apps:它收 webhook、有版本、有 image、會部署,形狀跟一般 service 相同,只是內部用了 LLM。

---

## 四、共用 library

內部共用 code 放 `libs/`,按語言分,用 live at head 解析相依。

| 方式 | 機制 | 用途 |
|------|------|------|
| Live at head | Go `go.work` / Python uv workspace 直接引用原始碼 | monorepo 內部共用(預設) |
| 發布版本 | lib 發成版本化套件,consumer pin 版本 | 給外部 repo 用 |

重點:

- internal lib **不需要自己的版本和 release**,因為不出貨、consumer 也不 pin。
- 改 lib 的把關靠 **CI affected 偵測**,而非版本號:改了 `libs/go/common` 要連帶重測所有 consumer。
- **跨語言共用 = 共用合約,不是共用 code**。Go 與 Python 要共用資料模型時,放 `libs/proto/` 的 schema,各語言各自產生 client,這組用 linked-versions 綁版本。

---

## 五、版本與發布策略

採用 release-please **manifest mode**,獨立版本為主,強耦合群組才 linked-versions。

```json
{
  "$schema": "https://raw.githubusercontent.com/googleapis/release-please/main/schemas/config.json",
  "separate-pull-requests": true,
  "packages": {
    "apps/service-a": { "release-type": "go", "component": "service-a" },
    "apps/service-b": { "release-type": "python", "component": "service-b" },
    "apps/web-console": { "release-type": "node", "component": "web-console" },
    "apps/webhook-agent": { "release-type": "python", "component": "webhook-agent" },
    "infra/terraform": { "release-type": "terraform-module", "component": "infra" }
  },
  "plugins": [
    {
      "type": "linked-versions",
      "groupName": "proto-bundle",
      "components": ["proto", "go-client", "py-client"]
    }
  ]
}
```

前提:**conventional commits 必須落實**。release-please 完全靠 commit message 判斷 bump,所以要搭 commitlint + PR title 檢查,否則版本會算錯。internal lib(如 `libs/go/common`)不放進 packages,或只給 `simple` type 純追 changelog。Helm chart 起步**不列為 package**(不獨立管 chart 版本);等到要把 chart 發佈到 chart registry 給外部用,再加入 `"release-type": "helm"`。

### 版本流:從 commit 到部署

版本只被 release-please 算一次,然後往下游流,沿途的 image tag 都是衍生值,不重新產生。

```
conventional commit
  → release-please 算出 service-a 1.2.3,打 tag service-a-v1.2.3
    → build workflow 建 image: registry/service-a:1.2.3
      → deploy 把 env values 的 image.tag 設成 1.2.3
        → helm upgrade / kubectl apply 用這個 tag
```

### 哪個欄位才真正控制部署版本

| 欄位 | 位置 | 控制部署版本? | 誰更新 |
|------|------|--------------|--------|
| `image.tag` | `deploy/envs/<env>/<app>.yaml` | **是,這個才算數** | 部署環節(非 release-please) |
| Chart.yaml `appVersion` | chart | 否,只是 metadata | 可選 |
| Chart.yaml `version` | chart | 否 | 起步先別管 |
| pyproject / package.json version | app 原始碼 | 否 | release-please |

實際決定叢集跑哪顆 image 的是 values 裡的 `image.tag`,所以**要同步的就是它**;Chart.yaml 的版本不影響實際部署。

### image tag 同步機制(關鍵)

release-please **只更新原始碼裡的版本檔並打 tag,不會去動 `deploy/envs/` 的 image tag**。image tag 同步是部署環節的事,兩者分開。

強理由:stage 與 prd 是在不同時間點拿到同一版本(stage 先上,prd 核可後才上)。release-please 一次 bump 只發生一次,無法表達 per-env 的版本差異,所以 per-env image tag 一定由部署環節同步,而非發版。

**鐵則(對應原則 7):image 不存在前,manifest 絕不指向它。** 競態的成因只有一個 — manifest 在 image build 完成前就被改成新 tag,ArgoCD 或 kubectl 一讀就 ImagePullBackOff。解法不是同步邏輯多聰明,而是排序:

> 先 build + push image,確認存在,**才**更新 / 注入 tag。

副作用很好:build 失敗 → 更新 tag 那一步(`needs: build`)根本不會跑 → 叢集維持舊的可用版本。失敗的 build 不可能污染環境。

| 方式 | image tag 存哪 | yaml 寫死 tag? | 誰同步 | 競態防護 |
|------|---------------|----------------|--------|----------|
| push + `--set`(近期) | 不存 git | 否 | 部署時注入 | deploy 在 build 之後跑 |
| ArgoCD,CI 寫回(之後) | env values(git 是源頭) | 是 | CI `yq`+commit | `needs: build` 後才寫 |
| ArgoCD Image Updater | env values | 是 | Image Updater 反查 registry | 只寫 registry 裡存在的 tag |

近期(push):**yaml 不寫死 tag**,部署時用 `--set image.tag=X` 注入(kubectl 用 `kustomize edit set image`)。零 git 寫回、零再觸發、零競態。

之後(ArgoCD):tag 必須在 git。要嘛 CI 在 `needs: build` 之後 `yq` 寫回(commit 加 `[skip ci]`,release-please workflow 設 `paths-ignore: ['deploy/envs/**']` 防再觸發),要嘛交給 ArgoCD Image Updater 從 registry 反查 — 它只寫存在的 image,設計上不可能指向不存在的 tag,最省事。

完整 workflow 範例見第六章。

---

## 六、CI/CD 與部署

部署平台用 GitHub Actions(GHA)。核心設計:**兩種部署方式共用前段 pipeline,只差「最後一哩」**。build、test、release、推 image 全部共用,差別只在誰把版本套到叢集,讓之後從方式一換到方式二不必重構。

### 環境

兩個環境:**stage**(驗證)與 **prd**(正式),可選 dev 做開發自動部署。升環境一律部署同一顆 image,不重新 build。

| 環境 | 觸發 | 閘門 |
|------|------|------|
| stage | tag 建立後自動 | 測試通過 |
| prd | 從 stage 手動 promote | GHA environment 需審核(或 PR 審核) |

GitHub 端用 **Environments** 功能:`production` 設 required reviewers,deploy job 跑到 prd 時會暫停等人核可。

### 共用前段 pipeline

不論哪種部署方式,前段(造出帶版本號的 image)都一樣:

1. `ci.yml`(PR):lint、test、build affected。
2. `release.yml`(main):release-please 切版、打 tag。
3. build job(tag 觸發):build + push `service-a:X.Y.Z`。

差別只在「最後一哩」— 拿到已存在的 image 後怎麼套到叢集(方式一 `helm --set`,方式二寫回 git 給 ArgoCD)。下面範例為了好讀,把 build 與 deploy-stage 放在同一個 workflow;實務上 build job 可抽成 reusable workflow 給兩種方式共用。

`deploy/charts/` 管「怎麼部署」,`deploy/envs/` 管「哪個版本上哪個環境」。兩種方式都吃這份結構。

### 兩種部署方式比較

| 面向 | 方式一:GHA + kubectl/helm(push) | 方式二:ArgoCD(pull / GitOps) |
|------|----------------------------------|-------------------------------|
| 模型 | CI 主動 apply 到叢集 | 叢集主動拉 git 同步 |
| 真實來源 | 最後一次 apply | git repo |
| 叢集憑證 | CI 需要(安全面較大) | CI 不需要(只需 git + registry) |
| 漂移偵測 | 無 | 有,可自動修復 |
| rollback | 重跑舊版 workflow | git revert |
| 稽核 | workflow log | git 歷史 |
| 維運成本 | 低,馬上能動 | 需安裝維護 ArgoCD |
| 即時性 | 立即 | 視 sync 週期(可 webhook 觸發即時) |
| 適合階段 | 起步、團隊小 | 服務變多、要治理 |

**建議路徑**:起步用方式一快速上線,服務數量與治理需求成長後切方式二。因為前段共用且都吃 `deploy/charts` + `deploy/envs`,遷移成本低。

### 方式一:GHA + helm/kubectl(push)

關鍵:**先 build 成功才部署,tag 用 `--set` 注入,manifest 不寫死 tag**(對應原則 7)。

```yaml
# .github/workflows/release-deploy.yml(service-a)
on:
  push:
    tags: ['service-a-v*']            # release-please 打的 tag 觸發

jobs:
  build:
    runs-on: self-hosted
    outputs:
      version: ${{ steps.v.outputs.version }}
    steps:
      - uses: actions/checkout@v4
      - id: v
        run: echo "version=${GITHUB_REF_NAME#service-a-v}" >> "$GITHUB_OUTPUT"
      - run: |
          IMG=registry.local/service-a:${{ steps.v.outputs.version }}
          docker build -t "$IMG" apps/service-a && docker push "$IMG"

  deploy-stage:
    needs: build                      # build 成功才部署 → image 保證存在
    environment: staging
    runs-on: self-hosted
    steps:
      - uses: actions/checkout@v4
      - run: |
          helm upgrade --install service-a deploy/charts/service-a \
            -f deploy/envs/stage/service-a.yaml \
            --set image.tag=${{ needs.build.outputs.version }} \
            -n service-a-stage
```

prd **不綁在同一個 run**,改用獨立、手動觸發、帶版本輸入的 workflow(promotion 模型,避免 run 掛著等核可、也讓 prd 可獨立部署):

```yaml
# .github/workflows/deploy-prd.yml
on:
  workflow_dispatch:
    inputs:
      version: { description: "已在 stage 驗證的版本", required: true }
jobs:
  deploy:
    environment: production           # required reviewers = 上 prod 閘門
    runs-on: self-hosted
    steps:
      - uses: actions/checkout@v4
      - run: |
          helm upgrade --install service-a deploy/charts/service-a \
            -f deploy/envs/prd/service-a.yaml \
            --set image.tag=${{ inputs.version }} \
            -n service-a-prd
```

kubectl 版把部署步驟換成:`kustomize edit set image registry.local/service-a=registry.local/service-a:$VERSION && kubectl apply -k deploy/envs/<env>`。

憑證:CI 需要叢集存取權。地端 vSphere/K8s 建議用 per-environment scoped ServiceAccount token 存成 GHA environment secret,或內網架 self-hosted runner。這是方式一最大的安全面與維運成本。

### 方式二:ArgoCD(pull / GitOps)

CI 不碰叢集,只在 **build 成功後**把 image tag 寫回 git;ArgoCD 偵測後同步:

```yaml
# build 之後才 bump(對應原則 7)
bump-stage:
  needs: build
  steps:
    - uses: actions/checkout@v4
    - run: |
        yq -i '.image.tag = "${{ needs.build.outputs.version }}"' \
          deploy/envs/stage/service-a.yaml
        git commit -am "chore(deploy): service-a ${{ needs.build.outputs.version }} [skip ci]"
        git push
```

`[skip ci]` 加上 release-please workflow 的 `paths-ignore: ['deploy/envs/**']`,寫回就不會回頭觸發發版。prd 用 PR + 審核 merge 來 promote(把同一個 tag 寫進 `deploy/envs/prd/`,PR review 即閘門)。

或直接用 **ArgoCD Image Updater**:從 registry 反查 tag,只寫 registry 裡存在的 image,連寫回 workflow 都省,且設計上不可能造成競態。

ArgoCD 每個(app × env)一個 Application:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: service-a-stage
  namespace: argocd
spec:
  source:
    repoURL: https://github.com/org/monorepo.git
    targetRevision: main
    path: deploy/charts/service-a
    helm:
      valueFiles: [../../envs/stage/service-a.yaml]
  destination:
    server: https://kubernetes.default.svc
    namespace: service-a-stage
  syncPolicy:
    automated: { prune: true, selfHeal: true }
```

app/env 變多時,用 **ApplicationSet** 一個 generator 自動產生每個組合的 Application,不必手寫每一個。

### 前端 web-console

流程骨架相同,差在產物:靜態檔上傳物件儲存 + CDN,或包成 nginx 容器。建議前端也容器化,讓兩種部署方式對前後端一致。

---

## 七、webhook agent 注意事項

| 議題 | 做法 |
|------|------|
| 來源 timeout vs LLM 慢 | 先回 200 ACK,工作丟 queue,背景 worker 跑 agent |
| 重送 | 用事件 ID 去重(冪等) |
| 安全 | 入口驗 HMAC 簽章 |
| 擴展不對稱 | receiver 輕可水平擴展;worker 吃 LLM 配額要控並發 |

v1 先做單一進程(async 背景任務,一顆 image)。等 worker 真的要獨立擴展再拆成 `webhook-agent/{receiver,worker}`(巢狀、共用版本、各自一顆 image)。

---

## 八、建構工具與 CI

| 項目 | 選擇 | 理由 |
|------|------|------|
| 建構協調 | Task + 各語言原生工具 | 零魔法、無鎖定;痛了再評估 Moon |
| affected 偵測 | git-diff 路徑過濾 | 只建/測/部署有變動的單元,含 lib 相依邊 |
| commit 規範 | commitlint + PR title 檢查 | release-please 的前提 |

---

## 九、倉庫慣例:Jenkins、文件、設定

### Jenkins pipeline 放哪?

先問「該不該進這個 monorepo」,再決定位置。**不要放 infra/** — infra/ 是 IaC,放 CI 定義會混淆關注點。

| 選項 | 適用 | 取捨 |
|------|------|------|
| 獨立 repo | Shared Library、給多系統用的大量 legacy | Jenkins 從 repo root 載入 library 最自然;不污染新 monorepo |
| 本 repo `ci/jenkins/` | 只服務本 repo apps 的 pipeline | 集中管理;但 Shared Library 放子目錄對 Jenkins 載入有摩擦 |
| `infra/` | 不建議 | 關注點錯置 |

建議:

- **Jenkins Shared Library 留獨立 repo**。Jenkins 慣例從 repo root 載入 `vars/`、`src/`,放進 monorepo 子目錄會有摩擦。
- **大量服務其他系統的 legacy pipeline 也留獨立 repo**,別把 legacy 倒進乾淨的新 monorepo;要遷移就漸進改寫成 GHA。
- **只有「專屬本 repo apps」的 pipeline 才放 `ci/jenkins/`**。長期目標仍是收斂到 GHA,`ci/` 只是過渡或共存的容身處。

一句話:本 repo 的 CI 是 GHA(放 `.github/workflows/`),Jenkins 是另一個體系,默認分開,真要納入才放 `ci/`,永遠不放 infra/。

### 文件放哪?

跨切面文件集中在 `docs/`,單一 app 的說明就近放在該 app 內。

| 文件類型 | 位置 |
|---------|------|
| repo 入口 / 總覽 | 根目錄 `README.md` |
| 跨切面教學 / 操作指南 | `docs/guides/` |
| 架構決策(本文件) | `docs/adr/` |
| 維運手冊 | `docs/runbooks/` |
| 單一 app 說明 | 該 app 資料夾的 `README.md` |
| 貢獻流程 / 安全政策 | `CONTRIBUTING.md`、`SECURITY.md` |

原則:**就近原則** — 改 app 的人會先看 app 內的 README,跨切面知識才往 `docs/` 找。教學文件屬於跨切面,放 `docs/guides/`。

### 設定檔放哪?

| 設定類型 | 位置 | 原因 |
|---------|------|------|
| 工具設定(lint / format / commit / release) | 根目錄 | 多數工具只認 repo root 或往上找,別藏 |
| app 執行期設定 | 各 app 內 | 隨 app 一起部署 |
| 環境部署值 | `deploy/envs/` | 已規劃 |
| 共用程式設定(常數 / flag) | `libs/` 或各 app | 那是 code,不是「設定檔」 |

工具設定檔(`.editorconfig`、`commitlint.config.js`、`.golangci.yml`、`ruff` 設定等)一律放根目錄,因為 lint/format 工具預設從 repo root 或往上搜尋。不要為了「乾淨」把它們塞進子目錄,那會讓工具找不到。

### 業界標準的 meta 檔

一個完整 repo 該有的治理檔:

| 檔案 | 用途 |
|------|------|
| `README.md` | 入口:這是什麼、怎麼開始 |
| `CONTRIBUTING.md` | 開發流程、commit 規範、PR 規則 |
| `SECURITY.md` | 漏洞回報管道 |
| `.github/CODEOWNERS` | 每個目錄的擁有 team,自動指派 reviewer |
| `.github/pull_request_template.md` | 統一 PR 格式 |
| `.github/dependabot.yml` | 相依更新 |
| `.editorconfig` / `.gitattributes` | 跨編輯器 / 跨平台一致性 |

CODEOWNERS 在 monorepo 特別重要:它把「擁有權邊界」落實到目錄層級,讓不同 team 對自己的 app 有 review 權,避免 monorepo 變成沒人負責的公共地。

---

## 十、重大決策總表

| 決策 | 選擇 | 理由 |
|------|------|------|
| 組織軸線 | 依部署單元,apps 平鋪 | 用最穩定的軸,避開易變的技術層 |
| 巢狀時機 | 只在一起切版一起上線 | 巢狀=版本綁定,非打包 |
| 共用 code | libs/ 按語言分,live at head | 不發版、改了立即生效 |
| 跨語言共用 | 共用 proto 合約,非 code | 不同語言無法直接共用 |
| agent 位置 | apps/ | 它是可部署服務,非按技術分類 |
| 版本策略 | 獨立版本 + 強耦合才 linked | 部署解耦 |
| 發版工具 | release-please manifest mode | 多語言、各 package 各自 release-type |
| release vs deploy | 拆開 | release 造 artifact,deploy 選版本上環境 |
| 部署平台 | GitHub Actions | 與 GHES 整合 |
| 部署方式 | 起步 push(GHA+helm),成長後切 ArgoCD | 前段共用,遷移成本低 |
| image tag 同步 | build 成功後才注入/寫回;manifest 不在 image 存在前指向它 | 消除部署競態(ImagePullBackOff) |
| 環境升級 | 部署同一顆 image,不重 build(push 用 `--set` 注入,ArgoCD 改 git tag) | 同一顆 image 跨 stage/prd |
| prd 部署 | 獨立 workflow_dispatch + environment 審核(promotion 模型) | prd 可獨立部署,不被 stage run 拖著、不掛著等核可 |
| Helm chart 版本 | 起步不列為 release package | 一個 app 一個版本,避免雙 PR |
| 部署模型 | GitOps(charts + envs) | 升版只動 envs,不重 build |
| internal lib 版本 | 不發版,靠 affected CI 把關 | 版本是給出貨/外部 pin 用的 |
| 建構工具 | Task + 原生工具起步 | 避免單行道,被痛逼到再升級 |
| Jenkins 位置 | 不放 infra;Shared Library/legacy 留獨立 repo,專屬本 repo 才 `ci/jenkins/` | CI 定義非 IaC;避免污染新 repo |
| 文件 | 跨切面放 `docs/`,app 說明就近放 app 內 | 就近原則 |
| 工具設定檔 | 放根目錄 | 工具只認 repo root |
| 擁有權 | `.github/CODEOWNERS` 落實到目錄 | monorepo 需明確邊界,避免無人負責 |
