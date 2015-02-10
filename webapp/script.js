$(function() {
  $('.add').click(function() {
    edit({
      name: "My Security Light",
      sensors: [],
      lights: [],
      timeout: 5
    })
  });

  $('.save').click(function(e) {
    e.preventDefault()
    // TODO: Save the thing

    var securityLight = {lights:[], sensors:[]}

    $('form').serializeArray().forEach(function(field) {
      if (securityLight[field.name]) {
        securityLight[field.name].push(field.value)
      } else {
        securityLight[field.name] = field.value
      }
    });

    if (securityLight.active == 'always') {
      delete(securityLight.timeStart)
      delete(securityLight.timeEnd)
    }

    delete(securityLight.active)

    console.log("saving", securityLight)

    $.post('/api/security-lights', securityLight, function(_, status, response) {

        if (status != 'success') {
          alert("Failed to save: " + response.responseText)
        } else {
          list();
        }
    }, 'json');

  });

  $('.cancel').click(function(e) {
    e.preventDefault()
    list()
  });

  $(document).on('click', 'button.edit', function(e){
    edit(securityLights[$(e.target).data('id')]);
  })

  $(document).on('click', 'button.delete', function(e){
    var id = $(e.target).data('id');
    $.ajax({
      type: 'DELETE',
      url: '/api/security-lights/' + id,
      dataType: 'json',
      success: function() {
        list()
      }
    });
  })
})

var securityLights = {};

function list() {

  $('#edit').hide()
  $('#list').show()

  $('#securityLights').empty();
  securityLights = {};
  $.get("/api/security-lights", function(l) {
    console.log("got lights", l)
    l.forEach(function(light) {
      securityLights[light.id] = light;
      $('#securityLights').append('<li>' + light.name + ' <button class="edit" data-id="' + light.id + '">Edit</button> <button class="delete" data-id="' + light.id + '">Delete</button></li>');
    })
  });

}

function edit(securityLight) {

  function fill(el, url, selected) {
    el.empty();
    $.get(url, function(items) {
      items.forEach(function(item) {
        el.append('<option value="' + item.id + '"' + (selected.indexOf(item.id) > -1?' selected':'') + '>' + item.name + '</option>');
      })
    })
  }

  fill($('[name=sensors]'), "/api/sensors", securityLight.sensors)
  fill($('[name=lights]'), "/api/lights", securityLight.lights)

  $('[name=id]').val(securityLight.id || "")
  $('[name=name]').val(securityLight.name || "")
  $('[name=timeout]').val(securityLight.timeout || "5")
  if (securityLight.timeStart) {
      $('[name=active][value=between]').click()
      $('[name=timeStart]').val(securityLight.timeStart)
      $('[name=timeEnd]').val(securityLight.timeEnd)
  } else {
      $('[name=active][value=always]').click()
      $('[name=timeStart]').val('sunrise')
      $('[name=timeEnd]').val('sunset')
  }

  $('#list').hide()
  $('#edit').show()
}

// Times dropdowns
function pad(a,b){return(1e15+a+"").slice(-b)}

var times = ['midnight', 'sunrise', 'dawn', 'midday', 'sunset', 'dusk']

for (var h = 0; h < 24; h++) {
  for (var m = 0; m < 60; m += 15) {
    times.push(pad(h, 2) + ':' + pad(m, 2));
  }
}
$(function(){
  var els = $('.times');
  times.forEach(function(t) {
    els.append('<option>'+t+'</option>');
  });

  /* Don't work right with 'sunrise' etc.
  var endOptions = $('[name=timeEnd] option');

  $('[name=timeStart]').change(function(e) {
    console.log("Change", e);
    var selected = $(e.srcElement).find('option').not(function(){ return !this.selected })[0].index;

    endOptions.forEach(function(el, index) {
      if (index > selected) {
        $(el).removeAttr('disabled');
      } else {
        $(el).attr('disabled', 'disabled');
      }
    });

    if (endOptions.not(function(){ return !this.selected })[0].disabled) {
      endOptions.not(function(){ return this.disabled})[0].selected = true
    }
  });
  */
})






// Http mocking
if (location.protocol != "file:") {
  list();
} else {
  mock();
}
