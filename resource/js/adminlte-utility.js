/**
 * AdminLTEプラグインのためのユーティリティークラス
 */
class AdminLTEUtility {
  constructor() {
    //
  }

  /**
   * create ck editor
   */
  summernote(editorID, options = {}) {
    return $("#" + editorID).summernote(options);
  }

  /**
   * 指定のIDの要素の中にソート値のinput hiddenを出力する
   */
  createSortInput(submitID, inputName = "ids", dataID = "id") {
    const $submit = $("#" + submitID);

    $submit.on("click", function () {
      const $form = $(this).parent("form");
      const elementID = $(this).data("element-id");

      $form.children('input[name^="' + inputName + '"]').remove();
      (elementID
        ? $("#" + elementID + " [data-" + dataID + "]")
        : $("[data-" + dataID + "]")
      ).each(function (index, element) {
        const id = $(element).data(dataID);
        $form.append(
          $("<input>").val(id).attr("type", "hidden").attr("name", inputName)
        );
      });
    });
  }
}

// ブラウザ上では$Utilityで利用できる
window.$Utility = new AdminLTEUtility();
