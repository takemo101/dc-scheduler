/**
 * ダイアログ要素生成クラス
 */
class DialogElementBuilder {
  constructor() {
    //
  }

  /**
   * ボタン要素生成
   */
  buildButtonElement() {
    return $(
      '<button data-modal="remove" type="button" style="display: none;" data-toggle="modal" data-target="#dialog-modal"></button>'
    );
  }

  /**
   * HTML要素生成
   */
  buildHTMLElement(options = {}) {
    options = $.extend(
      {
        title: "タイトル",
        message: "メッセージ",
        yes: "はい",
        no: "いいえ",
        class: "",
      },
      options
    );

    var noHtml = options.no
      ? '<button type="button" class="btn btn-secondary btn-outline-light" data-dismiss="modal">' +
        options.no +
        "</button>"
      : "";
    var yesHtml =
      '<button type="button" class="btn btn-primary btn-outline-light" data-dismiss="modal" data-yes="true">' +
      options.yes +
      "</button>";
    var html =
      '<div id="dialog-modal" style="display: none;" class="modal fade" data-modal="remove" style="display: block; padding-right: 17px;" aria-modal="true">' +
      '<div class="modal-dialog">' +
      '<div class="modal-content ' +
      options.class +
      '">' +
      '<div class="modal-header">' +
      '<h4 class="modal-title">' +
      options.title +
      "</h4>" +
      '<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">×</span></button>' +
      "</div>" +
      '<div class="modal-body"><p>' +
      options.message +
      "</p></div>" +
      '<div class="modal-footer justify-content-between">' +
      noHtml +
      yesHtml +
      "</div>" +
      "</div>" +
      "</div>" +
      "</div>";
    return $(html);
  }
}

/**
 * ダイアログ管理クラス
 */
class DialogUtility {
  constructor() {
    this.can = false;
    this.builder = new DialogElementBuilder();
  }

  /**
   * submit 状態かどうか
   */
  canSubmit() {
    return this.can;
  }

  /**
   * submit 状態をリセット
   */
  resetSubmit() {
    this.can = false;
  }

  /**
   * submit 実行
   */
  executeSubmit($form, element) {
    if ($form.reportValidity() && !this.canSubmit()) {
      this.can = true;
      $form.submit();
    }
  }

  /**
   * ダイアログの作成
   */
  createDefaultDialog(element, options = {}) {
    const self = this;

    $('[data-modal="remove"]').remove();

    const $form = $($(element).data("form"));

    // フォームをチェック
    if (!$form.reportValidity()) {
      return;
    }

    const $btn = self.builder.buildButtonElement();

    $("body").append($btn);
    $("body").append(
      self
        .buildHTMLElement(options)
        .find('[data-yes="true"]')
        .on("click", function () {
          self.executeSubmit($form, element);
        })
    );

    $btn.trigger("click");
  }

  /**
   * 削除ダイアログ
   */
  delete(element) {
    this.createDefaultDialog(element, {
      title: "DELETE",
      message: "削除しますか？",
    });
  }

  /**
   * 復元ダイアログ
   */
  restore(element) {
    this.createDefaultDialog(element, {
      title: "RESTORE",
      message: "復元しますか？",
    });
  }

  /**
   * 更新ダイアログ
   */
  update(element) {
    this.createDefaultDialog(element, {
      title: "UPDATE",
      message: "更新しますか？",
    });
  }

  /**
   * 追加ダイアログ
   */
  store(element) {
    this.createDefaultDialog(element, {
      title: "STORE",
      message: "追加しますか？",
    });
  }

  /**
   * メッセージダイアログ
   */
  message(element) {
    const self = this;

    $('[data-modal="remove"]').remove();

    const $form = $($(element).data("form"));

    // フォームをチェック
    if (!$form.reportValidity()) {
      return;
    }

    const $btn = this.builder.buildButtonElement();

    $("body").append($btn);
    $("body").append(
      self
        .buildHTMLElement({
          title: $(element).data("title"),
          message: $(element).data("message"),
        })
        .find('[data-yes="true"]')
        .on("click", function () {
          self.executeSubmit($form, element);
        })
    );

    $btn.trigger("click");
  }

  /**
   * インフォダイアログ
   */
  info(element) {
    var self = this;

    $('[data-modal="remove"]').remove();

    $("body").append(
      self.buildHTMLElement({
        title: $(element).data("title"),
        message: $(element).data("message"),
        no: null,
      })
    );
  }
}

// ブラウザ上では$Dialogで利用できる
var $Dialog = new DialogUtility();
window.$Dialog = $Dialog;

$(function () {
  /**
   * data-class="js-dialog" と data-type="" は必ず設定する
   */
  $('[data-class="js-dialog"]').on("click", function (e) {
    if (!$Dialog.canSubmit()) {
      var type = $(this).data("type");
      switch (type) {
        case "delete":
          $Dialog.delete(this);
          break;
        case "restore":
          $Dialog.restore(this);
          break;
        case "update":
          $Dialog.update(this);
          break;
        case "store":
          $Dialog.store(this);
          break;
        case "message":
          $Dialog.message(this);
          break;
        case "info":
          $Dialog.info(this);
          break;
        default:
      }

      return false;
    }
  });
});
