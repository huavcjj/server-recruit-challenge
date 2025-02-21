## 課題 5

### 1. データベースの整合性
テーブル `albums` の `singer_id` に `FOREIGN KEY` 制約があるが、`ON DELETE` / `ON UPDATE` が指定されていません。

#### **改善策:**
- `ON DELETE CASCADE ON UPDATE CASCADE`を追加し、歌手が削除または更新された際に、関連するアルバムも適切に削除/更新されるようにする。

  `FOREIGN KEY (singer_id) REFERENCES singers(id) ON DELETE CASCADE ON UPDATE CASCADE`

**データの変更履歴の管理:**
- デバッグや変更履歴の管理のため、`albums` および `singers` テーブルに `created_at`、`updated_at` カラムを追加する。
- `TIMESTAMP` 型を使用するとストレージサイズを抑えられ、インデックスの効率も向上する。
- ただし、2038年問題を考慮する場合、長期的なデータ保存が必要なら `DATETIME` を使用する。

`
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;
`

**検索速度の向上:**
- `albums` テーブルの `singer_id` に `INDEX` を追加し、検索の最適化を図る。

```sql
CREATE INDEX idx_singer_id ON albums(singer_id);
```

### 2. 環境変数の管理
現在、データベース接続情報やサーバーのポート番号などの設定情報がハードコードされている。

**改善策:**
- `.env` ファイルを導入し、環境変数として管理することで、環境ごとに設定を変更できるようにする。
- これにより、開発環境・本番環境で異なる設定を簡単に適用可能。

### 3. フロントエンドとのデータフォーマットの分離
フロントエンドに対し、内部のデータ構造（`Model`）をそのまま返すと、変更に弱くなる。

**改善策:**
- DTO（Data Transfer Object）を通じてAPIのレスポンス形式を固定する。
- 例えば、`model.Album` に `price` を追加した場合、DTO を利用すればフロントエンド側に影響を与えずに済む。
- DTO 層でバリデーションを一元化し、無効なデータが Service 層に流れるのを防ぐ。

### 4. Model層のValidate()をDTO層に移動
Model層のデータは内部のビジネスロジックやDBとのやり取りに使われるべきで、バリデーションを直接持つべきではない。

**改善策:**
- DTO層にValidate()メソッドを追加する、またはバリデーションライブラリを導入して、DTO層でバリデーションを行う。
- これにより、Model層はデータの保持と変換に専念し、DTO層はデータのバリデーションに専念する。

### 5. AlbumのJSONレスポンスの最適化
現在のJSONレスポンスは以下のような構造になっている。これは、ネストが深くなり、フロントエンドでのデータ取得が複雑になる。

```json
[
  {
    "id": 1,
    "title": "Alice's 1st Album",
    "singer": {
      "id": 1,
      "name": "Alice"
    }
  },
  {
    "id": 2,
    "title": "Alice's 2nd Album",
    "singer": {
      "id": 1,
      "name": "Alice"
    }
  },
  {
    "id": 3,
    "title": "Bella's 1st Album",
    "singer": {
      "id": 2,
      "name": "Bella"
    }
  }
]
```

**改善策:**
- フロントエンドでのデータ取得を容易にするため、DTO層を介して以下のようなフラットな構造に変更する。
- これにより、フロントエンド側でのデータ取得が容易になり、ネストが深くなるのを防ぐ。

```json
[
  {
    "id": 1,
    "title": "Alice's 1st Album",
    "singer_id": 1,
    "singer_name": "Alice"
  },
  {
    "id": 2,
    "title": "Alice's 2nd Album",
    "singer_id": 1,
    "singer_name": "Alice"
  },
  {
    "id": 3,
    "title": "Bella's 1st Album",
    "singer_id": 2,
    "singer_name": "Bella"
  }
]
```

### 6. セキュリティ対策
ingerIDやAlbumIDがintのため、IDの予測が容易。これにより、セキュリティ上のリスクが生じる可能性がある。

**改善策:**
- UUIDを使用して、IDの予測を困難にする。しかし、UUIDは文字列であり、データベースのインデックス効率が低下する。APIのURLにUUIDを含めると、URLが長くなり、可読性が低下する。

### 7. Service層でのバリデーション
Controller層はリクエストの受け取りとレスポンスの返却に集中すべきで、バリデーションの責務を持つべきではない。
(現時点では、Controller層でDTOのバリデーションを行っている。これは、Service層のインターフェースが Model 層に直接データを渡す設計になっているため、現在のコード構造を壊さずに実装されている。)

**改善策:**
- DTOのバリデーションをService層に移動し、リクエストのバリデーションを行う。
- これにより、Controller層がシンプルになり、コードの責務が明確になる。

### 8. API 仕様の自動生成とドキュメント化
フロントエンド開発者が API 仕様を把握しやすくするため、OpenAPI仕様に従ってAPIのドキュメントを自動生成する。

**改善策:**
- OpenAPI 仕様を採用し、API のドキュメントを自動生成する。
- OpenAPI のスキーマにバリデーションを定義することで、すべてのリクエストを自動的にチェック可能。
- OpenAPI で定義したスキーマが Go の `model` や `dto` にそのまま影響を与えるため、DTO 層を併用して柔軟なデータ変換を実現する。


### 9. OpenAPI と DTO の併用
OpenAPIでAPIの仕様を定義することで、リクエストの型やバリデーションを自動管理できるが、DTO層を併用することで、APIのレスポンスデータを柔軟に管理できる。

**改善策:**
- OpenAPIでAPI仕様を定義し、DTOを使ってGo内部のデータ変換を行う。
- これにより、API 仕様の標準化と柔軟なデータ変換を両立できる。

