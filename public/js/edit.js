var editor = ace.edit('editor');

function xhr(method, url, fn) {
  var xhr = new XMLHttpRequest();
  xhr.open(method, url, true);
  xhr.setRequestHeader('Content-Type', 'text/plain');
  xhr.setRequestHeader('X-Requested-With', 'XMLHttpRequest');
  xhr.onreadystatechange = function () {
    if (xhr.readyState == 4) {
      fn(xhr.status, xhr.responseText);
    }
  };
  return xhr;
}

function start() {
  var url = location.href.replace('/edit/', '/edit-get/');
  xhr('GET', url, function (status, text) {
    if (status == 200) {
      editor.setValue(text);
      editor.clearSelection();
      editor.scrollToLine(0);
      editor.gotoLine(1);
      editor.focus();
    } else {
      editor.setValue('ERROR!');
    }
  }).send();
}

function save() {
  var url = location.href.replace('/edit/', '/edit-put/');
  xhr('PUT', url, function (status, text) {
    if (status == 200) {
      alert('save ok');
    } else {
      alert('save error!');
    }
  }).send(editor.getValue());
}

editor.commands.addCommand({
  name: 'save',
  bindKey: { win: 'Ctrl-S', mac: 'Command-S' },
  exec: save,
  readOnly: false
});

start();
