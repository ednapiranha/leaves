(function() {
  var deletes = document.querySelectorAll('.delete');

  for (var i = 0; i < deletes.length; i++) {
    console.log(deletes[i])
    deletes[i].onclick = function (ev) {
      ev.preventDefault();

      var delRev = confirm("Are you sure you want to delete this review?");

      if (delRev) {
        document.location.href = this.href;
      }
    }
  }
})()
