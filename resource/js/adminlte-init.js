// for file input
$(function() {
    bsCustomFileInput.init();
});
// select2
$('.select2').select2({
    theme: 'bootstrap4'
});
// label element required
$('label.required').each(function(index, element) {
    let $element = $(element);
    const text = $element.text();
    $element.html(text + '<span class="badge badge-danger ml-1">必須</span>');
});
// popover
$('[data-toggle="hint-popover"]').popover({
    trigger: 'focus'
});