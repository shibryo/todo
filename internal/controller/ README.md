# Controller
controllerを記述するために意識していることを記述する。

参考に、Modelで記述するべきことをWikipediaのMVCモデルの記述から引用する
> ユーザからの入力（通常イベントとして通知される）をモデルへのメッセージへと変換してモデルに伝える要素である。すなわち、UIからの入力を担当する。モデルに変更を引き起こす場合もあるが、直接に描画を行ったり、モデルの内部データを直接操作したりはしない。

私自身は、UIの入出力とモデルへのメッセージのやり取りを表現するレイヤーとして使っている。

## DTO
まず、データの入出力には専用のModelを使います。
一つ目の理由は、Modelで定義してるTodoとTodoViewではオブジェクトを使う目的が違うからです。前者は業務ロジックを表現するために、後者はUIで入出力を行うためだからです。
二つ目の理由は、入出力内容の変更によってModelが影響を受けないようにするためです。
三つ目の理由は、入出力に必要なデータはStingやint非常に簡単な型で表されるからです。もし、Modelのモデルと入出力が同じモデルを使う場合は強制的にModelのモデルがstringやintで表され、Modelの表現力が弱まります。コードを読んでも意味が伝わりづらくなります。

以下のコードではModelのモデルで受け取った情報をUI専用のViewModelに移し替えて返しています。
```go
func(ctrl *TodoController) FindTodoByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
        // Modelのモデルで値を取得
		todo, err := ctrl.todoRepository.FindByID(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
        // View用のモデルに詰め替えてWebUIへ返す。
		todoView := TodoView{
			ID:         uint64(todo.ID),
			Title:      todo.Title.AsGoString(),
			Completed:  todo.Completed.AsGoBool(),
			LastUpdate: todo.LastUpdate.AsGoString(),
			CreatedAt:  todo.CreatedAt.AsGoString(),
		}
		return c.JSON(http.StatusOK, todoView)
	}
}
```



参考：
データ詰め替え戦略
https://scrapbox.io/kawasima/%E3%83%87%E3%83%BC%E3%82%BF%E8%A9%B0%E3%82%81%E6%9B%BF%E3%81%88%E6%88%A6%E7%95%A5


## 「計算と判定のモデル」と「データの記録と参照」を分ける

```go
func(ctrl *TodoController) UpdateTodo() echo.HandlerFunc {
	return func(c echo.Context) error {
        id, err := strconv.ParseUint(c.Param("id"), 10, 64)
        // ....

		err = ctrl.todoRepository.Update(todo)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		resultTodoView := TodoView{
			ID:         uint64(todo.ID),
			Title:      todo.Title.AsGoString(),
			Completed:  todo.Completed.AsGoBool(),
			LastUpdate: todo.LastUpdate.AsGoString(),
			CreatedAt:  todo.CreatedAt.AsGoString(),
		}
		return c.JSON(http.StatusOK, resultTodoView)
	}
}
```

参考：
ドメイン駆動設計を理解する３つのキーワード
https://masuda220.hatenablog.com/entry/2019/03/04/144154