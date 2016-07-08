(function() {
  // Delete review
  var deletes = document.querySelectorAll('.delete');

  for (var i = 0; i < deletes.length; i++) {
    deletes[i].onclick = function (ev) {
      ev.preventDefault();

      var delRev = confirm("Are you sure you want to delete this review?");

      if (delRev) {
        document.location.href = this.href;
      }
    }
  }

  // Like review
  var like = document.querySelector('.like');
  var xmlhttp = new XMLHttpRequest();

  like.onclick = function (ev) {
    ev.preventDefault();

    xmlhttp.onreadystatechange = function() {
      if (xmlhttp.readyState == XMLHttpRequest.DONE) {
        if (xmlhttp.status == 200) {
          if (like.classList.contains('on')) {
            like.classList.remove('on');
          } else {
            like.classList.add('on');
          }
        }
      }
    }

    xmlhttp.open("GET", like.href, true);
    xmlhttp.send();
  }
})()
