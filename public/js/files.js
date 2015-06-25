var editable = ['html', 'htm', 'css', 'js', 'java', 'txt'];

function isGoodName(name) {
  return name.match(/^\w[-.\w]*$/);
}

function isEditable(name) {
  if (! isGoodName(name)) {
    return false;
  }
  var index = name.lastIndexOf('.');
  if (index > 0) {
    var ext = name.slice(index + 1).toLowerCase();
    if (editable.indexOf(ext) >= 0) {
      return true;
    }
  }
  return false;
}

function createFile() {
  var name = prompt('Create new file as:\n(must end with .html, .css, .js, .java, or .txt)', '');
  if (name == null || name == '') {
    return;
  }
  if (! isEditable(name)) {
    alert('Sorry, bad filename.');
    return;
  }
  location.href = '/edit/' + name;
}

function copyFile(button) {
  var source = button.getAttribute('data-name');
  var name = prompt('Copy "' + source + '" to:', '');
  if (name == null || name == '') {
    return;
  }
  if (! isGoodName(name)) {
    alert('Sorry, bad filename.');
    return;
  }
  document.getElementById('copy-form-source').value = source;
  document.getElementById('copy-form-dest').value = name;
  document.getElementById('copy-form').submit();
}
