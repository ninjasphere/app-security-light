$(function() {
  $('.add').click(function() {
    edit({
      name: "My Security Light",
      enabled: true,
      sensors: [],
      lights: [],
      timeStart: "sunrise",
      timeEnd: "sunset",
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

    securityLight.timeout = parseInt(securityLight.timeout)

    securityLight.enabled = securityLight.enabled === "true"

    delete(securityLight.active)

    console.log("saving", securityLight)

    $.ajax({
      type: "POST",
      url: '/api/security-lights',
      data: JSON.stringify(securityLight),
      dataType: 'json',
      success: function(_, status, response) {

        if (status != 'success') {
          alert("Failed to save: " + response.responseText)
        } else {
          list();
        }
      }
    });

  });

  $('.cancel').click(function(e) {
    e.preventDefault()
    list()
  });

  $(document).on('click', '.things .thing', function(e){
    var t = $(e.target).parents('.thing');

    if (!$(e.target).is('input')) {
      t.find('input').click()
    }

    t.toggleClass('selected', t.find('input').prop('checked'))
  })

  $(document).on('click', '.activeButtons .button', function(e){
    var t = $(e.target);

    if (!t.hasClass('button')) {
      t = t.parents('.button');
    }
    t.addClass('selected')
    $('.activeButtons .button').not(t).removeClass('selected').find('input').removeAttr('checked')
    t.find('input').attr('checked', 'checked')
  })

  $(document).on('click', 'a.edit', function(e){
    edit(securityLights[$(e.target).data('id')]);
  })

  $(document).on('click', '.delete', function(e){
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
  $('#list').css('display', 'block')
  refreshSecurityLights()
}

//setInterval(refreshSecurityLights, 5000);
function refreshSecurityLights() {
  $.get("/api/security-lights", function(l) {

    securityLights = l;
    $('#securityLights').empty();

    console.log("got lights", l)
    for (var id in l) {
      light = securityLights[id];
      $('#securityLights').append('<li><span class="name">' + light.name + '</span> <a href="#" class="edit" data-id="' + light.id + '">edit</a></li>');
    }
  });
}

function edit(securityLight) {
  var thingTpl = $($('.thing')[0])

  function fill(el, url, selected, name) {
    el.empty();

    $.get(url, function(items) {
      items.forEach(function(item) {
        var t = thingTpl.clone();
        t.find('input').attr('value', item.id).attr('name', name);
        t.find('.name').text(item.name)
        t.toggleClass('selected', selected.indexOf(item.id) > -1);
        t.find('input').prop('checked', selected.indexOf(item.id) > -1);

        el.append(t);
      })
    })
  }

  fill($('#sensors'), "/api/sensors", securityLight.sensors, 'sensors')
  fill($('#lights'), "/api/lights", securityLight.lights, 'lights')

  $('[name=id]').val(securityLight.id || "")
  $('.delete').data('id', securityLight.id)

  $('.delete').toggle(!!securityLight.id)

  $('[name=enabled]').val(securityLight.enabled || "true")
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

var times = ['dawn', 'sunrise', 'sunset', 'dusk']

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
})





$(function() {
  // Http mocking
  if (location.protocol != "file:") {
    list();
  } else {
    mock();
  }
})
