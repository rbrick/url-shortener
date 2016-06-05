window.addEventListener('load', function() {
    var shortenForm = document.getElementById('shortenForm');
    shortenForm.addEventListener('submit', function(event) {
      event.preventDefault();
      onSubmit();
    });
})


function onSubmit() {
  var longUrl = document.getElementById('longUrl');
  var errorHolder = document.getElementById('errorHolder');
  var shortened = document.getElementById('shortened');
  // console.log('LETS SEE HOW FAR WE\'VE COME');
  create(longUrl.value, function(event){
    console.log('LETS SEE HOW FAR WE\'VE COME');
    update(JSON.parse(event.responseText));
  });
}

function update(res) {
  if (res['error']) {
    shortened.innerHTML = '';
    errorHolder.innerHTML += `<div class="alert alert-danger alert-dismissible" role="alert">
<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>
<strong>Error!</strong> ${res['error']}
</div>`;
  } else {
     errorHolder.innerHTML = '';
     shortened.innerHTML = '/' + res['id'];
  }
}

function lookup(id, callback) {
   makeGet('/api/lookup', {'q': id}, callback);
}

function create(longUrl, callback) {
  makePost('/api/create', {'longUrl': longUrl}, callback);
}

function onLookup(event) {
    console.log(event.responseText);
    var res = JSON.parse(event.responseText);
    if (res['error']) {
       alert(res['error']);
    } else {
       alert(res['longUrl']);
    }
}
